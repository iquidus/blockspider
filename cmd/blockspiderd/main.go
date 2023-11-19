package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/iquidus/blockspider/params"
	"github.com/iquidus/blockspider/rpc"
	"github.com/iquidus/blockspider/state"
)

var (
	cfg        params.Config
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

	readConfig(&cfg)

	mainLogger.Debug("printing config", "cfg", cfg)

	rpcClient := rpc.NewRPCClient(&cfg.Rpc)
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

	startBlock, err := rpcClient.GetBlockByHeight(cfg.Crawler.Start)
	if err != nil {
		log.Error("could not retrieve start block", "err", err)
		os.Exit(1)
	}

	state, err := state.Init(&cfg.State, &cfg.ChainId, startBlock)
	if err != nil {
		log.Error("could not initialize state", "err", err)
		os.Exit(1)
	}
	
	// TODO(iquidus): init kafka here, check for topics, create if they dont exist.
	go startCrawler(&cfg.Crawler, state, rpcClient, appLogger)
	
	quit := make(chan int)
	<-quit
}
