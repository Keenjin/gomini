package ktoml

import (
	"github.com/BurntSushi/toml"
	"os"
)

func LoadToml(cfgPath string, v interface{}) error {
	_, err := toml.DecodeFile(cfgPath, v)
	return err
}

func SaveToml(cfgPath string, v interface{}) error {
	f, err := os.Create(cfgPath)
	if err != nil {
		return err
	}
	defer f.Close()

	cfg := toml.NewEncoder(f)
	return cfg.Encode(v)
}
