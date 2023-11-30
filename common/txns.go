package common

import (
	"github.com/iquidus/blockspider/util"
)

type RawTransaction struct {
	BlockHash            string `json:"blockHash"`
	BlockNumber          string `json:"blockNumber"`
	From                 string `json:"from"`
	Gas                  string `json:"gas"`
	GasPrice             string `json:"gasPrice"`
	MaxFeePerGas         string `bson:"maxFeePerGas" json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string `bson:"maxPriorityFeePerGas" json:"maxPriorityFeePerGas,omitempty"`
	Hash                 string `json:"hash"`
	Input                string `json:"input"`
	Nonce                string `json:"nonce"`
	To                   string `json:"to"`
	TransactionIndex     string `json:"transactionIndex"`
	Value                string `json:"value"`
	Type                 string `bson:"type" json:"type,omitempty"`
	ChainId              string `json:"chainId"`
	V                    string `json:"v"`
	R                    string `json:"r"`
	S                    string `json:"s"`
}

type RawTransactionReceipt struct {
	BlockHash         string   `bson:"blockHash" json:"blockHash"`
	BlockNumber       string   `bson:"blockNumber" json:"blockNumber"`
	ContractAddress   string   `bson:"contractAddress" json:"contractAddress"`
	CumulativeGasUsed string   `bson:"cumulativeGasUsed" json:"cumulativeGasUsed"`
	From              string   `bson:"from" json:"from"`
	EffectiveGasPrice string   `bson:"effectiveGasPrice" json:"effectiveGasPrice"`
	GasUsed           string   `bson:"gasUsed" json:"gasUsed"`
	Logs              []RawLog `bson:"logs" json:"logs"`
	LogsBloom         string   `bson:"logsBloom" json:"logsBloom"`
	Status            string   `bson:"status" json:"status"`
	To                string   `bson:"to" json:"to"`
	TransactionHash   string   `bson:"transactionHash" json:"transactionHash"`
	TransactionIndex  string   `bson:"transactionIndex" json:"transactionIndex"`
	Type              string   `bson:"type" json:"type,omitempty"`
}

func (rt *RawTransaction) Convert(receipt RawTransactionReceipt) Transaction {
	return Transaction{
		// from txn
		From:                 rt.From,
		Gas:                  util.DecodeHex(rt.Gas),
		GasPrice:             util.DecodeHex(rt.GasPrice),
		Hash:                 rt.Hash,
		Index:                util.DecodeHex(rt.TransactionIndex),
		MaxFeePerGas:         util.DecodeHex(rt.MaxFeePerGas),
		MaxPriorityFeePerGas: util.DecodeHex(rt.MaxPriorityFeePerGas),
		Nonce:                util.DecodeHex(rt.Nonce),
		To:                   rt.To,
		Value:                util.DecodeValueHex(rt.Value),
		// from receipt
		Status:            util.DecodeHex(receipt.Status),
		GasUsed:           util.DecodeHex(receipt.GasUsed),
		CumulativeGasUsed: util.DecodeHex(receipt.CumulativeGasUsed),
		EffectiveGasPrice: util.DecodeHex(receipt.EffectiveGasPrice),
		CreatedContract:   receipt.ContractAddress,
	}
}

type Transaction struct {
	// from txn
	From                 string `bson:"from" json:"from"`
	Gas                  uint64 `bson:"gas" json:"gas"`
	GasPrice             uint64 `bson:"gasPrice" json:"gasPrice"`
	Hash                 string `bson:"hash" json:"hash"`
	Index                uint64 `bson:"index" json:"index"`
	MaxFeePerGas         uint64 `bson:"maxFeePerGas" json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas uint64 `bson:"maxPriorityFeePerGas" json:"maxPriorityFeePerGas,omitempty"`
	Nonce                uint64 `bson:"nonce" json:"nonce"`
	To                   string `bson:"to" json:"to"`
	Value                string `bson:"value" json:"value"`
	// from receipt
	Status            uint64 `json:"status"`
	GasUsed           uint64 `bson:"gasUsed" json:"gasUsed"`
	CumulativeGasUsed uint64 `bson:"cumulativeGasUsed" json:"cumulativeGasUsed,omitempty"`
	EffectiveGasPrice uint64 `bson:"effectiveGasPrice" json:"effectiveGasPrice,omitempty"`
	CreatedContract   string `bson:"createdContract" json:"createdContract,omitempty"`
}
