package fn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
)

func FnPost(protocolPath string, filePath string, idFilePath string) string {
	
	ctx := context.Background()
	

	var cfg Protocol	
	data,_:=os.ReadFile(protocolPath)
	json.Unmarshal(data, &cfg)

	var idf IDFile
	dataID,_:=os.ReadFile(idFilePath)
	json.Unmarshal(dataID,&idf)


	h, _:= libp2p.New()
	fmt.Printf("Connexion %s", h.ID().String())

	maddr, _ := ma.NewMultiaddr(cfg.RelayAddr)
	ri, _ := peer.AddrInfoFromP2pAddr(maddr)
	h.Connect(ctx, *ri)
	fmt.Println("✅ Connecté au relay")

 	s,_:=	h.NewStream(ctx, ri.ID,protocol.ID(cfg.Protocol))
 	defer s.Close()

	f,_:=os.Open(filePath)

	_,err:=io.Copy(s, f)
	if err!=nil{
		fmt.Printf("error sending")
	}

	return fmt.Sprintf("%s a été envoyé via le protocol %s",filePath,cfg.Groupname)
}