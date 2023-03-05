package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/logger"
	"github.com/jackc/pgx/v5"
	"strings"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Config struct {
	LogLevel    string `toml:"log_level"`
	Username    string `toml:"username"`
	Password    string `toml:"password"`
	Database    string `toml:"database"`
	Host        string `toml:"host"`
	SSLMode     string `toml:"ssl_mode"`
	Port        int    `toml:"port"`
	ConnMaxOpen int    `toml:"conn_max_open"`
	ConnMaxIdle int    `toml:"conn_max_idle"`
	Logging     bool   `toml:"logging"`
	Tracing     bool   `toml:"tracing"`
	Migration   bool   `toml:"migration"`
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Database, c.Username, c.Password, c.SSLMode,
	)
}

// NewDB create new standard library *sql.DB object (with configurable query logger).
func NewDB(cfg Config, zlogger zerolog.Logger) (*sql.DB, error) {
	conCfg, err := pgx.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("pgsql: parse DSN failed: %w", err)
	}
	// default log to none
	traceLogger := &tracelog.TraceLog{
		LogLevel: tracelog.LogLevelNone,
	}
	moduleName := "pgsql"
	if cfg.Logging && zlogger.GetLevel() != zerolog.Disabled {
		logLevel, logLevelErr := tracelog.LogLevelFromString(strings.TrimSpace(strings.ToLower(cfg.LogLevel)))
		if logLevelErr != nil {
			return nil, fmt.Errorf("pgsql: parse log level '%s' failed: %w", cfg.LogLevel, logLevelErr)
		}
		traceLogger.Logger = &lg{
			moduleName: moduleName,
			logger:     zlogger.With().Str("module", moduleName).Logger(),
		}
		traceLogger.LogLevel = logLevel
	}
	conCfg.Tracer = traceLogger
	db := stdlib.OpenDB(*conCfg)
	if cfg.Tracing {
		db = otelsql.OpenDB(stdlib.GetConnector(*conCfg),
			otelsql.WithDBName(cfg.Database),
			otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		)
	}
	db.SetMaxOpenConns(cfg.ConnMaxOpen)
	db.SetMaxIdleConns(cfg.ConnMaxIdle)
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("pgsql: DB ping failed: %w", err)
	}
	return db, nil
}

type lg struct {
	moduleName string
	logger     zerolog.Logger
}

func (pl *lg) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	var lvl zerolog.Level
	switch level {
	case tracelog.LogLevelNone:
		lvl = zerolog.NoLevel
	case tracelog.LogLevelError:
		lvl = zerolog.ErrorLevel
	case tracelog.LogLevelWarn:
		lvl = zerolog.WarnLevel
	case tracelog.LogLevelInfo:
		lvl = zerolog.InfoLevel
	case tracelog.LogLevelDebug:
		lvl = zerolog.DebugLevel
	case tracelog.LogLevelTrace:
		lvl = zerolog.TraceLevel
	default:
		lvl = zerolog.DebugLevel
	}
	// prioritize logger carried from given context,
	// because it is ideally prepended with trace or request id per request-scoped context.
	var event *zerolog.Event
	if ctxLg := logger.FromContext(ctx); ctxLg.GetLevel() != zerolog.Disabled {
		event = ctxLg.WithLevel(lvl).Str("module", pl.moduleName)
	} else {
		event = pl.logger.WithLevel(lvl)
	}
	if !event.Enabled() {
		return
	}
	// time field duplicated with common timestamp field, in pgx it's actual duration of the query execution.
	if duration, ok := data["time"]; ok {
		data["duration"] = duration
		delete(data, "time")
	}
	delete(data, "pid")
	delete(data, "commandTag")
	if query, ok := data["sql"]; ok {
		// use SQL query as the log message
		if sqlQuery, isString := query.(string); isString {
			msg = sqlQuery
		}
	}
	event.Fields(data).Msg(msg)
}
