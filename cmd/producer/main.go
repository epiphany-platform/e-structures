package main

import (
	"github.com/mkyc/go-stucts-versioning-tests/pkg/azbi"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	configName = "azbi-config.json"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile(filepath.Join(pwd, configName))
	if err != nil {
		return
	}

	existingConfig, err := azbi.Load(b)
	if err != nil {
		log.Fatal(err)
	}
	if existingConfig.Unused != nil && len(existingConfig.Unused) > 0 {
		for _, u := range existingConfig.Unused {
			log.Println(u)
		}
	}

	upgradedConfig, err := performBusinessLogic(*existingConfig)
	if err != nil {
		log.Fatal(err)
	}

	b, err = upgradedConfig.Save()
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(pwd, configName), b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func performBusinessLogic(config azbi.Config) (azbi.Config, error) {
	return config, nil
}
