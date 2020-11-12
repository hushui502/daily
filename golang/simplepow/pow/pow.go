package pow

import (
	"bytes"
	crand "crypto/rand"
	"crypto/sha256"
	"errors"
	"lib/simplepow/util"
	"log"
	"math"
	"math/big"
	"math/rand"
	"runtime"
)

type Block struct {
	Header *Header
	body Body
}

type Body struct {
	
}

type Header struct {
	Difficulty *big.Int
}

type ProofOfWorkManager struct {
	threads int
	rand    *rand.Rand
	logger  log.Logger
}

type ProofOfWorkManagerOptionFunc func(options *ProofOfWorkManagerOption)

type ProofOfWorkManagerOption struct {
	threads int
	rand *rand.Rand
}

func WithThreads(threads int) ProofOfWorkManagerOptionFunc {
	return func(options *ProofOfWorkManagerOption) {
		options.threads = threads
	}
}

func WithRand(rand *rand.Rand) ProofOfWorkManagerOptionFunc {
	return func(options *ProofOfWorkManagerOption) {
		options.rand = rand
	}
}

func NewProofOfWorkManager(looger log.Logger, opts ...ProofOfWorkManagerOptionFunc) (*ProofOfWorkManager, error) {
	option := ProofOfWorkManagerOption{}
	for _, o := range opts {
		o(&option)
	}
	if option.threads == 0 {
		option.threads = runtime.NumCPU()
	}
	if option.rand == nil {
		seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			return nil, err
		}
		option.rand = rand.New(rand.NewSource(seed.Int64()))
	}

	return &ProofOfWorkManager{
		threads: option.threads,
		rand: option.rand,
		logger:looger,
	}, nil
}

type Result struct {
	Nonce uint64
	AttemptNum int64
	Hash []byte
}

func (proofOfWorkManager *ProofOfWorkManager) Work(block *Block, result chan *Result) error {
	if proofOfWorkManager.threads < 0 {
		return errors.New("threads set error")
	}
	abort := make(chan string)
	for i := 0; i < proofOfWorkManager.threads; i++ {
		initNonce := uint64(proofOfWorkManager.rand.Int63())
		proofOfWorkManager.logger.Printf("Thread %d: initNonce: %d", i, initNonce)
		go func(id int, nonce uint64) {
			proofOfWorkManager.mine(block, id, initNonce, abort, result)
		}(i, initNonce)
	}

	return nil
}

func (proofWorkManager *ProofOfWorkManager) mine(block *Block, threadId int, nonce uint64, abort chan string, found chan *Result) {
	var (
		hash = proofWorkManager.hashHeader(block.Header)
		target = new(big.Int).Div(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)), block.Header.Difficulty)
		attempts = int64(0)
	)
	search:
		for {
			select {
			case <-abort:
				proofWorkManager.logger.Printf("Thread %d: abort", threadId)
				break search
			default:
				attempts++
				hash := sha256.Sum256(bytes.Join([][]byte{
				hash,
				util.MustToBuffer(nonce),
				}, []byte{}))
				realHash := hash[:]
				if attempts % 1000000 == 0 {
					proofWorkManager.logger.Printf("Thread %d: attempted %d currentHash: %s", threadId, attempts, util.BufferToHexString(realHash, true))
				}
				if new(big.Int).SetBytes(realHash).Cmp(target) < 0 {
					proofWorkManager.logger.Printf("Thread %d: Found result", threadId)
					select {
					case found <- &Result{
						Nonce:      nonce,
						AttemptNum: attempts,
						Hash:       realHash,
					}:
						proofWorkManager.logger.Printf("Thread %d: nonce found and reported", threadId)
					case <-abort:
						proofWorkManager.logger.Printf("Thread %d: nonce found but discarded", threadId)
					}
					go func() {
						proofWorkManager.logger.Printf("Thread %d: stop all calc", threadId)
						close(abort)
					}()
				}
				nonce++
			}

		}
}

func (proofOfWorkManager *ProofOfWorkManager) hashHeader(header *Header) []byte {
	hash := sha256.Sum256(bytes.Join([][]byte{
		util.MustToBuffer(header.Difficulty.Int64()),
	}, []byte{}))

	return hash[:]
}























