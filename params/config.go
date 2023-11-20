package params

import (
	"github.com/iquidus/blockspider/crawler"
	"github.com/iquidus/blockspider/rpc"
	"github.com/iquidus/blockspider/state"
)

type Config struct {
	ChainId uint64         `json:"chainId"`
	Crawler crawler.Config `json:"crawler"`
	Rpc     rpc.Config     `json:"rpc"`
	State   state.Config   `json:"state"`
}
