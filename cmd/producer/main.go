package main

import (
	"encoding/json"
	"github.com/mkyc/go-stucts-versioning-tests/pkg/azbi"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	azbiConfigName = "azbi-config.json"
)

func main() {

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile(filepath.Join(pwd, azbiConfigName))
	if err != nil {
		log.Fatal(err)
	}
	var existingConfig azbi.Config
	err = json.Unmarshal(b, &existingConfig)
	if err != nil {
		log.Fatal(err)
	}

	upgradedConfig, err := performBusinessLogic(existingConfig)

	b2, err := json.MarshalIndent(upgradedConfig, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(pwd, azbiConfigName), b2, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func performBusinessLogic(config azbi.Config) (azbi.Config, error) {
	return config, nil
}
