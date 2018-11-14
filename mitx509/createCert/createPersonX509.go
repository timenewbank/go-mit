package main

import (
	"path/filepath"
	"os"
	"log"
	x509 "github.com/timenewbank/go-mit/mitx509"
	"flag"
	"github.com/timenewbank/go-mit/node"
)

func main() {
	getCountry:=flag.String("C","","input the Country")
	getOrganization:=flag.String("OR","","input the Organization")
	getOrganizationalUnit:=flag.String("OU","","input the OrganizationalUnit")
	getEmailAddress:=flag.String("EA","","input the EmailAddress")
	getLocality:=flag.String("L","","input the Locality")
	getProvince:=flag.String("P","","input the Province")
	getCommonName:=flag.String("CN","","input the CommonName")
	getCrtName:=flag.String("CrtN","","input the CrtName")
	getKeyName:=flag.String("KeyN","","input the KeyName")
	getRootPath:=flag.String("R","","input the RootCertPath")
	getRootKeyPath:=flag.String("RK","","input RootKeyPath")
	getCABool:=flag.Bool("IsCA",false,"isCA")

	flag.Parse()

	if flag.Parsed(){
		if *getCountry==""{
			log.Panic("input person country")
		}
		if *getOrganization==""{
			log.Panic("input person organization")
		}
		if *getOrganizationalUnit==""{
			log.Panic("input person organizationalUnit")
		}
		if *getEmailAddress==""{
			log.Panic("input person emailAddress")
		}
		if *getRootPath==""||*getRootKeyPath==""{
			log.Panic("rootPath is incorrent")
		}
	}

	//dir info
	path:= node.DefaultDataDir()
	if path != "" {
		path=filepath.Join(path,x509.CertPath)
	}


	certPath:=path
	flag,error:=x509.PathExists(certPath)
	if error!=nil{
		log.Panic("error find exist")
	}
	if !flag{
		//there is no filepath create it
		err:=os.MkdirAll(certPath,os.ModePerm)
		if err!=nil{
			log.Println("create dir fail")
		}
	}


	//root info
	personInfo := x509.CertInfo{
		Country: []string{*getCountry},
		Organization: []string{*getOrganization},
		IsCA: *getCABool,
		OrganizationalUnit: []string{*getOrganizationalUnit},
		EmailAddress: []string{*getEmailAddress},
		Locality: []string{*getLocality},
		Province: []string{*getProvince},
		CommonName: *getCommonName,
		CrtName: filepath.Join(certPath,*getCrtName),
		KeyName: filepath.Join(certPath,*getKeyName)}

	crtBool:=x509.FileExists(personInfo.CrtName)
	keyBool:=x509.FileExists(personInfo.KeyName)

	//get the root cert
	rootCrt, rootPri, err := x509.Parse(*getRootPath, *getRootKeyPath)
	if err!=nil{
		log.Panic("parse root error")
	}

	if !crtBool&&!keyBool {
		error=x509.CreatePersonCERT(rootCrt,rootPri,personInfo)
		if error!=nil{
			log.Println("create person cert fail")
		}else{
			log.Println("crate person cert success")
		}
	}else{
		log.Println("there is exist cert or key")
	}


}
