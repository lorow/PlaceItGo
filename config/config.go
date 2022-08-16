package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
}

func GetConfig() (*Config, error) {
	config := Config{}

	err := envconfig.Process("", &config)
	return &config, err
}
