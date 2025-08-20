package fn

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	lcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p"
)

func FnGet(protocolPath, storeDir, privPath string)  {
	ctx := context.Background()

	//insère la clé privée générer à partir du peerid

	data:=unmarshalPriv(privPath)
	raw,_:=base64.StdEncoding.DecodeString(data.Priv)
	priv,_:=lcrypto.UnmarshalPrivateKey(raw)
	


	cfg:=unmarshalProtocol(protocolPath)
	keep:=keep(cfg)

	os.MkdirAll(storeDir, 0o755)


	h, _ := libp2p.New(libp2p.Identity(priv))
	if data.PeerId == h.ID().String() {
	fmt.Println("📡 PeerID is ok :", h.ID().String())
	}
	

	maddr, _ := ma.NewMultiaddr(cfg.RelayAddr)
	ri, _ := peer.AddrInfoFromP2pAddr(maddr)
	h.Connect(ctx, *ri)
	fmt.Println("✅ Connecté au relay")

	 handler := func (s network.Stream) {
		defer s.Close()
		
		reader:=bufio.NewReader(s)
		filename,_:=reader.ReadString('\n')
		filename=strings.TrimSpace(filename)
		filename=filepath.Base(filename)


		path:=filepath.Join(storeDir,filename)
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
				exec.Command("taskkill", "/PID", fmt.Sprint(os.Getpid()), "/F").Run()
				return
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

func keep(cfg Protocol) string{
		base := `c:\pulse_test\receivers`
		os.MkdirAll(base, 0o755)
		sum:=sha256.Sum256([]byte(cfg.Protocol))
		name:=hex.EncodeToString(sum[:8])
		keep := filepath.Join(base, name+".keep")
		os.WriteFile(keep, []byte("1"), 0o600)
		return keep
	}