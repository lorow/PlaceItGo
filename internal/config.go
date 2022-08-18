package internal

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	redis_url      string
	redis_database int
	redis_password string
}

func GetConfig() (*Config, error) {
	config := Config{}

	err := envconfig.Process("", &config)
	return &config, err
}
