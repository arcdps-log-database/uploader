package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Yi-Jiahe/evtc"
)

type Config struct {
	LogPath   string `json:"logPath"`
	BucketKey string `json:"bucketKey"`
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configDir := filepath.Join(homeDir, ".arcdps-log-uploader")
	configFilePath := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0750)

		config := Config{
			LogPath:   filepath.Join(homeDir, `Documents\Guild Wars 2\addons\arcdps\arcdps.cbtlogs`),
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

	filepath.Walk(config.LogPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".zevtc") {
			fmt.Printf("%v\n", info.Name())

			// Unzip compressed evtc file
			r, err := zip.OpenReader(path)
			if err != nil {
				log.Panic(err)
			}
			defer r.Close()

			for _, f := range r.File {
				rc, err := f.Open()
				if err != nil {
					log.Panic(err)
				}
				defer rc.Close()

				fmt.Printf("Parsing log...\n")

				// Parse evtc file
				_, err = evtc.Parse(rc)
				fmt.Printf("Log parsed.\n")
				if err != nil {
					fmt.Printf("Error occurred when parsing log.\n")
					// log.Panic(err)
				}

				// fmt.Printf("%v\n", chain)
			}
		}

		return nil
	})
}
