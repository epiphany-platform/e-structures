package imh

import (
	"fmt"
	"github.com/epiphany-platform/e-structures/globals"
	"os"
	"path/filepath"
	"time"
)

type Modulator interface {
	globals.Initializer
	globals.Backuper
	globals.Loader
	globals.Saver
	globals.Printer
	globals.Validator
	globals.Upgrader
}

const (
	configFileName      = "config.json"
	stateFileName       = "state.json"
	backupDirectoryName = "backup"
)

type InfrastructureModuleHelper struct {
	ModuleDirectoryPath string
	ModuleVersion       string
}

func (h InfrastructureModuleHelper) Initialize(config Modulator, state Modulator) (Modulator, Modulator, error) {
	// check if required fields are set
	if h.ModuleDirectoryPath == "" {
		return nil, nil, fmt.Errorf("setup module directory path first")
	}
	if h.ModuleVersion == "" {
		return nil, nil, fmt.Errorf("setup module version first")
	}

	// ensure module directory
	err := os.MkdirAll(h.ModuleDirectoryPath, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	// create default files paths
	configFilePath := filepath.Join(h.ModuleDirectoryPath, configFileName)
	stateFilePath := filepath.Join(h.ModuleDirectoryPath, stateFileName)

	// load state file
	err = state.Load(stateFilePath)
	if os.IsNotExist(err) {
		// if no state loaded then init it
		state.Init(h.ModuleVersion)
	} else if err != nil {
		return nil, nil, fmt.Errorf("load config failed: %v", err)
	}
	// load config file
	err = config.Load(configFilePath)
	if os.IsNotExist(err) {
		// if no config loaded then init it
		config.Init(h.ModuleVersion)
	} else if err != nil {
		return nil, nil, fmt.Errorf("load state failed: %v", err)
	}

	// backup
	err = backup(h.ModuleDirectoryPath, config, state)
	if err != nil {
		return nil, nil, err
	}

	return config, state, nil
}

func (h InfrastructureModuleHelper) Load(config Modulator, state Modulator) (Modulator, Modulator, error) {
	// check if required fields are set
	if h.ModuleDirectoryPath == "" {
		return nil, nil, fmt.Errorf("setup module directory path first")
	}
	if h.ModuleVersion == "" {
		return nil, nil, fmt.Errorf("setup module version first")
	}

	// ensure module directory
	err := os.MkdirAll(h.ModuleDirectoryPath, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	// create default files paths
	configFilePath := filepath.Join(h.ModuleDirectoryPath, configFileName)
	stateFilePath := filepath.Join(h.ModuleDirectoryPath, stateFileName)

	// load state file
	err = state.Load(stateFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("load config failed: %v", err)
	}
	// load config file
	err = config.Load(configFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("load state failed: %v", err)
	}

	// backup
	err = backup(h.ModuleDirectoryPath, config, state)
	if err != nil {
		return nil, nil, err
	}

	return config, state, nil
}

func backup(moduleDirectoryPath string, config Modulator, state Modulator) error {
	backupDirectoryPath := filepath.Join(moduleDirectoryPath, backupDirectoryName)

	// ensure backups directory
	err := os.MkdirAll(backupDirectoryPath, os.ModePerm)
	if err != nil {
		return err
	}

	// prepare timestamp
	t := time.Now().Format("20060102-150405.999999")
	// backup config to
	err = config.Backup(filepath.Join(backupDirectoryPath, fmt.Sprintf("config-%s.json", t)))
	if err != nil {
		return fmt.Errorf("config backup failed: %v", err)
	}
	err = state.Backup(filepath.Join(backupDirectoryPath, fmt.Sprintf("state-%s.json", t)))
	if err != nil {
		return fmt.Errorf("state backup failed: %v", err)
	}

	return nil
}
