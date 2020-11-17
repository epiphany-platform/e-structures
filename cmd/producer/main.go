package main

import (
	"fmt"
	"github.com/mkyc/go-stucts-versioning-tests/pkg/azbi"
	"github.com/mkyc/go-stucts-versioning-tests/pkg/state"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	configFile = "azbi-config.json"
	stateFile  = "state.json"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	newConfig := initializeConfig()

	b, err := newConfig.Save()
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(pwd, configFile), b, 0644)
	if err != nil {
		log.Fatal(err)
	}

	newState := initializeState(newConfig)

	b, err = newState.Save()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(pwd, stateFile), b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func initializeState(config *azbi.Config) *state.State {
	s := state.NewState()
	s.AzBI.Config = config
	s.AzBI.Status = state.Initialized
	rgName := fmt.Sprintf("%s-rg", *config.Params.Name)
	vnetName := fmt.Sprintf("%s-vnet", *config.Params.Name)
	s.AzBI.Output = &azbi.Output{
		PrivateIps: []string{},
		PublicIps:  []string{"123.234.345.456", "234.345.456.567", "345.456.567.678"},
		RgName:     &rgName,
		VmNames:    []string{"vm1", "vm2", "vm3"},
		VnetName:   &vnetName,
	}
	return s
}

func initializeConfig() *azbi.Config {
	return azbi.NewConfig()
}
