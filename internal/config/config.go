package config

import (
	"fmt"
	"os"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/appinfo"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/logger"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/pgsql"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/token"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/tracer"
	"github.com/pelletier/go-toml"
)

type API struct {
	Host     string `toml:"host"`
	RESTPort int    `toml:"rest_port"`
}

type Config struct {
	API        API           `toml:"api"`
	AppInfo    appinfo.Info  `toml:"app_info"`
	Logger     logger.Config `toml:"logger"`
	PostgreSQL pgsql.Config  `toml:"postgres"`
	Tracer     tracer.Config `toml:"tracer"`
	JWT        token.Config  `toml:"jwt"`
}

func LoadEnvFromFile(path string) (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return cfg, fmt.Errorf("error open config file: %w", err)
	}
	err = toml.NewDecoder(file).Decode(cfg)
	if err != nil {
		return cfg, fmt.Errorf("error parsing toml: %w", err)
	}

	return cfg, nil
}
