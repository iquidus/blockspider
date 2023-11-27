package common

import (
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
		GasUsed:         util.DecodeHex(b.GasUsed),
		GasLimit:        util.DecodeHex(b.GasLimit),
		MixHash:         b.MixHash,
		StateRoot:       b.StateRoot,
		TotalDifficulty: b.TotalDifficulty,
		Miner:           b.Miner,
		Difficulty:      b.Difficulty,
		Sha3Uncles:      b.Sha3Uncles,
	}
}

type Block struct {
	Number          uint64           `bson:"number" json:"number"`
	Timestamp       uint64           `bson:"timestamp" json:"timestamp"`
	Hash            string           `bson:"hash" json:"hash"`
	ParentHash      string           `bson:"parentHash" json:"parentHash"`
	Transactions    []Transaction    `bson:"transactions" json:"transactions,omitempty"`
	RawTransactions []RawTransaction `bson:"-" json:"-"`
	BaseFeePerGas   string           `bson:"baseFeePerGas" json:"baseFeePerGas,omitempty"`
	GasUsed         uint64           `bson:"gasUsed" json:"gasUsed,omitempty"`
	GasLimit        uint64           `bson:"gasLimit" json:"gasLimit,omitempty"`
	MixHash         string           `bson:"mixHash" json:"mixHash,omitempty"`
	StateRoot       string           `bson:"stateRoot" json:"stateRoot,omitempty"`
	TotalDifficulty string           `bson:"totalDifficulty" json:"totalDifficulty,omitempty"`
	Sha3Uncles      string           `bson:"sha3Uncles" json:"sha3Uncles,omitempty"`
	Miner           string           `bson:"miner" json:"miner,omitempty"`
	Difficulty      string           `bson:"difficulty" json:"difficulty,omitempty"`
}
