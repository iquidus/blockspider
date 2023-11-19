package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/iquidus/blockspider/params"
)

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
