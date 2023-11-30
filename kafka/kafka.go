package kafka

import "github.com/iquidus/blockspider/common"

type EventsConfig struct {
	Broker    string   `json:"broker"`
	Topic     string   `json:"topic"`
	Addresses []string `json:"addresses"`
	Topics    []string `json:"topics"`
}

type BlocksConfig struct {
	Broker string `json:"broker"`
	Topic  string `json:"topic"`
}

type Config struct {
	Blocks BlocksConfig   `json:"blocks"`
	Events []EventsConfig `json:"events"`
}

type BlocksPayload struct {
	Method string       `json:"method"`
	Block  common.Block `json:"block"`
}

type EventsPayload struct {
	Method string         `json:"method"`
	Events []common.Log `json:"events"`
}
