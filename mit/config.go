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

package mit

import (
	"math/big"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"github.com/timenewbank/go-mit/common"
	"github.com/timenewbank/go-mit/common/hexutil"
	"github.com/timenewbank/go-mit/consensus/mithash"
	"github.com/timenewbank/go-mit/core"
	"github.com/timenewbank/go-mit/mit/downloader"
	"github.com/timenewbank/go-mit/mit/gasprice"
	"github.com/timenewbank/go-mit/params"
)

// DefaultConfig contains default settings for use on the Mit main net.
var DefaultConfig = Config{
	//SyncMode: downloader.FastSync,
	SyncMode: downloader.FullSync,
	Mithash: mithash.Config{
		CacheDir:       "mithash",
		CachesInMem:    2,
		CachesOnDisk:   3,
		DatasetsInMem:  1,
		DatasetsOnDisk: 2,
	},
	NetworkId:     666,
	LightPeers:    100,
	DatabaseCache: 768,
	TrieCache:     256,
	TrieTimeout:   5 * time.Minute,
	GasPrice:      big.NewInt(1800 * params.Shannon),//18000-1800

	TxPool: core.DefaultTxPoolConfig,
	GPO: gasprice.Config{
		Blocks:     20,
		Percentile: 60,
	},
}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}
	if runtime.GOOS == "windows" {
		DefaultConfig.Mithash.DatasetDir = filepath.Join(home, "AppData", "Mithash")
	} else {
		DefaultConfig.Mithash.DatasetDir = filepath.Join(home, ".mithash")
	}
}

//go:generate gencodec -type Config -field-override configMarshaling -formats toml -out gen_config.go

type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, the Mit main net block is used.
	Genesis *core.Genesis `toml:",omitempty"`

	// Protocol options
	NetworkId uint64 // Network ID to use for selecting peers to connect to
	SyncMode  downloader.SyncMode
	NoPruning bool

	// Light client options
	LightServ  int `toml:",omitempty"` // Maximum percentage of time allowed for serving LES requests
	LightPeers int `toml:",omitempty"` // Maximum number of LES client peers

	// Database options
	SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int
	TrieCache          int
	TrieTimeout        time.Duration

	// Mining-related options
	Etherbase    common.Address `toml:",omitempty"`
	MinerThreads int            `toml:",omitempty"`
	ExtraData    []byte         `toml:",omitempty"`
	GasPrice     *big.Int

	// Mithash options
	Mithash mithash.Config

	// Transaction pool options
	TxPool core.TxPoolConfig

	// Gas Price Oracle options
	GPO gasprice.Config

	// Enables tracking of SHA3 preimages in the VM
	EnablePreimageRecording bool

	// Miscellaneous options
	DocRoot string `toml:"-"`
}

type configMarshaling struct {
	ExtraData hexutil.Bytes
}
