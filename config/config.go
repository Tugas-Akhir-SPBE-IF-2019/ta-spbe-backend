package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
)

type API struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type Postgres struct {
	Host      string `toml:"host"`
	Port      int    `toml:"port"`
	Database  string `toml:"database"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	Migration bool   `toml:"migration"`
}

type OAuth struct {
	ClientId     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RedirectURL  string `toml:"redirect_uri"`
}

type JWT struct {
	Secret string `toml:"secret"`
}

type Config struct {
	API   API      `toml:"api"`
	DB    Postgres `toml:"postgres"`
	OAuth OAuth    `toml:"oauth"`
	JWT   JWT      `toml:"jwt"`
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
