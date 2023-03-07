package rest

import (
	"database/sql"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	indicatorassessmenthandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/indicator_assessment"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/middleware"
	storepgsql "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store/pgsql"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/metric"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/token"
	"github.com/rs/zerolog"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(
	cfg *config.Config,
	zlogger zerolog.Logger,
	sqlDB *sql.DB,
) http.Handler {
	r := chi.NewRouter()

	reg := prometheus.NewRegistry()
	m := metric.NewMetrics(reg)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	jwt := token.NewJWT(cfg.JWT)
	r.Use(
		middleware.HTTPTracer,
		middleware.RequestID(zlogger),
		middleware.HTTPLogger,
		middleware.HTTPMetric(m),
	)

	indicatorAssessmentStore := storepgsql.NewIndicatorAssessment(sqlDB)
	indicatorAssessmentHandler := indicatorassessmenthandler.NewIndicatorAssessmentHandler(sqlDB, indicatorAssessmentStore)

	r.Get("/metrics", promHandler.ServeHTTP)
	r.Get("/assessments/index", indicatorAssessmentHandler.GetIndicatorAssessmentIndexList)
	r.Route("/assessments", func(r chi.Router) {
		r.Use(middleware.JWTAuth(jwt))
		r.Get("/{id}", indicatorAssessmentHandler.GetIndicatorAssessmentResultGetIndicatorAssessmentIndexList)

	})

	return r
}
