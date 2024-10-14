package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	apiMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/bernardosecades/zaslink/internal/api/handler/health"
	"github.com/bernardosecades/zaslink/internal/api/handler/secret"
	"github.com/bernardosecades/zaslink/internal/api/middleware"
	"github.com/bernardosecades/zaslink/internal/config"
	"github.com/bernardosecades/zaslink/internal/events"
	"github.com/bernardosecades/zaslink/internal/repository"
	"github.com/bernardosecades/zaslink/internal/service"
	"github.com/bernardosecades/zaslink/pkg/observability"
	observabilityMiddleware "github.com/bernardosecades/zaslink/pkg/observability/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const timeOutHandlers = 30 * time.Second

func main() {
	ctx := context.Background()

	// LOGGER
	loggerOutput := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(loggerOutput).With().Timestamp().Logger()

	// CONFIG
	builder := config.Builder{}

	builder.Port(os.Getenv("PORT"))
	builder.SecretKey(os.Getenv("SECRET_KEY"))
	builder.DefaultPassword(os.Getenv("DEFAULT_PASSWORD"))
	builder.MongoDBURI(os.Getenv("MONGODB_URI"))
	builder.MongoDBName(os.Getenv("MONGODB_NAME"))
	builder.NatsURL(os.Getenv("NATS_URL"))

	cfg, err := builder.Build()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize configuration")
	}

	// OBSERVABILITY (OPEN TELEMETRY)

	/* TRACES */
	// TODO change sporter stdout to jaeger (zipkin is other option)
	consoleTraceExporter, err := observability.NewTraceExporter()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize observability trace exporter")
	}

	tracerProvider := observability.NewTraceProvider(consoleTraceExporter)
	defer func() { _ = tracerProvider.Shutdown(ctx) }()
	otel.SetTracerProvider(tracerProvider)

	/* METRICS */
	prometheusMetricExporter, err := prometheus.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize prometheus exporter")
	}

	// Create the resource to be observed
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("share service secret"),
			semconv.ServiceVersion("v0.0.1"),
		),
	)
	if err != nil {
		panic(err)
	}

	meterProvider := apiMetric.NewMeterProvider(
		apiMetric.WithResource(res),
		apiMetric.WithReader(prometheusMetricExporter),
	)
	defer func() { _ = meterProvider.Shutdown(ctx) }()
	otel.SetMeterProvider(meterProvider)

	meter := otel.Meter(
		"share secret api",
		metric.WithInstrumentationVersion("v0.0.1"),
	)

	prop := observability.NewPropagator()
	otel.SetTextMapPropagator(prop)

	// HANDLERS
	opts := options.Client().ApplyURI(cfg.MongoDBURI).SetConnectTimeout(5 * time.Second)
	client, _ := mongo.Connect(opts)

	secretRepo := repository.NewMongoDbSecretRepository(ctx, client, cfg.MongoDBName)
	secretService := service.NewSecretService(secretRepo, events.NewNatsPublisher(cfg.NatsURL), cfg.DefaultPassword, cfg.SecretKey)

	secretHandler := secret.NewHandler(secretService)
	healthHandler := health.NewHandler(cfg.MongoDBURI)

	// TODO move the router to internal package
	// ROUTER
	router := mux.NewRouter()

	v1 := router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/secret/{id}", secretHandler.RetrieveSecret).Methods(http.MethodGet)
	v1.HandleFunc("/secret/{private_id}", secretHandler.DeleteSecret).Methods(http.MethodDelete)
	v1.HandleFunc("/secret", secretHandler.CreateSecret).Methods(http.MethodPost)
	v1.HandleFunc("/ping", func(writer http.ResponseWriter, _ *http.Request) {
		_, err = writer.Write([]byte("pong"))
		if err != nil {
			fmt.Println(err)
		}
	}).Methods(http.MethodGet)

	router.HandleFunc("/healthz", healthHandler.Healthz).Methods(http.MethodGet)
	// Used to expose metrics in prometheus format
	router.Handle("/metrics", promhttp.Handler())

	// MIDDLEWARE
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger(logger))
	router.Use(observabilityMiddleware.NewMetricMiddleware(meter))

	// TODO Instrument the HTTP Server
	// SERVER
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// SHUTDOWN
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Printf("Received signal: %v. Initiating graceful shutdown...", sigChan)

		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal().Msg(fmt.Sprintf("HTTP shutdown error: %v", err))
		}
	}()

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://docs.zaslink.com",
			"https://www.zaslink.com",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	routerWithCors := c.Handler(router)

	// TIMEOUT endpoints handlers
	routerWithMiddlewares := http.TimeoutHandler(routerWithCors, timeOutHandlers, "Timeout!")
	http.Handle("/", routerWithMiddlewares)

	logger.Info().
		Str("PORT", server.Addr).
		Msg("HTTP server listening on port")

	// RUN SERVER
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal().
			Err(err).
			Msg("failed to start server")
	}
}
