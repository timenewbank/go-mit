package main

import (
	x509 "github.com/timenewbank/go-mit/mitx509"
	"path/filepath"
	"log"
	"fmt"
	"encoding/hex"
	"flag"
)

func main() {

	certPath:=flag.String( "F","","input the certPath")
	rootPath:=flag.String("RF","","input the rootCertPath")
	flag.Parse()

	//dir info
	if *certPath==""||*rootPath==""{
		log.Panic("no filepath")
	}


	crtString:=x509.ReadFile(filepath.Join(*certPath))
	crtA,_:=x509.ParseCrtString(crtString)
	//fmt.Println("name====>",crtA.Subject.CommonName)

	//handle the root path
	rootString:=x509.ReadFile(filepath.Join(*rootPath))
	rootA,_:=x509.ParseCrtString(rootString)
	error:=crtA.CheckSignatureFrom(rootA)
	if error!=nil{
		fmt.Println("error",error)
	}else{
		bytes:=[]byte(crtString)
		strHex:=hex.EncodeToString(bytes)
		fmt.Println("0x"+strHex)
	}





}


