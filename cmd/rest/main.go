package main

import (
	"flag"
	"fmt"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest"
	"github.com/nsqio/go-nsq"

	"log"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/filesystem"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/jsonmanipulator"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/logger"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/messagequeue"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/pgsql"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/smtpmailer"
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
		zlogger.Error().Err(err).Msgf("rest: main failed to setup open telemetry tracer: %s", err)
		return
	}

	// -----------------------------------------------------------------------------------------------------------------
	// INFRASTRUCTURE OBJECTS
	// -----------------------------------------------------------------------------------------------------------------
	// PGSQL
	sqlDB, sqlDBErr := pgsql.NewDB(cfg.PostgreSQL, zlogger)
	if sqlDBErr != nil {
		zlogger.Error().Err(sqlDBErr).Msgf("rest: main failed to construct pgsql: %s", sqlDBErr)
		return
	}

	migrate := flag.Bool("migrate", cfg.PostgreSQL.Migration, "do migration")
	if *migrate {
		if migrateErr := pgsql.Migrate(sqlDB, cfg.PostgreSQL.Database); err != nil {
			zlogger.Error().Err(migrateErr).Msgf("rest: migration failed to migrate: %s", migrateErr)
			return
		}
	}

	// SMTPMailer
	smtpMailer, smtpMailerErr := smtpmailer.NewSimpleMailer(cfg.SMTPMailer)
	if smtpMailerErr != nil {
		zlogger.Error().Err(smtpMailerErr).Msgf("rest: main failed to construct smtp mailer client: %s", smtpMailerErr)
		return
	}

	// FileSystem
	fileSystemClient, fileSystemClientErr := filesystem.NewSimpleFSIO()
	if fileSystemClientErr != nil {
		zlogger.Error().Err(fileSystemClientErr).Msgf("rest: main failed to construct file system client: %s", fileSystemClientErr)
		return
	}

	// JSON Client
	jsonClient, jsonClientErr := jsonmanipulator.NewSimpleJSONManipulator()
	if jsonClientErr != nil {
		zlogger.Error().Err(jsonClientErr).Msgf("rest: main failed to construct json client: %s", jsonClientErr)
		return
	}

	// Message Queue (NSQ)
	nsqConfig := nsq.NewConfig()
	nsqdAddress := fmt.Sprintf("%s:%d", cfg.MessageQueue.Host, cfg.MessageQueue.Port)
	nsqProducer, err := nsq.NewProducer(nsqdAddress, nsqConfig)
	if err != nil {
		zlogger.Error().Err(err).Msgf("rest: main failed to construct nsq producer: %s", err)
	}

	messageQueue, messageQueueErr := messagequeue.NewMessageQueueNSQ(nsqProducer, nil)
	if messageQueueErr != nil {
		zlogger.Error().Err(messageQueueErr).Msgf("rest: main failed to construct nsq client: %s", messageQueueErr)
		return
	}

	// -----------------------------------------------------------------------------------------------------------------
	// SERVER SETUP AND EXECUTE
	// -----------------------------------------------------------------------------------------------------------------
	restServerHandler := rest.New(cfg, zlogger, sqlDB, smtpMailer, fileSystemClient, jsonClient, messageQueue)

	zlogger.Info().Msgf("REST Server started on port %d", cfg.API.RESTPort)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.API.RESTPort), restServerHandler)
}
