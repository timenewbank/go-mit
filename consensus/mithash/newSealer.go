package mithash

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"sync"

	"github.com/timenewbank/go-mit/common"
	"github.com/timenewbank/go-mit/consensus"
	"github.com/timenewbank/go-mit/core/types"
	"github.com/timenewbank/go-mit/log"
	"encoding/binary"
	"github.com/timenewbank/go-mit/crypto"
	"github.com/timenewbank/go-mit/core/state"
)




// Seal implements consensus.Engine, attempting to find a nonce that satisfies
// the block's difficulty requirements.
func (mithash *Mithash) NewSeal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{},stateDb *state.StateDB) (*types.Block, error) {
	// If we're running a fake PoW, simply return a 0 nonce immediately
	if mithash.config.PowMode == ModeFake || mithash.config.PowMode == ModeFullFake {
		header := block.Header()
		header.Nonce, header.MixDigest = types.BlockNonce{}, common.Hash{}
		return block.WithSeal(header), nil
	}
	// If we're running a shared PoW, delegate sealing to it
	if mithash.shared != nil {
		return mithash.shared.Seal(chain, block, stop)
	}
	//get the stateDB and get the balance of coinbase
	currentHeader:=chain.CurrentHeader()
	//delay the header balance
	//if block.Number().Cmp(common.PoSDis)<=0{
	//	currentHeader=chain.GetHeaderByNumber(common.Big0.Uint64())
	//}else {
	//	currentHeader=chain.GetHeaderByNumber(new(big.Int).Sub(block.Number(),common.PoSDis).Uint64())
	//}

	newStateDb,_:=state.New(currentHeader.Root,stateDb.Database())
	//newStateDb.Reset(currentHeader.Root)
	balance:=newStateDb.GetBalance(block.Header().Coinbase)

	//balance:=mit.ChainStateReader().BalanceAt()


	// Create a runner and the multiple search threads it directs
	abort := make(chan struct{})
	found := make(chan *types.Block)

	mithash.lock.Lock()
	threads := mithash.threads
	if mithash.rand == nil {
		seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			mithash.lock.Unlock()
			return nil, err
		}
		mithash.rand = rand.New(rand.NewSource(seed.Int64()))
	}
	mithash.lock.Unlock()
	if threads == 0 {
		threads = runtime.NumCPU()
	}
	if threads < 0 {
		threads = 0 // Allows disabling local mining without extra logic around local/remote
	}

	var pend sync.WaitGroup
	for i := 0; i < threads; i++ {
		pend.Add(1)
		go func(id int, nonce uint64) {
			defer pend.Done()
			mithash.newMine(block, id, nonce, abort, found,balance)
		}(i, uint64(mithash.rand.Int63()))
	}
	// Wait until sealing is terminated or a nonce is found
	var result *types.Block
	select {
	case <-stop:
		// Outside abort, stop all miner threads
		close(abort)
	case result = <-found:
		// One of the threads found a block, abort all others
		close(abort)
	case <-mithash.update:
		// Thread count was changed on user request, restart
		close(abort)
		pend.Wait()
		return mithash.Seal(chain, block, stop)
	}
	// Wait for all miners to terminate and return the block
	pend.Wait()
	return result, nil
}

// mine is the actual proof-of-work miner that searches for a nonce starting from
// seed that results in correct final block difficulty.
func (mithash *Mithash) newMine(block *types.Block, id int, seed uint64, abort chan struct{}, found chan *types.Block,balance *big.Int) {
	balanceTarget:=new(big.Int).Mul(balance,big.NewInt(1))
	//if the balance==0
	//fmt.Println("coinbase===balanceValue===>",block.Header().Coinbase,balanceTarget)
	balanceTarget=balanceTarget.Div(balanceTarget,big.NewInt(1000000000000000000))
	if balanceTarget.Cmp(common.Big0)<=0{
		balanceTarget=common.Big1
	}
	// Extract some data from the header
	var (
		header  = block.Header()
		hash    = header.HashNoNonce().Bytes()
		target  = new(big.Int).Div(maxUint256, header.Difficulty)
		posTarget=new(big.Int).Div(new(big.Int).Mul(maxUint256,balanceTarget),header.Difficulty)
		number  = header.Number.Uint64()
		dataset = mithash.dataset(number)
	)
	//
	//fmt.Println("coinbase===balance===>",balanceTarget)

	// Start generating random nonces until we abort or find a good one
	var (
		attempts = int64(0)
		nonce    = seed
	)
	logger := log.New("miner", id)
	logger.Trace("Started mithash search for new nonces", "seed", seed)
search:
	for {
		select {
		case <-abort:
			// Mining terminated, update stats and abort
			logger.Trace("Mithash nonce search aborted", "attempts", nonce-seed)
			mithash.hashrate.Mark(attempts)
			break search

		default:
			// We don't have to update hash rate on every nonce, so update after after 2^X nonces
			attempts++
			if (attempts % (1 << 15)) == 0 {
				mithash.hashrate.Mark(attempts)
				attempts = 0
			}
			// Compute the PoW value of this nonce
			digest, result := hashimotoFull(dataset.dataset, hash, nonce)
			//Add PoS calculation result
			posResult:=posMine(nonce,hash,header.Time,header.ParentHash,header.UncleHash)
			if new(big.Int).SetBytes(result).Cmp(target) <= 0 &&new(big.Int).SetBytes(posResult).Cmp(posTarget) <= 0{
				// Correct nonce found, create a new header with it
				header = types.CopyHeader(header)
				header.Nonce = types.EncodeNonce(nonce)
				header.MixDigest = common.BytesToHash(digest)

				// Seal and return a block (if still needed)
				select {
				case found <- block.WithSeal(header):
					//fmt.Println("new pos+pow get the answer====》")
					logger.Trace("Mithash nonce found and reported", "attempts", nonce-seed, "nonce", nonce)
				case <-abort:
					logger.Trace("Mithash nonce found but discarded", "attempts", nonce-seed, "nonce", nonce)
				}
				break search
			}
			nonce++
		}
	}
	// Datasets are unmapped in a finalizer. Ensure that the dataset stays live
	// during sealing so it's not unmapped while being read.
	runtime.KeepAlive(dataset)
}



//pos count
func posMine(nonce uint64,hash []byte,blockTime *big.Int,prevHash common.Hash,uncleHash common.Hash) ([]byte){
	// Combine header+nonce into a 64 byte seed
	seed := make([]byte, 40)
	copy(seed, hash)
	binary.LittleEndian.PutUint64(seed[32:], nonce)

	seed = crypto.Keccak512(seed)
	//fmt.Println("seed长度",len(seed))

	currentBlockTimeByte:=make([]byte,common.HashLength)
	copy(currentBlockTimeByte,blockTime.Bytes())
	currentBlockTimeByte=crypto.Keccak512(currentBlockTimeByte)

	posBytes:=append(seed,currentBlockTimeByte...)

	prevHashByte:=make([]byte,common.HashLength)
	copy(prevHashByte,prevHash.Bytes())
	prevHashByte=crypto.Keccak512(prevHashByte)

	posBytes=append(posBytes,prevHashByte...)

	uncleHashByte:=make([]byte,common.HashLength)
	copy(uncleHashByte,uncleHash.Bytes())
	uncleHashByte=crypto.Keccak512(uncleHashByte)

	posBytes=append(posBytes,uncleHashByte...)

	return crypto.Keccak256(append(seed,posBytes...))
}

