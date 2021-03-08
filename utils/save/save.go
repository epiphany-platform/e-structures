package save

import (
	hi "github.com/epiphany-platform/e-structures/hi/v0"
	st "github.com/epiphany-platform/e-structures/state/v0"
	"io/ioutil"
)

func State(path string, state *st.State) error {
	bytes, err := state.Marshal()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func HiConfig(path string, config *hi.Config) error {
	bytes, err := config.Marshal()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
