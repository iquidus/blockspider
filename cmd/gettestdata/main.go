package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/disk"
	"github.com/iquidus/blockspider/params"
)

var (
	configFileName string

	height uint64

	mainLogger log.Logger
)

const (
	configFlagDefault = "config.json"
	configFlagDesc    = "specify name of config file (should be in working dir)"

	heightFlagDefault = 0
	heightFlagDesc    = "block number to retrieve for testdada"
)

func init() {
	flag.StringVar(&configFileName, "c", configFlagDefault, configFlagDesc)
	flag.StringVar(&configFileName, "config", configFlagDefault, configFlagDesc)

	flag.Uint64Var(&height, "n", heightFlagDefault, heightFlagDesc)
	flag.Uint64Var(&height, "number", heightFlagDefault, heightFlagDesc)

	flag.Parse()

	mainLogger = *log.New(os.Stdout, "", log.LstdFlags)
}

func generateReceipts(rpc *common.RPCClient) {
	// get block by number
	rawBlock, err := rpc.GetBlockByHeight(height)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	// write raw block to file
	err = disk.WriteJsonFile[common.RawBlock](rawBlock, fmt.Sprintf("./testdata/eth-block-%d.json", height), 0644)
	if err != nil {
		log.Fatal("Error writing to file: ", err)
	}
	// get transaction receipts
	receipts := make([]common.RawTransactionReceipt, len(rawBlock.Transactions))
	for i, txn := range rawBlock.Transactions {
		receipt, err := rpc.GetTransactionReceipt(txn.Hash)
		if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
		}
		receipts[i] = *receipt
	}
	// write receipts to file
	err = disk.WriteJsonFile[[]common.RawTransactionReceipt](receipts, fmt.Sprintf("./testdata/eth-txn-receipts-%d.json", height), 0644)
	if err != nil {
		log.Fatal("Error writing to file: ", err)
	}
}

func main() {
	mainLogger.Print("blockspider/gettestdata ", params.VersionWithMeta)
	// Read config
	var cfg params.Config
	configPath, err := filepath.Abs(configFileName)
	if err != nil {
		mainLogger.Fatal("Error: could not parse config filepath", "err", err)
	}
	err = disk.ReadJsonFile[params.Config](configPath, &cfg)
	if err != nil {
		mainLogger.Fatal("Error: could read config file", "err", err)
	}
	rpcClient := common.NewRPCClient(&cfg.Rpc)
	generateReceipts(rpcClient)
}
