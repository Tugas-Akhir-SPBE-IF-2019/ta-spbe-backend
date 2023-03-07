package rest

import (
	"database/sql"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	assessmenthandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/assessment"
	authhandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/auth"
	indicatorassessmenthandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/indicator_assessment"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/middleware"
	storepgsql "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store/pgsql"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/filesystem"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/jsonmanipulator"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/messagequeue"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/metric"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/smtpmailer"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/token"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/whatsapp"
	"github.com/rs/zerolog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(
	cfg *config.Config,
	zlogger zerolog.Logger,
	sqlDB *sql.DB,
	smtpMailer smtpmailer.Client,
	fileSystemClient filesystem.Client,
	jsonClient jsonmanipulator.Client,
	messageQueue messagequeue.Client,
	whatsAppClient whatsapp.Client,

) http.Handler {
	r := chi.NewRouter()

	reg := prometheus.NewRegistry()
	m := metric.NewMetrics(reg)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	jwt := token.NewJWT(cfg.JWT)
	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
		}),
		middleware.HTTPTracer,
		middleware.RequestID(zlogger),
		middleware.HTTPLogger,
		middleware.HTTPMetric(m),
	)

	assessmentStore := storepgsql.NewAssessment(sqlDB)
	indicatorAssessmentStore := storepgsql.NewIndicatorAssessment(sqlDB)
	userStore := storepgsql.NewUser(sqlDB)

	indicatorAssessmentHandler := indicatorassessmenthandler.NewIndicatorAssessmentHandler(sqlDB, indicatorAssessmentStore, userStore, smtpMailer)
	authHandler := authhandler.NewAuthHandler(sqlDB, userStore, cfg.OAuth, jwt)
	assessmentHandler := assessmenthandler.NewAssessmentHandler(sqlDB, assessmentStore, cfg.API, userStore, smtpMailer, fileSystemClient, jsonClient, messageQueue, whatsAppClient)
	r.Route("/auth", func(r chi.Router) {
		r.Get("/", token.HandleMain)
		r.Post("/google", authHandler.Google)
		r.Get("/google/callback", authHandler.GoogleCallback)
	})

	r.Get("/metrics", promHandler.ServeHTTP)
	r.Get("/assessments/index", indicatorAssessmentHandler.GetIndicatorAssessmentIndexList)
	r.Route("/assessments", func(r chi.Router) {
		r.Use(middleware.JWTAuth(jwt, cfg.DevSettings))
		r.Get("/", assessmentHandler.GetSPBEAssessmentList)
		r.Get("/{id}", indicatorAssessmentHandler.GetIndicatorAssessmentResultGetIndicatorAssessmentIndexList)
		r.Patch("/{id}/validate", indicatorAssessmentHandler.ValidateIndicatorAssessmentResult)
		r.Post("/documents/upload", assessmentHandler.UploadSPBEDocument)

	})

	r.Post("/assessments/result/callback", indicatorAssessmentHandler.ResultCallback)

	return r
}
