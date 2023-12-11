package crawler

import (
	"os"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/iquidus/blockspider/cache"
	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/kafka"
	"github.com/iquidus/blockspider/state"
)

type Config struct {
	Interval    string       `json:"interval"`
	MaxRoutines int          `json:"routines"`
	CacheLimit  int          `json:"cache"`
	Start       uint64       `json:"start"`
	Kafka       kafka.Config `json:"kafka"`
}

type Crawler struct {
	// backend *storage.MongoDB
	rpc     *common.RPCClient
	cfg     *Config
	logChan chan *logObject
	state   *state.State
	cache   *cache.BlockStack[common.Block]
	logger  log.Logger
	writer  *kafka.Writer
}

func NewCrawler(cfg *Config, state *state.State, rpc *common.RPCClient, logger log.Logger) *Crawler {
	bc := cache.New[common.Block](&cfg.CacheLimit)
	// create kafka writer
	kw := kafka.NewWriter(cfg.Kafka.Broker, nil, 1)

	return &Crawler{rpc, cfg, make(chan *logObject), state, bc, logger, kw}
}

func runCrawler(ticker *time.Ticker, c Crawler) {
	c.RunLoop()
	for {
		<-ticker.C
		c.RunLoop()
	}
}

func Start(crawler *Crawler, cfg *Config, logger log.Logger) {
	blockInterval, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		logger.Error("can't parse crawler duration", "d", cfg.Interval, "err", err)
		os.Exit(1)
	}
	blockTicker := time.NewTicker(blockInterval)
	logger.Info("Crawler interval set", "d", cfg.Interval)
	go runCrawler(blockTicker, *crawler)
}
