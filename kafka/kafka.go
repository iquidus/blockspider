package kafka

import "github.com/iquidus/blockspider/common"

type TopicParams struct {
	Topic     string   `json:"topic"`
	Addresses []string `json:"addresses"`
	Topics    []string `json:"topics"`
}

type Config struct {
	Broker string        `json:"broker"`
	Params []TopicParams `json:"params"`
}

type Payload struct {
	Status  string       `json:"status"`
	Block   common.Block `json:"block"`
	Version int          `json:"version"`
}
