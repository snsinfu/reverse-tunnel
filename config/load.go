package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// Load loads a YAML document from a file and fills in corresponding fields of
// conf variable. Default values given in the variable are honored.
func Load(path string, conf interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, conf)
}
