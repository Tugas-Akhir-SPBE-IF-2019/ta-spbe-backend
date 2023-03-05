package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	indicatorassessmenthandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/indicator_assessment"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/middleware"
	storepgsql "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store/pgsql"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/logger"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/metric"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/pgsql"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/tracer"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type Message struct {
	Message string `json:"data"`
}

var sqlDB *sql.DB

func main() {
	cfgPath := flag.String("c", "config.toml", "path to config file")
	cfg, err := config.LoadEnvFromFile(*cfgPath)
	if err != nil {
		log.Fatalln(err)
	}

	zlogger := logger.New(cfg.Logger).With().
		Logger()

	var sqlDBErr error
	sqlDB, sqlDBErr = pgsql.NewDB(cfg.PostgreSQL, zlogger)
	if sqlDBErr != nil {
		zlogger.Error().Err(sqlDBErr).Msgf("grpc: main failed to construct pgsql: %s", sqlDBErr)
		return
	}

	// // Tracer
	// tp, err := tracerProvider("http://jaeger:14268/api/traces")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Register our TracerProvider as the global so any imported
	// // instrumentation in the future will default to using it.
	// otel.SetTracerProvider(tp)

	if err = tracer.SetTracer(cfg.Tracer, cfg.AppInfo); err != nil {
		zlogger.Error().Err(err).Msgf("grpc: main failed to setup open telemetry tracer: %s", err)
		return
	}

	helloHandler := func(w http.ResponseWriter, r *http.Request) {
		sqlDB.ExecContext(r.Context(), "SELECT * FROM world;")

		response, _ := json.Marshal(Message{
			Message: "hello world",
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}

	r := chi.NewRouter()

	reg := prometheus.NewRegistry()
	m := metric.NewMetrics(reg)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	r.Use(
		middleware.HTTPTracer,
		middleware.RequestID(zlogger),
		middleware.HTTPLogger,
		middleware.HTTPMetric(m),
	)

	indicatorAssessmentStore := storepgsql.NewIndicatorAssessment(sqlDB)
	indicatorAssessmentHandler := indicatorassessmenthandler.NewIndicatorAssessmentHandler(sqlDB, indicatorAssessmentStore)

	r.Get("/metrics", promHandler.ServeHTTP)

	r.Get("/",
		helloHandler)

	r.Get("/hellos",
		helloHandler)

	r.Get("/handler",
		helloHandler)

	r.Get("/index", indicatorAssessmentHandler.GetIndicatorAssessmentIndexList)

	http.ListenAndServe(":3001", r)
}

const (
	service     = "go-gopher-opentelemetry"
	environment = "development"
	id          = 1
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}
