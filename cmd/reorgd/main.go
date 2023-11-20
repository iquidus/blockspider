package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iquidus/blockspider/cache"
	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/util"
	"golang.org/x/crypto/sha3"
)

const (
	BLOCKTIME = 2     // desired block time in seconds
	JSONRPC   = "2.0" // const for rpc responses
)

// some reusable bigInts
var (
	big1    = big.NewInt(1)
	big5    = big.NewInt(5)
	big7    = big.NewInt(7)
	big2048 = big.NewInt(2048)
)

type Request struct {
	Id      string        `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type Response struct {
	Id      string        `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Result  string        `json:"result,omitempty"`
	Error   *JsonRpcError `json:"error,omitempty"`
}

type BlockResponse struct {
	Id      string           `json:"id"`
	Jsonrpc string           `json:"jsonrpc"`
	Result  *common.RawBlock `json:"result"`
	Error   *JsonRpcError    `json:"error,omitempty"`
}

type LogResponse struct {
	Id      string   `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
	Result  []string `json:"result"`
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// returns a random keccak256 hash with 0x prefix
func randomKeccakHash() string {
	var seed, _ = rand.Int(rand.Reader, big2048)
	var hash = sha3.NewLegacyKeccak256()
	hash.Write(seed.Bytes())
	keccak := hash.Sum(nil)
	return fmt.Sprintf("0x%x", keccak)
}

// returns a dummy PoW block chained with the given parent block
func createPowBlock(parent *common.RawBlock) common.RawBlock {
	// start as if genesis
	var number uint64 = 0
	var parentHash string = "0x"
	timestamp := util.EncodeUint64(uint64(time.Now().Unix()))

	// check if a parent block is provided
	if parent != nil {
		// chain new block to parent
		number = util.DecodeHex(parent.Number)
		number++
		parentHash = parent.Hash
	}

	// generate a random "blockhash"
	var hash = randomKeccakHash()
	// return a block
	return common.RawBlock{
		Number:     util.EncodeUint64(number),
		Timestamp:  timestamp,
		Hash:       hash,
		ParentHash: parentHash,
		// TODO(iquidus): randomly generate values below
		Difficulty:      util.EncodeUint64(uint64(438231850248)),
		TotalDifficulty: util.EncodeUint64(uint64(2142877125748580710)),
		Size:            util.EncodeUint64(uint64(542)),
		GasUsed:         util.EncodeUint64(uint64(0)),
		GasLimit:        util.EncodeUint64(uint64(8000000)),
		Nonce:           util.EncodeUint64(number),
		BaseFeePerGas:   util.EncodeUint64(uint64(80000000)),
		ExtraData:       "0x",
	}
}

// main function (app entry)
func main() {
	// create a "blockchain" (a block stack using a doubly-linked list)
	blockchain := cache.New[common.RawBlock](nil)
	// create a map for faster lookups by block number
	blockmap := make(map[string]common.RawBlock)
	// generate a genesis block (no parent block)
	genesis := createPowBlock(nil)
	// add to blockchain
	blockchain.Push(genesis)
	blockmap[genesis.Number] = genesis

	// start "miner"
	miner := time.NewTicker(BLOCKTIME * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-miner.C:
				// roll a dice to determine reorg length (1-6 blocks)
				diceRoll1, _ := rand.Int(rand.Reader, big5)
				diceRoll1.Add(diceRoll1, big1)
				reorgLength := diceRoll1.Uint64()
				// make sure theres enough blocks in chain for this reorg
				if uint64(blockchain.Count()) > reorgLength {
					// roll a second dice (1-6)
					diceRoll2, _ := rand.Int(rand.Reader, big5)
					diceRoll2.Add(diceRoll2, big1)
					// combine results of both dice rolls
					diceRollCombined := big.NewInt(0)
					diceRollCombined.Add(diceRoll1, diceRoll2)
					// if we have rolled a 7, reorg
					if diceRollCombined.Cmp(big7) == 0 {
						// drop old blocks
						for i := 0; i < int(reorgLength); i++ {
							oldBlock, _ := blockchain.Pop()
							delete(blockmap, oldBlock.Number)
							fmt.Printf("Dropped old block, number: %d, hash: %s\n", util.DecodeHex(oldBlock.Number), oldBlock.Hash)
						}
						// add new blocks
						for i := 0; i < int(reorgLength); i++ {
							parent, _ := blockchain.Peak()
							newBlock := createPowBlock(&parent)
							blockchain.Push(newBlock)
							blockmap[newBlock.Number] = newBlock
							fmt.Printf("Mined new block, number: %d, hash: %s\n", util.DecodeHex(newBlock.Number), newBlock.Hash)
						}
					}
				}
				// create a new block using the chains head as parent
				parent, _ := blockchain.Peak()
				newBlock := createPowBlock(&parent)
				// add block to chain
				blockchain.Push(newBlock)
				blockmap[newBlock.Number] = newBlock
				// log
				fmt.Printf("Mined new block, number: %d, hash: %s\n", util.DecodeHex(newBlock.Number), newBlock.Hash)
			}
		}
	}()
	// start api
	router := setupRouter(blockchain, blockmap)
	router.Run(":8079")
}

func setupRouter(blockchain *cache.BlockStack[common.RawBlock], blockmap map[string]common.RawBlock) *gin.Engine {
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.POST("/", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		var req Request
		json.Unmarshal(body, &req)

		switch req.Method {
		case "web3_clientVersion":
			var res Response
			res.Id = req.Id
			res.Jsonrpc = "2.0"
			res.Result = "Reorgd/v0.0.1-develop"
			res.Error = nil
			c.JSON(http.StatusOK, res)
		case "eth_blockNumber":
			// return number of head block (as hex string)
			head, err := blockchain.Peak()
			if err == nil {
				res := Response{
					Id:      req.Id,
					Jsonrpc: JSONRPC,
					Result:  head.Number,
					Error:   nil,
				}
				c.JSON(http.StatusOK, res)
			} else {
				c.JSON(http.StatusOK, "0x0")
			}
		case "eth_getBlockByNumber":
			res := BlockResponse{Id: req.Id, Jsonrpc: JSONRPC, Result: nil}
			// get blocknumber (hex string) value from params
			// if key in blockmap exists return block as result
			bn := fmt.Sprintf("%v", req.Params[0])
			switch bn {
			case "latest":
				b, err := blockchain.Peak()
				if err == nil {
					res.Result = &b
				}
			case "earliest":
				b, ok := blockmap["0x0"]
				if ok {
					res.Result = &b
				}
			case "pending":
				res.Error.Code = -39001
				res.Error.Message = "-39001: Unknown block"
			case "finalized":
				res.Error.Code = -39001
				res.Error.Message = "-39001: Unknown block"
			case "safe":
				res.Error.Code = -39001
				res.Error.Message = "-39001: Unknown block"
			default:
				b, ok := blockmap[bn]
				if ok {
					res.Result = &b
				}
			}
			c.JSON(http.StatusOK, res)
		case "eth_getLogs":
			// TODO(iquidus): generate random (but useful) log data (e.g: random erc20 transfers)
			res := LogResponse{Id: req.Id, Jsonrpc: JSONRPC}
			logs := make([]string, 0)
			res.Result = logs
			c.JSON(http.StatusOK, res)
		default:
			// TODO(iquidus): handle this properly
			fmt.Printf("Unhandled method: %s, params: %v", req.Method, req.Params)
			c.JSON(http.StatusOK, nil)
		}
	})

	return router
}
