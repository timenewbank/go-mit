package mitx509

import (
	"crypto/x509/pkix"
	"crypto/x509"
	"crypto/rsa"
	"math/big"
	"time"
	"crypto/rand"
	rd "math/rand"
	"os"
	"encoding/pem"
	"io/ioutil"
	"errors"
)

const (
	X509Version=1.0
	CertPath="cert"
	RootCrtName="tnb_root.crt"
	RootKeyName="tnb_root.key"
)

//struct for x509
type CertInfo struct {
	Country            	[]string
	Organization       	[]string
	OrganizationalUnit 	[]string
	EmailAddress       	[]string
	Province           	[]string
	Locality           	[]string
	CommonName         	string
	CrtName				string
	KeyName   			string
	IsCA               	bool
	Names              	[]pkix.AttributeTypeAndValue
}

//rootCert only for M.I.T
func CreateRootCERT(info CertInfo) error {
	Crt := newCertificate(info)
	Key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	var buf []byte
	buf, err = x509.CreateCertificate(rand.Reader, Crt, Crt, &Key.PublicKey, Key)

	if err != nil {
		return err
	}

	err = write(info.CrtName, "CERTIFICATE", buf)
	if err != nil {
		return err
	}

	buf = x509.MarshalPKCS1PrivateKey(Key)
	return write(info.KeyName, "PRIVATE KEY", buf)
}


func CreatePersonCERT(RootCa *x509.Certificate, RootKey *rsa.PrivateKey, info CertInfo) error {
	Crt := newCertificate(info)
	Key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	var buf []byte
	if RootCa == nil || RootKey == nil {
		//create self
		return errors.New("not found the root cert")
	} else {
		//use CA cert
		buf, err = x509.CreateCertificate(rand.Reader, Crt, RootCa, &Key.PublicKey, RootKey)
	}
	if err != nil {
		return err
	}
	var keyBuf []byte
	err = write(info.CrtName, "CERTIFICATE", buf)
	if err != nil {
		return err
	}

	keyBuf = x509.MarshalPKCS1PrivateKey(Key)
	return write(info.KeyName, "PRIVATE KEY", keyBuf)
}

func newCertificate(info CertInfo) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(rd.Int63()),
		Subject: pkix.Name{
			Country:            info.Country,
			Organization:       info.Organization,
			OrganizationalUnit: info.OrganizationalUnit,
			Province:           info.Province,
			CommonName:         info.CommonName,
			Locality:           info.Locality,
			ExtraNames:         info.Names,
		},
		NotBefore:             time.Now(),//start time
		NotAfter:              time.Now().AddDate(10, 0, 0),//10 years
		BasicConstraintsValid: true, //
		IsCA:           info.IsCA,   //is CA
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},    //
		KeyUsage:       x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		EmailAddresses: info.EmailAddress,
		Version:		X509Version,
	}
}


//
func write(filename, Type string, p []byte) error {
	File, err := os.Create(filename)
	defer File.Close()
	if err != nil {
		return err
	}
	var b *pem.Block = &pem.Block{Bytes: p, Type: Type}
	return pem.Encode(File, b)
}


func Parse(crtPath, keyPath string) (rootcertificate *x509.Certificate, rootPrivateKey *rsa.PrivateKey, err error) {
	rootcertificate, err = ParseCrt(crtPath)
	if err != nil {
		return
	}
	rootPrivateKey, err = ParseKey(keyPath)
	return
}

func ParseCrt(path string) (*x509.Certificate, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	return x509.ParseCertificate(p.Bytes)
}


func ParseCrtString(str string) (*x509.Certificate, error) {
	p := &pem.Block{}
	p,_= pem.Decode([]byte(str))
	return x509.ParseCertificate(p.Bytes)
}

func ParseKey(path string) (*rsa.PrivateKey, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p, buf := pem.Decode(buf)
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}



// path is exist
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}



func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}



func ReadFile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}


func VerifyCrt(){

}
