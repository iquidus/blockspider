package main

import (
	"github.com/iquidus/blockspider/crawler"
	"github.com/iquidus/blockspider/rpc"
	"github.com/iquidus/blockspider/state"

	"github.com/ethereum/go-ethereum/log"
)

func startCrawler(cfg *crawler.Config, s *state.State, rpc *rpc.RPCClient, logger log.Logger) {
	blockCrawler := crawler.NewCrawler(cfg, s, rpc, logger.New())
	logger.Info("Starting crawler")
	crawler.Start(blockCrawler, cfg, logger)
}
