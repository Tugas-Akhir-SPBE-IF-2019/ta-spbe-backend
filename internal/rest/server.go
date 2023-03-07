package rest

import (
	"database/sql"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	authhandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/auth"
	indicatorassessmenthandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/indicator_assessment"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/middleware"
	storepgsql "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store/pgsql"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/metric"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/smtpmailer"
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
	smtpMailer smtpmailer.Client,
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
	userStore := storepgsql.NewUser(sqlDB)

	indicatorAssessmentHandler := indicatorassessmenthandler.NewIndicatorAssessmentHandler(sqlDB, indicatorAssessmentStore, userStore, smtpMailer)
	authHandler := authhandler.NewAuthHandler(sqlDB, userStore, cfg.OAuth, jwt)
	r.Route("/auth", func(r chi.Router) {
		r.Get("/", token.HandleMain)
		r.Post("/google", authHandler.Google)
		r.Get("/google/callback", authHandler.GoogleCallback)
	})

	r.Get("/metrics", promHandler.ServeHTTP)
	r.Get("/assessments/index", indicatorAssessmentHandler.GetIndicatorAssessmentIndexList)
	r.Route("/assessments", func(r chi.Router) {
		r.Use(middleware.JWTAuth(jwt))
		r.Get("/{id}", indicatorAssessmentHandler.GetIndicatorAssessmentResultGetIndicatorAssessmentIndexList)
		r.Patch("/{id}/validate", indicatorAssessmentHandler.ValidateIndicatorAssessmentResult)

	})

	r.Post("/assessments/result/callback", indicatorAssessmentHandler.ResultCallback)

	return r
}
