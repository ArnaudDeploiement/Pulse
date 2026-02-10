package transport

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	ma "github.com/multiformats/go-multiaddr"
)

// RelayInfo holds the running relay details.
type RelayInfo struct {
	PeerID string
	Addrs  []string
}

// StartRelay starts a libp2p host configured as a circuit relay v2 server.
// It blocks until SIGINT/SIGTERM is received.
func StartRelay(ctx context.Context, port int) (*RelayInfo, <-chan struct{}, error) {
	// Generate a dedicated key for the relay
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generating relay key: %w", err)
	}

	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)
	listen6Addr := fmt.Sprintf("/ip6/::/tcp/%d", port)

	h, err := libp2p.New(
		libp2p.Identity(priv),
		libp2p.ListenAddrs(
			ma.StringCast(listenAddr),
			ma.StringCast(listen6Addr),
		),
		libp2p.EnableRelayService(),
		libp2p.ForceReachabilityPublic(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("creating relay host: %w", err)
	}

	info := &RelayInfo{
		PeerID: h.ID().String(),
	}

	for _, addr := range h.Addrs() {
		full := addr.Encapsulate(ma.StringCast("/p2p/" + h.ID().String()))
		info.Addrs = append(info.Addrs, full.String())
	}

	done := make(chan struct{})
	go waitForShutdown(h, done)

	return info, done, nil
}

func waitForShutdown(h host.Host, done chan struct{}) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	h.Close()
	close(done)
}
