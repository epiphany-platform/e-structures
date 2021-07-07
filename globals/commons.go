package globals

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Backup(i interface{}, new string) error {
	if _, err := os.Stat(new); os.IsExist(err) {
		return err
	}
	bytes, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(new, bytes, 0644)
}

func Save(p Printer, path string) error {
	bytes, err := p.Print()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

func Print(v Validator) ([]byte, error) {
	if err := v.Valid(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(v, "", "\t")
}
