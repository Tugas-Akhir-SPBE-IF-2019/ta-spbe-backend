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

type SMTPClient struct {
	Debug         bool   `toml:"debug"`
	Host          string `toml:"host"`
	Port          int    `toml:"port"`
	AdminIdentity string `toml:"admin_identity"`
	AdminEmail    string `toml:"admin_email"`
	AdminPassword string `toml:"admin_password"`
}

type MessageBroker struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type Config struct {
	API           API           `toml:"api"`
	DB            Postgres      `toml:"postgres"`
	OAuth         OAuth         `toml:"oauth"`
	JWT           JWT           `toml:"jwt"`
	SMTPClient    SMTPClient    `toml:"smtp"`
	MessageBroker MessageBroker `toml:"messagebroker"`
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
