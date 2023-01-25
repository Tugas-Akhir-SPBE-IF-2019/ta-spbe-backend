package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
)

type API struct {
	Port int `toml:"port"`
}

type Postgres struct {
	Host      string `toml:"host"`
	Port      int    `toml:"port"`
	Database  string `toml:"database"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	Migration bool   `toml:"migration"`
}

type Config struct {
	API API      `toml:"api"`
	DB  Postgres `toml:"postgres"`
}

func LoadEnvFromFile(path string) (Config, error) {
	cfg := Config{}

	file, err := os.Open(path)
	if err != nil {
		return cfg, fmt.Errorf("error open config file: %w", err)
	}
	err = toml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("error parsing toml: %w", err)
	}

	return cfg, nil
}
