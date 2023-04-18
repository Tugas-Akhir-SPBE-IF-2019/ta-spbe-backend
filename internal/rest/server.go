package rest

import (
	"database/sql"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	assessmenthandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/assessment"
	authhandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/auth"
	indicatorassessmenthandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/indicator_assessment"
	institutionhandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/institution"
	userhandler "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/user"

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
	institutionStore := storepgsql.NewInstitution(sqlDB)

	indicatorAssessmentHandler := indicatorassessmenthandler.NewIndicatorAssessmentHandler(sqlDB, assessmentStore, indicatorAssessmentStore, userStore, smtpMailer, whatsAppClient, cfg.API)
	authHandler := authhandler.NewAuthHandler(sqlDB, userStore, cfg.OAuth, jwt)
	assessmentHandler := assessmenthandler.NewAssessmentHandler(sqlDB, assessmentStore, indicatorAssessmentStore, cfg.API, userStore, smtpMailer, fileSystemClient, jsonClient, messageQueue, whatsAppClient)
	institutionHandler := institutionhandler.NewInstitutionHandler(institutionStore)
	userHandler := userhandler.NewUserHandler(cfg.API, sqlDB, userStore, institutionStore, fileSystemClient, smtpMailer)

	r.Route("/auth", func(r chi.Router) {
		r.Get("/", token.HandleMain)
		r.Post("/google", authHandler.Google)
		r.Get("/google/callback", authHandler.GoogleCallback)
		r.Post("/google/validate", authHandler.GoogleValidate)
	})

	r.Get("/metrics", promHandler.ServeHTTP)
	r.Get("/assessments/index", indicatorAssessmentHandler.GetIndicatorAssessmentIndexList)
	r.Get("/institutions", institutionHandler.GetInstitutionList)
	r.Route("/assessments", func(r chi.Router) {
		r.Use(middleware.JWTAuth(jwt, cfg.DevSettings))
		r.Get("/", assessmentHandler.GetSPBEAssessmentList)
		r.Get("/{id}", indicatorAssessmentHandler.GetIndicatorAssessmentResultGetIndicatorAssessmentIndexList)
		r.Get("/{id}/histories", assessmentHandler.GetSPBEAssessmentStatusHistory)
		r.Get("/{id}/documents", assessmentHandler.GetSPBEAssessmentDocumentList)
		r.Get("/{id}/download", assessmentHandler.DownloadSupportDocuments)
		r.Patch("/{id}/validate", indicatorAssessmentHandler.ValidateIndicatorAssessmentResult)
		r.Post("/documents/upload", assessmentHandler.UploadSPBEDocument)

	})

	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.JWTAuth(jwt, cfg.DevSettings))
		r.Get("/profile", userHandler.GetUserProfile)
		r.Get("/evaluation", userHandler.GetUserEvaluationData)
		r.Get("/job", userHandler.GetUserJobData)
		r.Get("/institution", userHandler.GetUserCurrentInstitutionData)

		r.Get("/institution/{id}/approve", userHandler.VerifyUserCurrentInstitutionData)
		r.Get("/institution/{id}/reject", userHandler.RejectUserCurrentInstitutionData)

		r.Post("/evaluation", userHandler.AddUserEvaluationData)
		r.Post("/job", userHandler.AddUserJobData)
		r.Post("/institution", userHandler.AddUserInstitutionData)
		r.Put("/profile", userHandler.UpdateUserProfile)

		r.Delete("/institution/{id}", userHandler.DeleteUserCurrentInstitutionData)
	})

	r.Post("/assessments/result/callback", assessmentHandler.ResultCallback)

	// STATIC FILE SERVE (FOR DEVELOPMENT PURPOSE ONLY)
	fs := http.FileServer(http.Dir("static/supporting-documents"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	return r
}
