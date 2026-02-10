package transport

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	pcrypto "pulse/internal/crypto"
	"pulse/internal/group"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	rclient "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	ma "github.com/multiformats/go-multiaddr"
)

// SendProgress is sent on the progress channel during file transfer.
type SendProgress struct {
	PeerID string
	Done   bool
	Err    error
}

// ReceiveEvent is emitted when a file is received.
type ReceiveEvent struct {
	Filename string
	Size     int64
	From     string
	Err      error
}

// Header is the wire format for a file transfer.
type Header struct {
	Filename string
	Size     int64
	Hash     []byte // BLAKE3 hash
}

// SendFile sends a file to all members of a group via the relay.
func SendFile(ctx context.Context, priv crypto.PrivKey, g *group.Group, filePath string, progress chan<- SendProgress) error {
	defer close(progress)

	h, err := libp2p.New(libp2p.Identity(priv), libp2p.EnableRelay())
	if err != nil {
		return fmt.Errorf("creating host: %w", err)
	}
	defer h.Close()

	if err := connectToRelay(ctx, h, g.Relay); err != nil {
		return fmt.Errorf("connecting to relay: %w", err)
	}

	// Read file and compute hash
	payload, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	hash := pcrypto.HashBytes(payload)
	filename := filepath.Base(filePath)

	var wg sync.WaitGroup
	for _, pid := range g.Members {
		wg.Add(1)
		go func(peerIDStr string) {
			defer wg.Done()
			err := sendToPeer(ctx, h, g, peerIDStr, filename, payload, hash)
			progress <- SendProgress{PeerID: peerIDStr, Done: err == nil, Err: err}
		}(pid)
	}
	wg.Wait()

	return nil
}

func sendToPeer(ctx context.Context, h host.Host, g *group.Group, peerIDStr string, filename string, payload []byte, hash []byte) error {
	full := fmt.Sprintf("%s/p2p-circuit/p2p/%s", g.Relay, peerIDStr)
	maddr, err := ma.NewMultiaddr(full)
	if err != nil {
		return fmt.Errorf("parsing multiaddr: %w", err)
	}

	destInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return fmt.Errorf("parsing peer addr: %w", err)
	}

	// Connect with retry
	var connectErr error
	for attempt := 0; attempt < 3; attempt++ {
		h.Peerstore().AddAddrs(destInfo.ID, destInfo.Addrs, time.Minute)
		connectErr = h.Connect(ctx, *destInfo)
		if connectErr == nil {
			break
		}
		time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
	}
	if connectErr != nil {
		return fmt.Errorf("connecting to peer: %w", connectErr)
	}

	// Wait for connection
	for i := 0; i < 15; i++ {
		if h.Network().Connectedness(destInfo.ID) == network.Connected {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	s, err := h.NewStream(ctx, destInfo.ID, protocol.ID(g.Protocol))
	if err != nil {
		return fmt.Errorf("opening stream: %w", err)
	}
	defer s.Close()

	w := bufio.NewWriter(s)

	// Write header: filename\n
	fmt.Fprintln(w, filename)

	// Write size (8 bytes big endian)
	sizeBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(sizeBuf, uint64(len(payload)))
	w.Write(sizeBuf)

	// Write hash (32 bytes BLAKE3)
	w.Write(hash)

	// Write payload
	w.Write(payload)
	return w.Flush()
}

func connectToRelay(ctx context.Context, h host.Host, relayAddr string) error {
	relayMA, err := ma.NewMultiaddr(relayAddr)
	if err != nil {
		return fmt.Errorf("parsing relay address: %w", err)
	}

	relayInfo, err := peer.AddrInfoFromP2pAddr(relayMA)
	if err != nil {
		return fmt.Errorf("parsing relay peer info: %w", err)
	}

	// Retry with exponential backoff
	var connectErr error
	for attempt := 0; attempt < 4; attempt++ {
		connectErr = h.Connect(ctx, *relayInfo)
		if connectErr == nil {
			return nil
		}
		time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
	}
	return fmt.Errorf("relay connection failed after 4 attempts: %w", connectErr)
}

// ListenResult holds the outcome of a listen session.
type ListenResult struct {
	Events <-chan ReceiveEvent
	Stop   func()
}

// Listen starts listening for incoming files on a group protocol.
func Listen(ctx context.Context, priv crypto.PrivKey, g *group.Group, storeDir string) (*ListenResult, error) {
	if err := os.MkdirAll(storeDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating store directory: %w", err)
	}

	h, err := libp2p.New(libp2p.Identity(priv), libp2p.EnableRelay(), libp2p.EnableRelayService())
	if err != nil {
		return nil, fmt.Errorf("creating host: %w", err)
	}

	if err := connectToRelay(ctx, h, g.Relay); err != nil {
		h.Close()
		return nil, fmt.Errorf("connecting to relay: %w", err)
	}

	relayMA, _ := ma.NewMultiaddr(g.Relay)
	relayInfo, _ := peer.AddrInfoFromP2pAddr(relayMA)

	_, err = rclient.Reserve(ctx, h, *relayInfo)
	if err != nil {
		h.Close()
		return nil, fmt.Errorf("relay reservation failed: %w", err)
	}

	events := make(chan ReceiveEvent, 16)

	// Relay reservation renewal goroutine
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(90 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				_, err := rclient.Reserve(ctx, h, *relayInfo)
				if err != nil {
					events <- ReceiveEvent{Err: fmt.Errorf("relay renewal failed: %w", err)}
				}
			}
		}
	}()

	// Stream handler
	h.SetStreamHandler(protocol.ID(g.Protocol), func(s network.Stream) {
		defer s.Close()

		remotePeer := s.Conn().RemotePeer().String()

		// Verify sender is a group member
		isMember := false
		for _, m := range g.Members {
			if m == remotePeer {
				isMember = true
				break
			}
		}
		if !isMember {
			events <- ReceiveEvent{Err: fmt.Errorf("rejected connection from non-member: %s", remotePeer)}
			return
		}

		reader := bufio.NewReader(s)

		// Read filename
		filename, err := reader.ReadString('\n')
		if err != nil {
			events <- ReceiveEvent{Err: fmt.Errorf("reading filename: %w", err)}
			return
		}
		filename = strings.TrimSpace(filename)
		filename = filepath.Base(filename) // sanitize

		// Read size (8 bytes)
		sizeBuf := make([]byte, 8)
		if _, err := io.ReadFull(reader, sizeBuf); err != nil {
			events <- ReceiveEvent{Err: fmt.Errorf("reading size: %w", err)}
			return
		}
		size := int64(binary.BigEndian.Uint64(sizeBuf))

		// Read hash (32 bytes)
		hashBuf := make([]byte, 32)
		if _, err := io.ReadFull(reader, hashBuf); err != nil {
			events <- ReceiveEvent{Err: fmt.Errorf("reading hash: %w", err)}
			return
		}

		// Read payload
		path := filepath.Join(storeDir, filename)
		f, err := os.Create(path)
		if err != nil {
			events <- ReceiveEvent{Err: fmt.Errorf("creating file: %w", err)}
			return
		}

		n, err := io.Copy(f, io.LimitReader(reader, size))
		f.Close()
		if err != nil {
			os.Remove(path)
			events <- ReceiveEvent{Err: fmt.Errorf("receiving data: %w", err)}
			return
		}

		// Verify integrity
		received, err := os.ReadFile(path)
		if err != nil {
			events <- ReceiveEvent{Err: fmt.Errorf("reading received file: %w", err)}
			return
		}
		computedHash := pcrypto.HashBytes(received)
		if hex.EncodeToString(computedHash) != hex.EncodeToString(hashBuf) {
			os.Remove(path)
			events <- ReceiveEvent{Err: fmt.Errorf("integrity check failed for %s", filename)}
			return
		}

		events <- ReceiveEvent{
			Filename: filename,
			Size:     n,
			From:     remotePeer,
		}
	})

	// Graceful shutdown on signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	stopFn := func() {
		close(done)
		h.Close()
		close(events)
	}

	go func() {
		<-sigCh
		stopFn()
	}()

	return &ListenResult{
		Events: events,
		Stop:   stopFn,
	}, nil
}
