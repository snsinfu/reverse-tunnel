package config

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Load(r io.Reader, conf interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, conf)
}
