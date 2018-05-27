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
	"context"
	"math/big"

	"github.com/timenewbank/go-mit/accounts"
	"github.com/timenewbank/go-mit/common"
	"github.com/timenewbank/go-mit/common/math"
	"github.com/timenewbank/go-mit/core"
	"github.com/timenewbank/go-mit/core/bloombits"
	"github.com/timenewbank/go-mit/core/state"
	"github.com/timenewbank/go-mit/core/types"
	"github.com/timenewbank/go-mit/core/vm"
	"github.com/timenewbank/go-mit/mit/downloader"
	"github.com/timenewbank/go-mit/mit/gasprice"
	"github.com/timenewbank/go-mit/mitdb"
	"github.com/timenewbank/go-mit/event"
	"github.com/timenewbank/go-mit/params"
	"github.com/timenewbank/go-mit/rpc"
)

// EthApiBackend implements mitapi.Backend for full nodes
type MitApiBackend struct {
	mit *Mitereum
	gpo *gasprice.Oracle
}

func (b *MitApiBackend) ChainConfig() *params.ChainConfig {
	return b.mit.chainConfig
}

func (b *MitApiBackend) CurrentBlock() *types.Block {
	return b.mit.blockchain.CurrentBlock()
}

func (b *MitApiBackend) SetHead(number uint64) {
	b.mit.protocolManager.downloader.Cancel()
	b.mit.blockchain.SetHead(number)
}

func (b *MitApiBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.mit.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.mit.blockchain.CurrentBlock().Header(), nil
	}
	return b.mit.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *MitApiBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.mit.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.mit.blockchain.CurrentBlock(), nil
	}
	return b.mit.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *MitApiBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.mit.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.mit.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *MitApiBackend) GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	return b.mit.blockchain.GetBlockByHash(blockHash), nil
}

func (b *MitApiBackend) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error) {
	return core.GetBlockReceipts(b.mit.chainDb, blockHash, core.GetBlockNumber(b.mit.chainDb, blockHash)), nil
}

func (b *MitApiBackend) GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error) {
	receipts := core.GetBlockReceipts(b.mit.chainDb, blockHash, core.GetBlockNumber(b.mit.chainDb, blockHash))
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *MitApiBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.mit.blockchain.GetTdByHash(blockHash)
}

func (b *MitApiBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.mit.BlockChain(), nil)
	return vm.NewEVM(context, state, b.mit.chainConfig, vmCfg), vmError, nil
}

func (b *MitApiBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.mit.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *MitApiBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.mit.BlockChain().SubscribeChainEvent(ch)
}

func (b *MitApiBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.mit.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *MitApiBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.mit.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *MitApiBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.mit.BlockChain().SubscribeLogsEvent(ch)
}

func (b *MitApiBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.mit.txPool.AddLocal(signedTx)
}

func (b *MitApiBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.mit.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *MitApiBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.mit.txPool.Get(hash)
}

func (b *MitApiBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.mit.txPool.State().GetNonce(addr), nil
}

func (b *MitApiBackend) Stats() (pending int, queued int) {
	return b.mit.txPool.Stats()
}

func (b *MitApiBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.mit.TxPool().Content()
}

func (b *MitApiBackend) SubscribeTxPreEvent(ch chan<- core.TxPreEvent) event.Subscription {
	return b.mit.TxPool().SubscribeTxPreEvent(ch)
}

func (b *MitApiBackend) Downloader() *downloader.Downloader {
	return b.mit.Downloader()
}

func (b *MitApiBackend) ProtocolVersion() int {
	return b.mit.EthVersion()
}

func (b *MitApiBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *MitApiBackend) ChainDb() mitdb.Database {
	return b.mit.ChainDb()
}

func (b *MitApiBackend) EventMux() *event.TypeMux {
	return b.mit.EventMux()
}

func (b *MitApiBackend) AccountManager() *accounts.Manager {
	return b.mit.AccountManager()
}

func (b *MitApiBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.mit.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *MitApiBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.mit.bloomRequests)
	}
}
