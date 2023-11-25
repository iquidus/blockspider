package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
	"github.com/iquidus/blockspider/alchemy"
	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/kafka"
	"github.com/iquidus/blockspider/params"
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
	configFlagDesc    = "specify name of config file"

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

func readConfig(cfg *params.Config) {
	if configFileName == "" {
		mainLogger.Error("Invalid arguments", os.Args)
		os.Exit(1)
	}

	confPath, err := filepath.Abs(configFileName)
	if err != nil {
		mainLogger.Error("Error: could not parse config filepath", "err", err)
	}

	mainLogger.Info("Loading config", "path", confPath)

	configFile, err := os.Open(confPath)
	if err != nil {
		appLogger.Error("File error", "err", err.Error())
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		mainLogger.Error("Config error", "err", err.Error())
	}
}

// https://docs.alchemy.com/reference/custom-webhooks-faq
func isValidSignatureForStringBody(body []byte, signature string, signingKey []byte) bool {
	h := hmac.New(sha256.New, signingKey)
	h.Write(body)
	digest := hex.EncodeToString(h.Sum(nil))
	return digest == signature
}

func sendBlockMessage(blockWriter *kafka.Writer, block common.Block) error {
	var bp = kafka.BlocksPayload{
		Method: "PUSH",
		Block:  block,
	}

	payload, err := json.Marshal(bp)
	if err != nil {
		return err
	}

	err = blockWriter.WriteMessages(context.Background(), payload)

	if err != nil {
		log.Error("failed to write messages", "err", err)
	}

	return nil
}

func setupRouter(blockWriter *kafka.Writer, cfg params.TransmuteConfig) *gin.Engine {
	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies(cfg.TrustedProxies)

	// Define alchemy endpoint
	r.GET("/alchemy", func(c *gin.Context) {
		signature := c.Params.ByName("x-alchemy-signature")
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validate signature
		if !isValidSignatureForStringBody(body, signature, []byte(cfg.Alchemy.Secret)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}

		// Parse JSON
		event := new(alchemy.AlchemyEvent)
		if err := json.Unmarshal(body, event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Handle block event
		// convert to common block
		block := event.Data.Block.Convert()
		// send block to kafka
		err = sendBlockMessage(blockWriter, block)
		if err != nil {
			log.Info("failed to write messages", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}

func main() {
	log.Info("blockspider/transmuted ", "version", params.VersionWithMeta)
	readConfig(&cfg)
	// Create blockwriter
	bw := kafka.NewWriter(cfg.Crawler.Kafka.Blocks.Broker, &cfg.Crawler.Kafka.Blocks.Topic, 1)
	// Init gin router
	r := setupRouter(bw, cfg.Transmute)
	// Listen and Server
	r.Run(fmt.Sprintf(":%d", cfg.Transmute.Port))
}
