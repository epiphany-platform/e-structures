package save

import (
	awsbi "github.com/epiphany-platform/e-structures/awsbi/v0"
	"io/ioutil"

	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	azks "github.com/epiphany-platform/e-structures/azks/v0"
	hi "github.com/epiphany-platform/e-structures/hi/v0"
	st "github.com/epiphany-platform/e-structures/state/v0"
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

func AzBIConfig(path string, config *azbi.Config) error {
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

func AzKSConfig(path string, config *azks.Config) error {
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
func AwsBIConfig(path string, config *awsbi.Config) error {
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
