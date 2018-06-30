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
	"testing"
	"fmt"
)

var (
	country = []string{"Singapore"}
	organization = []string{"M.I.T Foundation"}
	organizationalUnit = []string{"0xbb85e976aaaf00f647a8ca0f5d8fa8583bb8d82e"}
	pub_key []byte
	priv_key []byte
)

func TestGenCa(t *testing.T) {
	x509mit := X509mit{pub_key,priv_key}
	x509 := x509mit.GenCa(country,organization,organizationalUnit)
	fmt.Println("Test X509 Infomation:",x509)
}




