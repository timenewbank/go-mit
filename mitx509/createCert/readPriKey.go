package main

import(
	"flag"
	"log"
	x509 "github.com/timenewbank/go-mit/mitx509"
	"path/filepath"
	"fmt"
	"encoding/hex"
)

func main() {
	keyPath:=flag.String( "KP","","input the keyPath")
	flag.Parse()

	//dir info
	if *keyPath==""{
		log.Panic("no filepath")
	}

	keyString:=x509.ReadFile(filepath.Join(*keyPath))
	//fmt.Println("key===>",keyString)
	bytes:=[]byte(keyString)
	strHex:=hex.EncodeToString(bytes)
	fmt.Println("0x"+strHex)
}