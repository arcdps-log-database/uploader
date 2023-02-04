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

func parseLog(EIPath string, EIConfigPath string, logPath string) error {
	cmd := exec.Command(EIPath, "-c", EIConfigPath, "-p", logPath)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

type Encounter struct {
	Duration      int64    `json:"durationMS"`
	TimeStart     string   `json:"timeStart"`
	EiEncounterID int64    `json:"eiEncounterID"`
	TriggerID     int64    `json:"triggerID"`
	Players       []Player `json:"players"`
}

type Player struct {
	Account    string `json:"Account"`
	Profession string `json:"profession"`
	Name       string `json:"name"`
}

func readParsedLog(tempDirPath string, name string) (Encounter, error) {
	var encounter Encounter

	filepath.Walk(tempDirPath, func(parsedLogPath string, parsedLogInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(parsedLogInfo.Name(), name) {
			parsedLogFile, err := os.Open(parsedLogPath)
			if err != nil {
				log.Fatal(err)
			}
			defer parsedLogFile.Close()

			byteValue, err := ioutil.ReadAll(parsedLogFile)
			if err != nil {
				log.Fatal(err)
			}

			json.Unmarshal(byteValue, &encounter)

			fmt.Printf("Removing parsed log: %s.\n", parsedLogPath)
			os.Remove(parsedLogPath)
		}

		return nil
	})

	return encounter, nil
}
