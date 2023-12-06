package common

import (
	"testing"

	"github.com/iquidus/blockspider/disk"
)

const (
	blockNumber  = 18721004
	logs         = 383
	txns         = 273
	blockPath    = "../testdata/eth-block-18721004.json"
	receiptsPath = "../testdata/eth-txn-receipts-18721004.json"
)

func TestConvert(t *testing.T) {
	// read test block from file
	var rawBlock RawBlock
	err := disk.ReadJsonFile[RawBlock](blockPath, &rawBlock)
	if err != nil {
		t.Error("Error reading file: ", err)
	}
	// read test receipts from file
	var receipts []RawTransactionReceipt
	err = disk.ReadJsonFile[[]RawTransactionReceipt](receiptsPath, &receipts)
	if err != nil {
		t.Error("Error reading file: ", err)
	}

	// convert raw block to common block
	block, err := rawBlock.Convert(nil, &receipts)
	if err != nil {
		t.Error("Error converting block: ", err)
	}

	// block number should be 18721004
	if block.Number != blockNumber {
		t.Errorf("TestConvert height = %d; want %d", block.Number, blockNumber)
	}

	// transaction count should be 273
	if block.TransactionCount != txns {
		t.Errorf("TestConvert txn count = %d; want %d", block.TransactionCount, txns)
	}

	// len(logs) should be 383
	lc := len(block.Logs)
	if lc != logs {
		t.Errorf("TestConvert log count = %d; want %d", lc, logs)
	}
}
