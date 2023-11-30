package common

import (
	"github.com/iquidus/blockspider/util"
)

type AlchemyConfig struct {
	Secret string `json:"secret"` // alchemy webhook signing key
}

type AlchemyWebhookBlockLogs struct {
	Data        string                    `bson:"data" json:"data"`
	Topics      []string                  `bson:"topics" json:"topics"`
	Index       uint64                    `bson:"index" json:"index"`
	Account     AlchemyWebhookAccount     `bson:"account" json:"account"`
	Transaction AlchemyWebhookTransaction `bson:"transaction" json:"transaction"`
}

func (l *AlchemyWebhookBlockLogs) Convert() Log {
	return Log{
		Address:     l.Account.Address,
		Topics:      l.Topics,
		Data:        l.Data,
		Index:       l.Index,
		Transaction: l.Transaction.Convert(),
	}
}

type AlchemyWebhookTransaction struct {
	Hash                 string                 `bson:"hash" json:"hash"`
	Nonce                uint64                 `bson:"nonce" json:"nonce"`
	Index                uint64                 `bson:"index" json:"index"`
	From                 AlchemyWebhookAccount  `bson:"from" json:"from"`
	To                   AlchemyWebhookAccount  `bson:"to" json:"to"`
	Value                string                 `bson:"value" json:"value"`
	GasPrice             string                 `bson:"gasPrice" json:"gasPrice"`
	MaxFeePerGas         string                 `bson:"maxFeePerGas" json:"maxFeePerGas"`
	MaxPriorityFeePerGas string                 `bson:"maxPriorityFeePerGas" json:"maxPriorityFeePerGas"`
	Gas                  uint64                 `bson:"gas" json:"gas"`
	Status               uint64                 `bson:"status" json:"status"`
	GasUsed              uint64                 `bson:"gasUsed" json:"gasUsed"`
	CumulativeGasUsed    uint64                 `bson:"cumulativeGasUsed" json:"cumulativeGasUsed"`
	EffectiveGasPrice    string                 `bson:"effectiveGasPrice" json:"effectiveGasPrice"`
	CreatedContract      *AlchemyWebhookAccount `bson:"createdContract" json:"createdContract"`
}

func (l *AlchemyWebhookTransaction) Convert() Transaction {
	return Transaction{
		Hash:                 l.Hash,
		Nonce:                l.Nonce,
		Index:                l.Index,
		From:                 l.From.Address,
		To:                   l.To.Address,
		Value:                util.DecodeValueHex(l.Value),
		GasPrice:             util.DecodeHex(l.GasPrice),
		MaxFeePerGas:         util.DecodeHex(l.MaxFeePerGas),
		MaxPriorityFeePerGas: util.DecodeHex(l.MaxPriorityFeePerGas),
		Gas:                  l.Gas,
		Status:               l.Status,
		GasUsed:              l.GasUsed,
		CumulativeGasUsed:    l.CumulativeGasUsed,
		EffectiveGasPrice:    util.DecodeHex(l.EffectiveGasPrice),
		CreatedContract:      l.CreatedContract.Address,
	}
}

type AlchemyWebhookAccount struct {
	Address string `bson:"address" json:"address"`
}

type AlchemyWebhookParentBlock struct {
	Hash string `bson:"hash" json:"hash"`
}

type AlchemyWebhookBlock struct {
	Hash             string                      `bson:"hash" json:"hash"`
	Number           uint64                      `bson:"number" json:"number"`
	Timestamp        uint64                      `bson:"timestamp" json:"timestamp"`
	Parent           AlchemyWebhookParentBlock   `bson:"parent" json:"parent"`
	BaseFeePerGas    string                      `bson:"baseFeePerGas" json:"baseFeePerGas,omitempty"`
	GasUsed          uint64                      `bson:"gasUsed" json:"gasUsed"`
	GasLimit         uint64                      `bson:"gasLimit" json:"gasLimit"`
	MixHash          string                      `bson:"mixHash" json:"mixHash"`
	StateRoot        string                      `bson:"stateRoot" json:"stateRoot"`
	Difficulty       string                      `bson:"difficulty" json:"difficulty"`
	TotalDifficulty  string                      `bson:"totalDifficulty" json:"totalDifficulty"`
	Nonce            string                      `bson:"nonce" json:"nonce"`
	TransactionCount uint64                      `bson:"transactionCount" json:"transactionCount"`
	TransactionsRoot string                      `bson:"transactionsRoot" json:"transactionsRoot"`
	ReceiptsRoot     string                      `bson:"receiptsRoot" json:"receiptsRoot"`
	LogsBloom        string                      `bson:"logsBloom" json:"logsBloom"`
	Transactions     []AlchemyWebhookTransaction `bson:"transactions" json:"transactions"`
	Logs             []AlchemyWebhookBlockLogs   `bson:"logs" json:"logs"`
}

type AlchemyEventData struct {
	Block AlchemyWebhookBlock `bson:"block" json:"block"`
}

type AlchemyEvent struct {
	Data           AlchemyEventData `bson:"data" json:"data"`
	SequenceNumber string           `bson:"sequenceNumber" json:"sequenceNumber"`
}

type AlchemyWebhook struct {
	WebhookId string       `bson:"webhookId" json:"webhookId"`
	Id        string       `bson:"id" json:"id"`
	CreatedAt string       `bson:"createdAt" json:"createdAt"`
	Type      string       `bson:"type" json:"type"`
	Event     AlchemyEvent `bson:"event" json:"event"`
}

func (b *AlchemyWebhookBlock) Convert() Block {
	baseFeePerGas := util.DecodeValueHex(b.BaseFeePerGas)
	// txns := make([]Transaction, len(b.Transactions))
	// for i, txn := range b.Transactions {
	// 	txns[i] = txn.Convert()
	// }
	logs := make([]Log, len(b.Logs))
	for i, log := range b.Logs {
		logs[i] = log.Convert()
	}
	return Block{
		Hash:             b.Hash,
		Number:           b.Number,
		Timestamp:        b.Timestamp,
		ParentHash:       b.Parent.Hash,
		BaseFeePerGas:    baseFeePerGas,
		GasUsed:          b.GasUsed,
		GasLimit:         b.GasLimit,
		MixHash:          b.MixHash,
		StateRoot:        b.StateRoot,
		Difficulty:       b.Difficulty,
		TotalDifficulty:  b.TotalDifficulty,
		Nonce:            b.Nonce,
		TransactionCount: b.TransactionCount,
		TransactionsRoot: b.TransactionsRoot,
		ReceiptsRoot:     b.ReceiptsRoot,
		LogsBloom:        b.LogsBloom,
		// Transactions:     txns,
		Logs: logs,
	}
}

/*
# Define your own custom webhook GraphQL query!
# For more example use cases & queries visit:
# https://docs.alchemy.com/reference/custom-webhooks-quickstart
#
{
  block {
    # Block hash is a great primary key to use for your data stores!
    hash,
    number,
    timestamp,
    parent{hash},
    baseFeePerGas,
    gasUsed,
    gasLimit,
    mixHash,
    stateRoot,
    totalDifficulty,
    # Add smart contract addresses to the list below to filter for specific logs
    logs(filter: {addresses: [], topics: []}) {
      data,
      topics,
      index,
      account {
        address
      },
      transaction {
        hash,
        nonce,
        index,
        from {
          address
        },
        to {
          address
        },
        value,
        gasPrice,
        maxFeePerGas,
        maxPriorityFeePerGas,
        gas,
        status,
        gasUsed,
        cumulativeGasUsed,
        effectiveGasPrice,
        createdContract {
          address
        }
      }
    }
  }
}
*/

/*
Headers
connection 	close
accept-encoding 	gzip,deflate
user-agent 	Apache-HttpClient/4.5.13 (Java/17.0.5)
host 	webhook.site
content-length 	322556
traceparent 	<string>
x-api-key 	<string>
x-alchemy-signature 	<string>
content-type 	application/json; charset=utf-8
*/
