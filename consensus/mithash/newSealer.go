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


var (
	// limit for pow+pos
	BaseLimitBalanceNumber		*big.Int=big.NewInt(1600000)	//the epoch number
	BaselimitPosBalance	*big.Int = big.NewInt(1e+7)	//the first balance limit
	BaselimitSubBalance *big.Int = big.NewInt(1e+6)	//the minuend
	LastLimitBalance *big.Int=big.NewInt(1e5)//the last limitbalance is 100000
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
		return mithash.NewSeal(chain, block, stop,stateDb)
	}
	// Wait for all miners to terminate and return the block
	pend.Wait()
	return result, nil
}

// mine is the actual proof-of-work miner that searches for a nonce starting from
// seed that results in correct final block difficulty.
func (mithash *Mithash) newMine(block *types.Block, id int, seed uint64, abort chan struct{}, found chan *types.Block,balance *big.Int) {
	balanceTarget:=new(big.Int).Mul(balance,big.NewInt(1))
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
	//get the pow+pos balance limit
	balanceForLimit:=limitBalance(header.Number)

	// Start generating random nonces until we abort or find a good one
	var (
		attempts = int64(0)
		nonce    = seed
	)
	logger := log.New("miner", id)
	logger.Trace("Started mithash search for new nonces", "seed", seed)
	if balanceTarget.Cmp(balanceForLimit)>=0{
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
	}else{
		logger.Warn("Please try to get more TNB for miner coinbase!")
		//abort
		<-abort
	}

	// Datasets are unmapped in a finalizer. Ensure that the dataset stays live
	// during sealing so it's not unmapped while being read.
	runtime.KeepAlive(dataset)
}


/**
the limit balance of the balance,has div 1e18
 */
/**
the limit balance of the balance,has div 1e18
 */
func limitBalance(blockNumber *big.Int) *big.Int{
	limitBalance:=big.NewInt(0)
	limitBalance.Add(BaselimitPosBalance,big.NewInt(0))
	blockEpoch:=big.NewInt(0)

	//sub 1
	usedBlockNumber:=big.NewInt(0)
	if blockNumber.Cmp(big.NewInt(1))>0{
		usedBlockNumber.Sub(blockNumber,big.NewInt(1))
	}else{
		usedBlockNumber.Add(blockNumber,usedBlockNumber)
	}

	blockEpoch=blockEpoch.Div(usedBlockNumber,BaseLimitBalanceNumber)
	//start reward minus
	startEpoch:=big.NewInt(0)
	endEpoch:=big.NewInt(10)
	if blockEpoch.Cmp(startEpoch)>0&&blockEpoch.Cmp(endEpoch)<0{
		y:=blockEpoch.Sub(blockEpoch,startEpoch)
		subBalance:=new(big.Int).Mul(BaselimitSubBalance,y)
		limitBalance=limitBalance.Sub(limitBalance,subBalance)
	}
	if blockEpoch.Cmp(endEpoch)>=0{
		limitBalance=LastLimitBalance
	}

	return limitBalance
}



//pos count
func posMine(nonce uint64,hash []byte,blockTime *big.Int,prevHash common.Hash,uncleHash common.Hash) ([]byte){
	// Combine header+nonce into a 64 byte seed
	seed := make([]byte, 40)
	copy(seed, hash)
	binary.LittleEndian.PutUint64(seed[32:], nonce)

	seed = crypto.Keccak512(seed)
	//fmt.Println("seed length",len(seed))

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

