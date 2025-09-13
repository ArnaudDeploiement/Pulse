package fn

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	lcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	rclient "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client" // ← AJOUT
	ma "github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p"
)

func FnGet(protocolPath, storeDir, privPath string) {
	ctx := context.Background()

	data := unmarshalPriv(privPath)
	raw, _ := base64.StdEncoding.DecodeString(data.Priv)
	priv, _ := lcrypto.UnmarshalPrivateKey(raw)

	cfg := unmarshalProtocol(protocolPath)
	keep := keep(cfg)

	os.MkdirAll(storeDir, 0o755)

	h, _ := libp2p.New(libp2p.Identity(priv), libp2p.EnableRelay(), libp2p.EnableRelayService())
	if data.PeerId == h.ID().String() {
		fmt.Println("📡 PeerID is ok :", h.ID().String())
	}

	maddr, _ := ma.NewMultiaddr(cfg.RelayAddr)
	ri, _ := peer.AddrInfoFromP2pAddr(maddr)
	h.Connect(ctx, *ri)
	fmt.Println("✅ Connecté au relay")

	res, err := rclient.Reserve(ctx, h, *ri)
	if err == nil {
		fmt.Println("🎫 Réservation OK, expires in:", res.Expiration) // TTL ~ 2min
	} else {
		fmt.Println("❌ Reservation échouée:", err)
	}

	handler := func(s network.Stream) {
		defer s.Close()

		reader := bufio.NewReader(s)
		filename, _ := reader.ReadString('\n')
		filename = strings.TrimSpace(filename)
		filename = filepath.Base(filename)

		path := filepath.Join(storeDir, filename)
		f, _ := os.Create(path)
		defer f.Close()
		io.Copy(f, reader)

		fmt.Println("📥 Reçu →", path)
	}

	h.SetStreamHandler(protocol.ID(cfg.Protocol), handler)

	go func() {
		t := time.NewTicker(2 * time.Second)
		for range t.C {
			if _, err := os.Stat(keep); os.IsNotExist(err) {
				fmt.Println("🛑 Stop détecté")
				h.Close()
				os.Exit(0)
			}
		}
	}()

	fmt.Println("👂 En écoute sur", cfg.Protocol, "→ dépôt :", storeDir)
	select {}
}

func unmarshalProtocol(protocolPath string) Protocol {
	data, _ := os.ReadFile(protocolPath)
	var cfg Protocol
	json.Unmarshal(data, &cfg)
	return cfg
}
func unmarshalPriv(privPath string) IdPeer {
	data, _ := os.ReadFile(privPath)
	var privcfg IdPeer
	json.Unmarshal(data, &privcfg)
	return privcfg
}

func keep(cfg Protocol) string {
	base := receiversDir()
	name := sanitize(cfg.Groupname)
	keep := filepath.Join(base, name+".keep")
	os.WriteFile(keep, []byte(cfg.Protocol), 0o600)
	return keep
}
