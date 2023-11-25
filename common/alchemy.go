package common

type AlchemyWebhookBlockLogs struct {
	Data        string                    `bson:"data" json:"data"`
	Topics      []string                  `bson:"topics" json:"topics"`
	Index       uint64                    `bson:"index" json:"index"`
	Account     AlchemyWebhookAccount     `bson:"account" json:"account"`
	Transaction AlchemyWebhookTransaction `bson:"transaction" json:"transaction"`
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

type AlchemyWebhookAccount struct {
	Address string `bson:"address" json:"address"`
}

type AlchemyWebhookParentBlock struct {
	Hash string `bson:"hash" json:"hash"`
}

type AlchemyWebhookBlock struct {
	Hash            string                    `bson:"hash" json:"hash"`
	Number          uint64                    `bson:"number" json:"number"`
	Timestamp       uint64                    `bson:"timestamp" json:"timestamp"`
	Parent          AlchemyWebhookParentBlock `bson:"parent" json:"parent"`
	BaseFeePerGas   string                    `bson:"baseFeePerGas" json:"baseFeePerGas,omitempty"`
	GasUsed         uint64                    `bson:"gasUsed" json:"gasUsed"`
	GasLimit        uint64                    `bson:"gasLimit" json:"gasLimit"`
	MixHash         string                    `bson:"mixHash" json:"mixHash"`
	StateRoot       string                    `bson:"stateRoot" json:"stateRoot"`
	TotalDifficulty string                    `bson:"totalDifficulty" json:"totalDifficulty"`
	Logs            []AlchemyWebhookBlockLogs `bson:"logs" json:"logs"`
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
