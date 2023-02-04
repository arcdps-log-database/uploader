package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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

	config := loadConfig(homeDir, configDir)

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

			err := parseLog(EIPath, EIConfigPath, path)
			if err != nil {
				log.Panic(err)
			}

			v := strings.Split(info.Name(), ".")
			name := v[0]

			encounter, err := readParsedLog(tempDirPath, name)
			if err != nil {
				log.Panic(err)
			}
			fmt.Printf("%+v\n", encounter)
		}

		return nil
	})

	fmt.Print("Removing temp directory.")
	os.RemoveAll(tempDirPath)
}
