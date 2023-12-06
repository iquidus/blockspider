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
	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/disk"
	"github.com/iquidus/blockspider/kafka"
	"github.com/iquidus/blockspider/params"
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

// https://docs.alchemy.com/reference/custom-webhooks-faq
func isValidSignatureForStringBody(body []byte, signature string, signingKey []byte) bool {
	h := hmac.New(sha256.New, signingKey)
	h.Write(body)
	digest := hex.EncodeToString(h.Sum(nil))
	return digest == signature
}

func sendBlockMessage(blockWriter *kafka.Writer, block common.Block) error {
	var bp = kafka.Payload{
		Status:  "ACCEPTED",
		Block:   block,
		Version: 1,
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
		event := new(common.AlchemyEvent)
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
	// Read config
	var cfg params.Config
	configPath, err := filepath.Abs(configFileName)
	if err != nil {
		mainLogger.Error("Error: could not parse config filepath", "err", err)
	}
	err = disk.ReadJsonFile[params.Config](configPath, &cfg)
	if err != nil {
		log.Error("Error: could read config file", "err", err)
	}
	// Create blockwriter
	kw := kafka.NewWriter(cfg.Crawler.Kafka.Broker, nil, 1)
	// Init gin router
	r := setupRouter(kw, cfg.Transmute)
	// Listen and Server
	r.Run(fmt.Sprintf(":%d", cfg.Transmute.Port))
}
