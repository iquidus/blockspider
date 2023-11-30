package common

import "github.com/iquidus/blockspider/util"

type RawLog struct {
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

func (l *RawLog) Convert(txn Transaction) Log {
	return Log{
		Address:     l.Address,
		Topics:      l.Topics,
		Data:        l.Data,
		Index:       util.DecodeHex(l.LogIndex),
		Transaction: txn,
	}
}

type Log struct {
	Address     string   `bson:"address" json:"address"`
	Topics      []string `bson:"topics" json:"topics"`
	Data        string   `bson:"data" json:"data"`
	Index       uint64   `bson:"index" json:"index"`
	Transaction Transaction
}

type LogRequest struct {
	Address   []string `bson:"address" json:"address"`
	Topics    []string `bson:"topics" json:"topics"`
	BlockHash string   `bson:"blockHash" json:"blockHash"`
}
