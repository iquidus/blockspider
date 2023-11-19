package common

import (
	"github.com/iquidus/blockspider/util"
)

type RawTransaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
	Type             string `bson:"type" json:"type,omitempty"`
	BaseFeePerGas    string `bson:"baseFeePerGas" json:"baseFeePerGas,omitempty"`
}

func (rt *RawTransaction) Convert() Transaction {
	return Transaction{
		BlockHash:        rt.BlockHash,
		BlockNumber:      util.DecodeHex(rt.BlockNumber),
		Hash:             rt.Hash,
		Input:            rt.Input,
		Value:            util.DecodeValueHex(rt.Value),
		Gas:              util.DecodeHex(rt.Gas),
		GasPrice:         util.DecodeHex(rt.GasPrice),
		Nonce:            rt.Nonce,
		TransactionIndex: util.DecodeHex(rt.TransactionIndex),
		From:             rt.From,
		To:               rt.To,
		Type:             rt.Type,
		BaseFeePerGas:    rt.BaseFeePerGas,
	}
}

type Transaction struct {
	BlockHash        string `bson:"blockHash" json:"blockHash"`
	BlockNumber      uint64 `bson:"blockNumber" json:"blockNumber"`
	Hash             string `bson:"hash" json:"hash"`
	Timestamp        uint64 `bson:"timestamp" json:"timestamp"`
	Input            string `bson:"input" json:"input"`
	Value            string `bson:"value" json:"value"`
	Gas              uint64 `bson:"gas" json:"gas"`
	GasPrice         uint64 `bson:"gasPrice" json:"gasPrice"`
	Nonce            string `bson:"nonce" json:"nonce"`
	TransactionIndex uint64 `bson:"transactionIndex" json:"transactionIndex"`
	From             string `bson:"from" json:"from"`
	To               string `bson:"to" json:"to"`
	Status           bool   `json:"status"`
	//
	GasUsed         uint64  `bson:"gasUsed" json:"gasUsed"`
	ContractAddress string  `bson:"contractAddress" json:"contractAddress"`
	Logs            []TxLog `bson:"logs" json:"logs"`
	//
	MaxFeePerGas         uint64 `bson:"maxFeePerGas" json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas uint64 `bson:"maxPriorityFeePerGas" json:"maxPriorityFeePerGas,omitempty"`
	Type                 string `bson:"type" json:"type,omitempty"`
	BaseFeePerGas        string `bson:"baseFeePerGas" json:"baseFeePerGas,omitempty"`
}

type TxLog struct {
	Address          string   `bson:"address" json:"address"`
	Topics           []string `bson:"topics" json:"topics"`
	Data             string   `bson:"data" json:"data"`
	BlockNumber      string   `bson:"blockNumber" json:"blockNumber"`
	TransactionIndex string   `bson:"transactionIndex" json:"transactionIndex"`
	TransactionHash  string   `bson:"transactionHash" json:"transactionHash"`
	BlockHash        string   `bson:"blockHash" json:"blockHash"`
	LogIndex         string   `bson:"logIndex" json:"logIndex"`
	Removed          bool     `bson:"removed" json:"removed"`
}

type TxLogRequest struct {
	Address   []string `bson:"address" json:"address"`
	Topics    []string `bson:"topics" json:"topics"`
	BlockHash string   `bson:"blockHash" json:"blockHash"`
}
