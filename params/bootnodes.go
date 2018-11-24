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

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Mit network.
var MainnetBootnodes = []string{
	//the main bootnode
	"enode://2a1fa5fb68c68536e0ed9c20eabd64b749264bcb12be24130861139e02cbd6aeb2dadec34b201e803586c042c5807e82951b0d7a14c9cafe9626a7fca68f6183@106.2.11.139:9999",
	"enode://880da7bfc175239d34d1a5b576e9cd54f24559ef4e9b35cd275942707fa77e0e109e5cc43f24e61e93fe8c392cdf9660389c5877c139183a836343d57b5b7379@103.242.67.59:9999",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Mit test network.
var TestnetBootnodes = []string{
	"enode://819656a6f60cda49fab06e04ac47f148fcb3a48cefec41e8e61795366c3771a514b1f267d62ae36230a96e665995530e34d639b0d27e1ff33ba5fa89e83a7549@103.242.67.65:9999",
	"enode://bbf4ba9aeea72dd1235d83828cc21a498db9dd2f4c2de93da9f607fffb2ee9bba540606ff0b576212cc0aef7024ce92bed0f38b153be65ba81bfdf21a2a3f0ef@103.242.67.72:9999",
	"enode://9ba2930e537939a981a5d908d96b45d8a7fb2c7364a3d0a7dc57686fa3c7d2ff5392eb015043f03e4e9d5a5b7f6e36cebf961fdd1072fc06473e6a11c779c6d9@103.242.67.72:9998",
}

var RinkebyBootnodes = []string{
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
}
