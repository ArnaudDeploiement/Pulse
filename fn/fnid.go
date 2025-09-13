package fn

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/crypto"
	lcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

func AddPeerId(peerID []string) string {
	basedir := filepath.Join(baseDir(), "IDFile")
	os.MkdirAll(basedir, 0o755)

	idf := IDFile{
		PeerId: peerID,
	}

	data, _ := json.MarshalIndent(idf, " ", " ")

	name := make([]byte, 8)
	rand.Read(name)
	file := base64.RawURLEncoding.EncodeToString(name) + "_GroupID.json"

	os.WriteFile(filepath.Join(basedir, file), data, 0o755)

	path := filepath.Join(basedir, file)
	fmt.Printf("Your file has been created : %s\n", path)

	return path
}

func Getid() string {
	basedir := filepath.Join(baseDir(), "IDFile")
	os.MkdirAll(basedir, 0o755)

	privkey, _, _ := crypto.GenerateEd25519Key(rand.Reader)
	id, _ := peer.IDFromPrivateKey(privkey)
	raw, _ := lcrypto.MarshalPrivateKey(privkey)

	peerid := IdPeer{
		PeerId: id.String(),
		Priv:   base64.StdEncoding.EncodeToString(raw),
	}

	data, _ := json.MarshalIndent(peerid, " ", " ")

	name := make([]byte, 8)
	rand.Read(name)
	file := base64.RawURLEncoding.EncodeToString(name) + "_PeerID.json"

	os.WriteFile(filepath.Join(basedir, file), data, 0o755)

	path := filepath.Join(basedir, file)
	fmt.Printf("Your file has been created : %s\n", path)

	return id.String()
}
