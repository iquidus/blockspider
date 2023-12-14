package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/log"
	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/crawler"
	"github.com/iquidus/blockspider/disk"
	"github.com/iquidus/blockspider/kafka"
	"github.com/iquidus/blockspider/params"
	"github.com/iquidus/blockspider/state"
)

var (
	appLogger  = log.Root()
	mainLogger log.Logger

	RootHandler *log.GlogHandler

	logLevel       string
	configFileName string
)

const (
	configFlagDefault = "config.json"
	configFlagDesc    = "specify name of config file (should be in working dir)"

	logLevelFlagDefault = "info"
	logLevelFlagDesc    = "set level of logs"
)

func init() {
	flag.StringVar(&configFileName, "c", configFlagDefault, configFlagDesc)
	flag.StringVar(&configFileName, "config", configFlagDefault, configFlagDesc)

	flag.StringVar(&logLevel, "ll", logLevelFlagDefault, logLevelFlagDesc)
	flag.StringVar(&logLevel, "logLevel", logLevelFlagDefault, logLevelFlagDesc)

	flag.Parse()

	RootHandler = log.NewGlogHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(true)))

	if logLevel == "debug" || logLevel == "d" || logLevel == "dbg" {
		RootHandler.Verbosity(log.LvlDebug)
	} else if logLevel == "trace" || logLevel == "t" {
		RootHandler.Verbosity(log.LvlTrace)
	} else {
		RootHandler.Verbosity(log.LvlInfo)
	}

	appLogger.SetHandler(RootHandler)

	mainLogger = log.Root().New()
}

func main() {
	log.Info(fmt.Sprint("blockspider ", params.VersionWithMeta))
	// Read config
	var cfg params.Config
	configPath, err := filepath.Abs(configFileName)
	if err != nil {
		log.Error("Error: could not parse config filepath", "err", err)
		os.Exit(1)
	}
	err = disk.ReadJsonFile[params.Config](configPath, &cfg)
	if err != nil {
		log.Error("Error: could read config file", "err", err)
		os.Exit(1)
	}

	mainLogger.Debug("printing config", "cfg", cfg)

	// check node connection
	rpcClient := common.NewRPCClient(&cfg.Rpc)
	version, err := rpcClient.Ping()
	if err != nil {
		switch err.(type) {
		case *url.Error:
			mainLogger.Error("rpc node offline", "err", err)
			os.Exit(1)
		default:
			mainLogger.Error(fmt.Sprintf("error pinging rpc node (%T)", err), "err", err)
		}
	}

	mainLogger.Info("connected to rpc server", "version", version)

	// Initialize state
	s, err := state.Init(&cfg.State, &cfg.ChainId)
	if err != nil {
		log.Error("could not initialize state", "err", err)
		os.Exit(1)
	}
	// Check if cache is empty
	if s.Cache.Count() == 0 {
		// empty cache, use start block
		log.Info("cache is empty, using start block", "number", cfg.Crawler.Start)
		// get start block from rpc
		rawStartBlock, err := rpcClient.GetBlockByHeight(cfg.Crawler.Start)
		if err != nil {
			log.Error("could not retrieve start block", "err", err)
			os.Exit(1)
		}
		if rawStartBlock.Hash == "" {
			// empty block with no err, possible future blocknumber, abort
			err = errors.New("block not found")
			log.Error("could not retrieve start block", "err", err)
			os.Exit(1)
		}
		// convert to common block
		startBlock, err := rawStartBlock.Convert(rpcClient, nil)
		if err != nil {
			log.Error("could not convert start block", "err", err)
			os.Exit(1)
		}
		// push to cache/state
		s.Cache.Push(startBlock)
	} else {
		// cache is not empty, resume from cached head
		cachedHead, err := s.Cache.Peak()
		if err != nil {
			log.Error("could not retrieve cached head", "err", err)
			os.Exit(1)
		}
		log.Info("resuming from cached block", "number", cachedHead.Number, "hash", cachedHead.Hash)
	}

	// Create kafka writer
	kw := kafka.NewWriter(cfg.Kafka.Broker, nil, 1)

	// Start crawler
	go startCrawler(&cfg.Crawler, s, rpcClient, kw, appLogger)

	quit := make(chan int)
	<-quit
}

func startCrawler(cfg *crawler.Config, s *state.State, rpc *common.RPCClient, writer *kafka.Writer, logger log.Logger) {
	blockCrawler := crawler.NewCrawler(cfg, s, rpc, writer, logger.New())
	logger.Info("Starting crawler")
	crawler.Start(blockCrawler, cfg, logger)
}
