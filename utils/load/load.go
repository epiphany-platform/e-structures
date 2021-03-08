package load

import (
	hi "github.com/epiphany-platform/e-structures/hi/v0"
	"io/ioutil"
	"os"

	st "github.com/epiphany-platform/e-structures/state/v0"
)

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
