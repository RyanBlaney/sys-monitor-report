package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogInterval int `yaml:"log_interval"`
	Thresholds  struct {
		CPU     int `yaml:"cpu"`
		Memory  int `yaml:"memory"`
		Disk    int `yaml:"disk"`
		Network int `yaml:"network"`
	}
}

func LoadConfig(path string) (Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}
