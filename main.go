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

type Encounter struct {
	Duration      int64     `json:"durationMS"`
	TimeStart     string    `json:"timeStart"`
	EiEncounterID int64     `json:"eiEncounterID"`
	TriggerID     int64     `json:"triggerID"`
	Players       []Players `json:"players"`
}

type Players struct {
	Account    string `json:"Account"`
	Profession string `json:"profession"`
	Name       string `json:"name"`
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
			fmt.Printf("Parsing Log: %s.\n", path)

			v := strings.Split(info.Name(), ".")
			name := v[0]

			cmd := exec.Command(EIPath, "-c", EIConfigPath, "-p", path)
			err := cmd.Run()
			if err != nil {
				log.Panic(err)
			}

			filepath.Walk(tempDirPath, func(parsedLogPath string, parsedLogInfo os.FileInfo, err error) error {
				if strings.HasPrefix(parsedLogInfo.Name(), name) {
					parsedLogFile, err := os.Open(parsedLogPath)
					if err != nil {
						log.Fatal(err)
					}
					defer configFile.Close()

					byteValue, err := ioutil.ReadAll(parsedLogFile)
					if err != nil {
						log.Fatal(err)
					}

					var encounter Encounter
					json.Unmarshal(byteValue, &encounter)

					fmt.Printf("%+v\n", encounter)

					fmt.Printf("Removing parsed log: %s.\n", parsedLogPath)
					os.Remove(parsedLogPath)
				}

				return nil
			})
		}

		return nil
	})

	fmt.Print("Removing temp directory.")
	os.RemoveAll(tempDirPath)
}
