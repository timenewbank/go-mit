// Copyright 2018 The go-mit Authors
// This file is part of the go-mit library.
//
// The go-mit library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-mit library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-mit library. If not, see <http://www.gnu.org/licenses/>.

package mitx509

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"crypto/rsa"
	"crypto/rand"
	"math/big"
	"log"
	"time"
	"fmt"
)


type X509mit struct {
	pub_key []byte
	priv_key []byte
}

func (x *X509mit) GenCa(country []string, organization []string, organizationalUnit []string) *X509mit{
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2018),
		Subject:pkix.Name{
			Country: country,
			Organization: organization,
			OrganizationalUnit: organizationalUnit,
		},
		NotBefore:	time.Now(),
		NotAfter:	time.Now().AddDate(10,0,0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		SubjectKeyId: []byte{1,2,3,4,5},
		KeyUsage: x509.KeyUsageDigitalSignature|x509.KeyUsageCertSign,
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 1024)
	pub := &priv.PublicKey
	pub_key, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)

	if err != nil {
		log.Println("create ca failed", err)
		return nil
	}
	priv_key := x509.MarshalPKCS1PrivateKey(priv)
	x509 := &X509mit{pub_key,priv_key}
	fmt.Printf("x509 %v",*x509)
	return x509
}

func (x *X509mit) ReadCa([]byte){
	return
}

func (x *X509mit) VerifyCa([]byte){
	return
}

func (x *X509mit) Is([]byte){
	return
}