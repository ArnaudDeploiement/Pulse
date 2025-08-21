package fn

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
)

func FnPost(protocolPath, filePath, idFilePath string) string {
	ctx := context.Background()

	// Config protocole
	var cfg Protocol
	b, _ := os.ReadFile(protocolPath)
	_ = json.Unmarshal(b, &cfg)

	// Liste des destinataires
	var idf IDFile
	bi, _ := os.ReadFile(idFilePath)
	_ = json.Unmarshal(bi, &idf)

	// Host émetteur avec support du relay
	h, _ := libp2p.New(
		libp2p.EnableRelay(), // nécessaire pour dialer /p2p-circuit
	)

	// Connexion au relay (une seule fois)
	relayMA, _ := ma.NewMultiaddr(cfg.RelayAddr)
	relayInfo, _ := peer.AddrInfoFromP2pAddr(relayMA)
	_ = h.Connect(ctx, *relayInfo)
	fmt.Println("✅ Connecté au relay:", relayInfo.ID)

	// Préparer le payload
	payload, _ := os.ReadFile(filePath)
	name := filepath.Base(filePath)

	// Envoi à chaque PeerID
	for _, pid := range idf.PeerId {
		// Adresse circuit complète: relay + /p2p-circuit + dest
		full := fmt.Sprintf("%s/p2p-circuit/p2p/%s", cfg.RelayAddr, pid)
		maddr, _ := ma.NewMultiaddr(full)
		destInfo, _ := peer.AddrInfoFromP2pAddr(maddr)

		// Donner l'adresse circuit au peerstore et forcer un dial via le relay
		h.Peerstore().AddAddrs(destInfo.ID, destInfo.Addrs, time.Minute)
		_ = h.Connect(ctx, *destInfo)

		// Attendre brièvement que l'état passe à Connected
		for i := 0; i < 10; i++ {
			if h.Network().Connectedness(destInfo.ID) == network.Connected {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}

		// Ouvrir un stream VERS le destinataire (pas le relay)
		s, _ := h.NewStream(ctx, destInfo.ID, protocol.ID(cfg.Protocol))
		if s == nil {
			fmt.Println("stream=nil → circuit non établi pour", pid)
			continue
		}

		// Envoyer: nom de fichier (ligne) puis le contenu
		w := bufio.NewWriter(s)
		fmt.Fprintln(w, name)
		w.Flush()

		s.Write(payload)
		s.Close()

		fmt.Println("📤 sent to:", pid)
	}

	return fmt.Sprintf("✅ %s envoyé via le protocole %s", filePath, cfg.Groupname)
}
