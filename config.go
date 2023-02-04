package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	LogPath   string `json:"logPath"`
	BucketKey string `json:"bucketKey"`
}

func loadConfig(homeDir string, configDir string) Config {
	configFilePath := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		config := Config{
			LogPath:   filepath.Join(homeDir, "Documents", "Guild Wars 2", "addons", "arcdps", "arcdps.cbtlogs", "Ankka"),
			BucketKey: "",
		}

		data, err := json.Marshal(
			config,
		)
		if err != nil {
			log.Fatal(err)
		}

		_ = ioutil.WriteFile(configFilePath, data, 0644)
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	json.Unmarshal(byteValue, &config)

	return config
}
