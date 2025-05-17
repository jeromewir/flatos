package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	CorcoranAPIKey string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() error {
	c.CorcoranAPIKey = os.Getenv("CORCORAN_API_KEY")

	if c.CorcoranAPIKey == "" {
		return fmt.Errorf("CORCORAN_API_KEY is not set")
	}

	return nil
}
