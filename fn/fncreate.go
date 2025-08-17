package fn

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)



func FnCreate(groupname string, relayAddr string, outdir string) string {
	basedir:=`C:/pulse_test`
	os.MkdirAll(basedir, 0o755); 
	
	if outdir == ""{ outdir=`.`}
	os.MkdirAll(outdir,0o755)
	
	protocol:=Protocol{
		Groupname: groupname,
		Protocol: endpointProtocol() ,
		RelayAddr: relayAddr,
		As:as(32),
	
	}

	dataProtocol,_:=json.MarshalIndent(protocol,"","  ")
	os.WriteFile(filepath.Join(basedir,"Protocol_"+groupname+".json"),dataProtocol,0600)
	
	return fmt.Sprintf("\n🌐 Protocol_%s has been created with success in %s\n📂 A repository directory has been created at %s\n\n🔥 Congratulations, you’ve created an endpoint for exchanging data\n\n🌟 Your protocol %s is now a standalone API! Welcome to Web3! :)\n\n", groupname,basedir,outdir, groupname)
}

func endpointProtocol() string{
	protoID:=make([]byte,16)
	rand.Read(protoID)
	return "/pulse/"+base64.RawURLEncoding.EncodeToString(protoID)+"/1.0"
}

func as(n int) string {
	b:=make([]byte,n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}