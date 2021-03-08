package load

import (
	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	azks "github.com/epiphany-platform/e-structures/azks/v0"
	hi "github.com/epiphany-platform/e-structures/hi/v0"
	"io/ioutil"
	"os"

	st "github.com/epiphany-platform/e-structures/state/v0"
)

// TODO handle problematic situation when some of config structures are "almost" empty and validation fails (https://github.com/epiphany-platform/e-structures/issues/10)
func State(path string) (*st.State, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return st.NewState(), nil
	} else {
		state := &st.State{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = state.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		return state, nil
	}
}

func AzBIConfig(path string) (*azbi.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return azbi.NewConfig(), nil
	} else {
		config := &azbi.Config{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = config.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
}

func AzKSConfig(path string) (*azks.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return azks.NewConfig(), nil
	} else {
		config := &azks.Config{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = config.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
}

func HiConfig(path string) (*hi.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return hi.NewConfig(), nil
	} else {
		config := &hi.Config{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = config.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
}
