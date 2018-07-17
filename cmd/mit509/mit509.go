// Copyright 2018 The go-mit Authors
// This file is part of go-mit.
//
// go-mit is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-mit is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-mit. If not, see <http://www.gnu.org/licenses/>.

// mit is the official command-line client for Mit.

package main

import (
	"github.com/timenewbank/go-mit/mitx509"
	"fmt"
)

var (
	country = []string{"Singapore"}
	orgnazation = []string{"M.I.T&TNB Foundation"};
	orgnazationalUnit = []string{"TNB Team"}
)

func main(){
	x509_issuer := &mitx509.X509mit{}
	x509_issuer = x509_issuer.GenCa(country,orgnazation,orgnazationalUnit)
	fmt.Printf("%v",x509_issuer)
	return;
}