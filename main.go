package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	// Ensure config dir exists
	configDir := filepath.Join(homeDir, ".arcdps-log-uploader")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0750)
	}

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

	EIPath := filepath.Join(homeDir, "GW2EI", "GuildWars2EliteInsights.exe")
	EIConfigPath := filepath.Join(homeDir, ".arcdps-log-uploader", "elite_insights.conf")

	tempDirPath := filepath.Join(configDir, "temp")
	os.Mkdir(tempDirPath, 0750)

	filepath.Walk(config.LogPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".zevtc") {
			fmt.Printf("%v\n", path)
			fmt.Printf("%v\n", info.Name())

			cmd := exec.Command(EIPath, "-c", EIConfigPath, "-p", path)
			err := cmd.Run()
			if err != nil {
				log.Panic(err)
			}

			// // Unzip compressed evtc file
			// r, err := zip.OpenReader(path)
			// if err != nil {
			// 	log.Panic(err)
			// }
			// defer r.Close()

			// for _, f := range r.File {
			// 	rc, err := f.Open()
			// 	if err != nil {
			// 		log.Panic(err)
			// 	}
			// 	defer rc.Close()

			// 	Parse evtc file
			// 	header, _, _, err := evtc.ParseHeader(rc)
			// 	if err != nil {
			// 		log.Panic(err)
			// 	}
			// 	fmt.Printf("%v %d\n", string(header.Date[:]), header.Boss)

			// 	// fmt.Printf("%v\n", chain)
			// }
		}

		return nil
	})

	os.RemoveAll(tempDirPath)
}
