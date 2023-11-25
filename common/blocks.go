package common

import (
	"github.com/iquidus/blockspider/util"
)

type RawBlockDetails struct {
	Number string `bson:"number" json:"number"`
	Hash   string `bson:"hash" json:"hash"`
}

func (rbn *RawBlockDetails) Convert() (uint64, string) {
	return util.DecodeHex(rbn.Number), rbn.Hash
}

type RawBlock struct {
	Number           string           `bson:"number" json:"number"`
	Timestamp        string           `bson:"timestamp" json:"timestamp"`
	Transactions     []RawTransaction `bson:"transactions" json:"transactions"`
	Hash             string           `bson:"hash" json:"hash"`
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

func (b *RawBlock) Convert() Block {
	baseFeePerGas := util.DecodeValueHex(b.BaseFeePerGas)
	return Block{
		Number:          util.DecodeHex(b.Number),
		Timestamp:       util.DecodeHex(b.Timestamp),
		Transactions:    make([]Transaction, len(b.Transactions)),
		RawTransactions: b.Transactions,
		Hash:            b.Hash,
		ParentHash:      b.ParentHash,
		BaseFeePerGas:   baseFeePerGas,
	}
}

type Block struct {
	Number          uint64           `bson:"number" json:"number"`
	Timestamp       uint64           `bson:"timestamp" json:"timestamp"`
	Transactions    []Transaction    `bson:"transactions" json:"transactions"`
	RawTransactions []RawTransaction `bson:"-" json:"-"`
	Hash            string           `bson:"hash" json:"hash"`
	ParentHash      string           `bson:"parentHash" json:"parentHash"`
	BaseFeePerGas   string           `bson:"baseFeePerGas" json:"baseFeePerGas,omitempty"`
}
