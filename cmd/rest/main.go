package main

import (
	"flag"
	"fmt"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest"

	"log"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/logger"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/pgsql"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/tracer"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// LOAD APPLICATION CONFIG FROM ENVIRONMENT VARIABLES
	// -----------------------------------------------------------------------------------------------------------------
	cfgPath := flag.String("c", "config.toml", "path to config file")
	cfg, err := config.LoadEnvFromFile(*cfgPath)
	if err != nil {
		log.Fatalln(err)
	}

	// -----------------------------------------------------------------------------------------------------------------
	// STRUCTURED LOGGER
	// -----------------------------------------------------------------------------------------------------------------
	zlogger := logger.New(cfg.Logger).With().
		Logger()

	// -----------------------------------------------------------------------------------------------------------------
	// SET OPEN TELEMETRY GLOBAL TRACER
	// -----------------------------------------------------------------------------------------------------------------
	if err = tracer.SetTracer(cfg.Tracer, cfg.AppInfo); err != nil {
		zlogger.Error().Err(err).Msgf("grpc: main failed to setup open telemetry tracer: %s", err)
		return
	}

	// -----------------------------------------------------------------------------------------------------------------
	// INFRASTRUCTURE OBJECTS
	// -----------------------------------------------------------------------------------------------------------------
	// PGSQL
	sqlDB, sqlDBErr := pgsql.NewDB(cfg.PostgreSQL, zlogger)
	if sqlDBErr != nil {
		zlogger.Error().Err(sqlDBErr).Msgf("grpc: main failed to construct pgsql: %s", sqlDBErr)
		return
	}

	migrate := flag.Bool("migrate", cfg.PostgreSQL.Migration, "do migration")
	if *migrate {
		if migrateErr := pgsql.Migrate(sqlDB, cfg.PostgreSQL.Database); err != nil {
			zlogger.Error().Err(migrateErr).Msgf("migration: migration failed to construct pgsql: %s", migrateErr)
			return
		}
	}

	// -----------------------------------------------------------------------------------------------------------------
	// SERVER SETUP AND EXECUTE
	// -----------------------------------------------------------------------------------------------------------------
	restServerHandler := rest.New(cfg, zlogger, sqlDB)

	zlogger.Info().Msgf("REST Server started on port %d", cfg.API.RESTPort)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.API.RESTPort), restServerHandler)
}
