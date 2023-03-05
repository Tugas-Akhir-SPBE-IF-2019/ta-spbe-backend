package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	assessmenthandler "ta-spbe-backend/api/handler/assessment"
	authhandler "ta-spbe-backend/api/handler/auth"
	indicatorassessmenthandler "ta-spbe-backend/api/handler/indicator_assessment"
	apimiddleware "ta-spbe-backend/api/middleware"
	filesystem "ta-spbe-backend/service/file-system"
	jsonmanipulator "ta-spbe-backend/service/json-manipulator"
	mailerservice "ta-spbe-backend/service/mailer"
	messagequeue "ta-spbe-backend/service/message-queue"
	"ta-spbe-backend/service/whatsapp"
	"time"

	"ta-spbe-backend/config"
	"ta-spbe-backend/database"
	prometheusinst "ta-spbe-backend/instrumentation/prometheus"
	"ta-spbe-backend/repository/pgsql"
	"ta-spbe-backend/token"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nsqio/go-nsq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var whatsappClient = &whatsapp.WhatsMeow{}

// Jaeger
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

// Prometheus
type Device struct {
	ID       int    `json:"id"`
	Mac      string `json:"mac"`
	Firmware string `json:"firmware"`
}

var dvs []Device

func init() {
	dvs = []Device{
		{1, "5F-33-CC-1F-43-82", "2.1.6"},
		{2, "EF-2B-C4-F5-D6-34", "2.1.6"},
	}
}

func main() {
	fmt.Println("test")
	log.Println("test")
	cfgPath := flag.String("c", "config.toml", "path to config file")
	cfg, err := config.LoadEnvFromFile(*cfgPath)
	if err != nil {
		log.Fatalln(err)
	}

	migrate := flag.Bool("migrate", cfg.DB.Migration, "do migration")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		url.QueryEscape(cfg.DB.Username),
		url.QueryEscape(cfg.DB.Password),
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	)

	db, err := sql.Open("postgres", dsn)

	log.Println("wow")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("wows")
	log.Println(db)
	err = db.Ping()
	if err != nil {
		log.Println(err.Error())
		log.Fatal(err)
	}

	log.Println("wowd")

	if *migrate {
		err := database.Migrate(db, cfg.DB.Database)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println("wowa")
	assessmentRepo, err := pgsql.NewAssessmentRepo(db)
	if err != nil {
		log.Fatalln(err)
	}

	indicatorAssessmentRepo, err := pgsql.NewIndicatorAssessmentRepo(db)
	if err != nil {
		log.Fatalln(err)
	}

	userRepo, err := pgsql.NewUserRepo(db)
	if err != nil {
		log.Fatalln(err)
	}

	nsqConfig := nsq.NewConfig()
	nsqdAddress := fmt.Sprintf("%s:%d", cfg.MessageBroker.Host, cfg.MessageBroker.Port)
	nsqProducer, err := nsq.NewProducer(nsqdAddress, nsqConfig)
	if err != nil {
		log.Fatalln(err)
	}

	messageQueue := &messagequeue.NSQ{
		Producer: nsqProducer,
	}
	mailerService, err := mailerservice.NewSimpleMailer(cfg.SMTPClient)
	if err != nil {
		log.Fatalln(err)
	}
	jsonEC := jsonmanipulator.EncoderDecoder{}

	fileSystemIO := filesystem.IO{}

	//WhatsMeow Init
	go func() {
		dbLog := waLog.Stdout("Database", "DEBUG", true)
		// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
		container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
		if err != nil {
			log.Println(err.Error())
		}
		// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
		deviceStore, err := container.GetFirstDevice()
		if err != nil {
			log.Println(err.Error())
		}
		clientLog := waLog.Stdout("Client", "DEBUG", true)
		client := whatsmeow.NewClient(deviceStore, clientLog)

		if client.Store.ID == nil {
			// No ID stored, new login
			qrChan, _ := client.GetQRChannel(context.Background())
			err = client.Connect()
			if err != nil {
				panic(err)
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					// Render the QR code here
					// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
					fmt.Println("QR code:", evt.Code)
				} else {
					fmt.Println("Login event:", evt.Event)
				}
			}
		} else {
			// Already logged in, just connect
			err = client.Connect()
			if err != nil {
				log.Printf(err.Error())
			}
		}

		whatsappClient = &whatsapp.WhatsMeow{
			Client: client,
		}
		//WhatsMeow Init
	}()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))

	//Prometheus test
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	fmt.Println(dvs)
	m.Devices.Set(float64(len(dvs)))

	mdh := manageDevicesHandler{metrics: m}
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	r.Get("/metrics", promHandler.ServeHTTP)
	r.Put("/devices/{id}", mdh.ServeHTTP)

	//End of prometheus

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Tugas Akhir Otomatisasi Penilaian Tingkat Kematangan Kebijakan SPBE IF 2019"))
	})

	jwt := token.NewJWT(cfg.JWT)
	r.Route("/auth", func(r chi.Router) {
		r.Get("/", token.HandleMain)
		r.Post("/google", authhandler.Google(cfg.OAuth))
		r.Get("/google/callback", authhandler.GoogleCallback(userRepo, cfg.OAuth, jwt))
	})

	authMW := apimiddleware.Auth(jwt)

	//Jaeger test
	spbeIndexHandlerJaegerTest := indicatorassessmenthandler.GetIndicatorAssessmentIndexList(indicatorAssessmentRepo, m)
	r.Get("/assessments/index", otelhttp.NewHandler(http.HandlerFunc(spbeIndexHandlerJaegerTest), "assessments-index").ServeHTTP)
	// r.Get("/assessments/index",  indicatorassessmenthandler.GetIndicatorAssessmentIndexList(indicatorAssessmentRepo, m))
	r.Route("/assessments", func(r chi.Router) {
		r.Use(authMW)
		r.Get("/", assessmenthandler.GetSPBEAssessmentList(assessmentRepo, m))
		r.Get("/{id}", indicatorassessmenthandler.GetIndicatorAssessmentResult(indicatorAssessmentRepo))
		r.Post("/documents/upload", assessmenthandler.UploadSPBEDocument(assessmentRepo, userRepo, messageQueue, mailerService, fileSystemIO, jsonEC, whatsappClient, cfg.API, cfg.SMTPClient))
		r.Patch("/{id}/validate", indicatorassessmenthandler.ValidateIndicatorAssessmentResult(indicatorAssessmentRepo))
	})

	r.Post("/assessments/result/callback", indicatorassessmenthandler.ResultCallback(indicatorAssessmentRepo, userRepo, mailerService, whatsappClient))

	//static file serve (for testing purpose only)
	fs := http.FileServer(http.Dir("static/supporting-documents"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Jaeger

	// Tracer
	tp, err := tracerProvider("http://jaeger:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	tr := tp.Tracer("component-main")

	ctx, span := tr.Start(ctx, "hello")
	defer span.End()

	// HTTP Handlers
	helloHandler := func(w http.ResponseWriter, r *http.Request) {
		// Use the global TracerProvider
		tr := otel.Tracer("hello-handler")
		_, span := tr.Start(r.Context(), "hello")
		span.SetAttributes(attribute.Key("mykey").String("value"))
		defer span.End()

		yourName := "Hardcoded name"
		fmt.Fprintf(w, "Hello %q!", yourName)
	}

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello")
	// Jaeger

	r.Get("/otel", otelHandler.ServeHTTP)

	log.Printf("Server is listening on port %d", 8080)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), r); err != nil {
		log.Println(err)

	}
}

//Prometheus metrics test

func NewMetrics(reg prometheus.Registerer) *prometheusinst.Metrics {
	m := &prometheusinst.Metrics{
		Devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "connected_devices",
			Help:      "Number of currently connected devices.",
		}),
		Upgrades: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "myapp",
			Name:      "device_upgrade_total",
			Help:      "Number of requests",
		}, []string{"type"}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "myapp",
			Name:      "request_duration_seconds",
			Help:      "Duration of the request.",
			// 4 times larger for apdex score
			// Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
			// Buckets: prometheus.LinearBuckets(0.1, 5, 5),
			Buckets: []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"status", "method"}),
	}
	reg.MustRegister(m.Devices, m.Upgrades, m.Duration)
	return m
}

type manageDevicesHandler struct {
	metrics *prometheusinst.Metrics
}

func upgradeDevice(w http.ResponseWriter, r *http.Request, m *prometheusinst.Metrics) {
	path := chi.URLParam(r, "id")

	id, err := strconv.Atoi(path)
	if err != nil || id < 1 {
		http.NotFound(w, r)
	}

	var dv Device
	err = json.NewDecoder(r.Body).Decode(&dv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range dvs {
		if dvs[i].ID == id {
			dvs[i].Firmware = dv.Firmware
		}
	}

	m.Upgrades.With(prometheus.Labels{"type": "router"}).Inc()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Upgrading..."))
}

func (mdh manageDevicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgradeDevice(w, r, mdh.metrics)
}
