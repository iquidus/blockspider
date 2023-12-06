package common

import (
	"errors"

	"github.com/iquidus/blockspider/util"
)

type RawBlock struct {
	Hash             string           `bson:"hash" json:"hash"`
	Number           string           `bson:"number" json:"number"`
	Timestamp        string           `bson:"timestamp" json:"timestamp"`
	Transactions     []RawTransaction `bson:"transactions" json:"transactions"`
	ParentHash       string           `bson:"parentHash" json:"parentHash"`
	Sha3Uncles       string           `bson:"sha3Uncles" json:"sha3Uncles"`
	Miner            string           `bson:"miner" json:"miner"`
	MixHash          string           `bson:"mixHash" json:"mixHash"`
	Difficulty       string           `bson:"difficulty" json:"difficulty"`
	TotalDifficulty  string           `bson:"totalDifficulty" json:"totalDifficulty"`
	Size             string           `bson:"size" json:"size"`
	GasUsed          string           `bson:"gasUsed" json:"gasUsed"`
	GasLimit         string           `bson:"gasLimit" json:"gasLimit"`
	Nonce            string           `bson:"nonce" json:"nonce"`
	Uncles           []string         `bson:"uncles" json:"uncles"`
	BaseFeePerGas    string           `bson:"baseFeePerGas" json:"baseFeePerGas,omitempty"`
	ExtraData        string           `bson:"extraData" json:"extraData"`
	LogsBloom        string           `bson:"logsBloom" json:"logsBloom"`
	ReceiptsRoot     string           `bson:"receiptsRoot" json:"receiptsRoot"`
	StateRoot        string           `bson:"stateRoot" json:"stateRoot"`
	TransactionsRoot string           `bson:"transactionsRoot" json:"transactionsRoot"`
}

// TODO(iquidus): refactor this, separate out txn receipts without introducing any additional looping
func (b *RawBlock) Convert(rpcClient *RPCClient, receipts *[]RawTransactionReceipt) (Block, error) {
	// make sure we have either an rpc client or txn receipts
	if receipts == nil && rpcClient == nil {
		return Block{}, errors.New("cannot convert block without receipts or rpc client")
	}

	baseFeePerGas := util.DecodeValueHex(b.BaseFeePerGas)
	// handle getting logs and txn receipts here
	txns := make([]Transaction, len(b.Transactions))
	var logs []Log
	for i, txn := range b.Transactions {
		var receipt RawTransactionReceipt
		if receipts != nil {
			receipt = (*receipts)[i]
		} else {
			// get transaction receipts
			r, err := rpcClient.GetTransactionReceipt(txn.Hash)
			if err != nil {
				return Block{}, err
			}
			receipt = *r
		}

		// convert raw txn to txn
		txns[i] = txn.Convert(receipt)

		// get logs
		for _, log := range receipt.Logs {
			logs = append(logs, log.Convert(txns[i]))
		}
	}

	return Block{
		Number:           util.DecodeHex(b.Number),
		Timestamp:        util.DecodeHex(b.Timestamp),
		Transactions:     txns,
		Hash:             b.Hash,
		ParentHash:       b.ParentHash,
		BaseFeePerGas:    baseFeePerGas,
		GasUsed:          util.DecodeHex(b.GasUsed),
		GasLimit:         util.DecodeHex(b.GasLimit),
		MixHash:          b.MixHash,
		StateRoot:        b.StateRoot,
		TotalDifficulty:  b.TotalDifficulty,
		Miner:            b.Miner,
		Difficulty:       b.Difficulty,
		Sha3Uncles:       b.Sha3Uncles,
		Nonce:            b.Nonce,
		TransactionCount: uint64(len(b.Transactions)),
		TransactionsRoot: b.TransactionsRoot,
		ReceiptsRoot:     b.ReceiptsRoot,
		LogsBloom:        b.LogsBloom,
		ExtraData:        b.ExtraData,
		Uncles:           b.Uncles,
		Logs:             logs,
	}, nil
}

type Block struct {
	Number           uint64        `bson:"number" json:"number"`
	Timestamp        uint64        `bson:"timestamp" json:"timestamp"`
	Hash             string        `bson:"hash" json:"hash"`
	ParentHash       string        `bson:"parentHash" json:"parentHash"`
	Transactions     []Transaction `bson:"transactions" json:"transactions,omitempty"`
	BaseFeePerGas    string        `bson:"baseFeePerGas" json:"baseFeePerGas,omitempty"`
	GasUsed          uint64        `bson:"gasUsed" json:"gasUsed,omitempty"`
	GasLimit         uint64        `bson:"gasLimit" json:"gasLimit,omitempty"`
	MixHash          string        `bson:"mixHash" json:"mixHash,omitempty"`
	StateRoot        string        `bson:"stateRoot" json:"stateRoot,omitempty"`
	TotalDifficulty  string        `bson:"totalDifficulty" json:"totalDifficulty,omitempty"`
	Sha3Uncles       string        `bson:"sha3Uncles" json:"sha3Uncles,omitempty"`
	Miner            string        `bson:"miner" json:"miner,omitempty"`
	Difficulty       string        `bson:"difficulty" json:"difficulty,omitempty"`
	Nonce            string        `bson:"nonce" json:"nonce,omitempty"`
	TransactionCount uint64        `bson:"transactionCount" json:"transactionCount,omitempty"`
	TransactionsRoot string        `bson:"transactionsRoot" json:"transactionsRoot,omitempty"`
	ReceiptsRoot     string        `bson:"receiptsRoot" json:"receiptsRoot,omitempty"`
	LogsBloom        string        `bson:"logsBloom" json:"logsBloom,omitempty"`
	ExtraData        string        `bson:"extraData" json:"extraData,omitempty"`
	Uncles           []string      `bson:"uncles" json:"uncles,omitempty"`
	Logs             []Log         `bson:"logs" json:"logs,omitempty"`
}

// TODO(iquidus): write a compact function for Block
// that returns a minimal representation of the block
// e.g: number, hash, parentHash, timestamp
