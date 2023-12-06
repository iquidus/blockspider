package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/params"
)

var (
	cfg            params.Config
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
	towrite, err := json.MarshalIndent(rawBlock, "", "  ")
	if err != nil {
		log.Fatal("Error marshalling block: ", err)
	}
	err = os.WriteFile(fmt.Sprintf("./testdata/eth-block-%d.json", height), towrite, 0644)
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
	towrite, err = json.MarshalIndent(receipts, "", "  ")
	if err != nil {
		log.Fatal("Error marshalling receipts: ", err)
	}
	err = os.WriteFile(fmt.Sprintf("./testdata/eth-txn-receipts-%d.json", height), towrite, 0644)
	if err != nil {
		log.Fatal("Error writing to file: ", err)
	}
}

func readConfig(cfg *params.Config) {
	if configFileName == "" {
		mainLogger.Fatal("Invalid arguments", os.Args)
		os.Exit(1)
	}

	confPath, err := filepath.Abs(configFileName)
	if err != nil {
		mainLogger.Fatal("Error: could not parse config filepath", "err", err)
	}

	mainLogger.Print("Loading config", "path", confPath)

	configFile, err := os.Open(confPath)
	if err != nil {
		mainLogger.Fatal("File error", "err", err.Error())
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		mainLogger.Fatal("Config error", "err", err.Error())
	}
}

func main() {
	mainLogger.Print("blockspider/gettestdata ", params.VersionWithMeta)
	readConfig(&cfg)
	rpcClient := common.NewRPCClient(&cfg.Rpc)
	generateReceipts(rpcClient)
}
