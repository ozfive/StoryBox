package utils

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Debug        bool   `yaml:"debug"`
	Address      string `yaml:"address"`
	Port         string `yaml:"port"`
	LogFilePath  string `yaml:"log_file_path"`
	DatabasePath string `yaml:"database_path"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
