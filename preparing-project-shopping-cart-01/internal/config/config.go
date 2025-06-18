package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerAddress string
}

func NewConfig() *Config {
	return &Config{
		ServerAddress: fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
	}
}
