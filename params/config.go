package params

import (
	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/crawler"
	"github.com/iquidus/blockspider/kafka"
	"github.com/iquidus/blockspider/state"
)

type TransmuteConfig struct {
	Port           uint64               `json:"port"`
	TrustedProxies []string             `json:"trustedProxies"`
	Alchemy        common.AlchemyConfig `json:"alchemy"`
}

type Config struct {
	ChainId   uint64           `json:"chainId"`
	Crawler   crawler.Config   `json:"crawler"`
	Rpc       common.RPCConfig `json:"rpc"`
	State     state.Config     `json:"state"`
	Kafka     kafka.Config     `json:"kafka"`
	Transmute TransmuteConfig  `json:"transmute"`
}
