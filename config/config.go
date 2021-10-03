// Package config provides function for loading config from environment variables.
package config

import (
	"fmt"
	"os"
)

var getenv = os.Getenv

type Config struct {
	BitFlyerAPIKey    string
	BitFlyerAPISecret string
}

type ErrRequired struct {
	envName string
}

func Load() (*Config, error) {
	var cfg Config

	if err := cfg.Reload(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg.BitFlyerAPIKey == "" {
		return ErrRequired{envName: "BITFLYER_API_KEY"}
	}

	if cfg.BitFlyerAPISecret == "" {
		return ErrRequired{envName: "BITFLYER_API_SECRET"}
	}

	return nil
}

func (cfg *Config) Reload() error {
	*cfg = Config{
		BitFlyerAPIKey:    getenv("BITFLYER_API_KEY"),
		BitFlyerAPISecret: getenv("BITFLYER_API_SECRET"),
	}

	return cfg.Validate()
}

func (err ErrRequired) Error() string {
	return fmt.Sprintf("%s is required", err.envName)
}
