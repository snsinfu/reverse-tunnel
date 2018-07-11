package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Load loads config variables from a YAML file. Default values given in conf
// variable are honored.
func Load(path string, conf interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, conf)
}
