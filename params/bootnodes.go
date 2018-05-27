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
	"enode://ba5d98d5f5fdc4375eda94ea2665684a25c38ae6f4678d27aa9d1e40cc37f1b1f20d09f7c339facbfffa47b59dc0b9a59743df3710d94dfe181692b77e8171b8@103.211.167.5:30303",
	"enode://cd11029b9eff87c4a81e5b56550cb41181621022a0b11ac2a623e079322cea27d8c016078ef5effdefcf357e0bb50f0f1436200241f1b1cefff416742aba9605@103.242.67.65:30303",

}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Mit test network.
var TestnetBootnodes = []string{
	"enode://ba5d98d5f5fdc4375eda94ea2665684a25c38ae6f4678d27aa9d1e40cc37f1b1f20d09f7c339facbfffa47b59dc0b9a59743df3710d94dfe181692b77e8171b8@103.211.167.5:30303",
	"enode://cd11029b9eff87c4a81e5b56550cb41181621022a0b11ac2a623e079322cea27d8c016078ef5effdefcf357e0bb50f0f1436200241f1b1cefff416742aba9605@103.242.67.65:30303",
}

var RinkebyBootnodes = []string{
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
}
