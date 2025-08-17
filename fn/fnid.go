package fn

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)


func AddPeerId(peerID []string) string {

	basedir:=`C:/pulse_test/IDFile`
	os.MkdirAll(basedir, 0o755); 


	idf := IDFile{
		PeerId: peerID,
	}

	data,_:=json.MarshalIndent(idf," ", " ")

	name:=make([]byte,8)
	rand.Read(name)
	file:=base64.RawURLEncoding.EncodeToString(name)+".json"
	
	os.WriteFile(filepath.Join(basedir,file),data,0o755)


	path:=filepath.Join(basedir,file)
	fmt.Printf("Your file has been created : %s\n", path)

	return path 

}