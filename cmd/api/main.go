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

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	apiMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/bernardosecades/sharesecret/internal/api/handler/health"
	"github.com/bernardosecades/sharesecret/internal/api/handler/secret"
	"github.com/bernardosecades/sharesecret/internal/api/middleware"
	"github.com/bernardosecades/sharesecret/internal/repository"
	"github.com/bernardosecades/sharesecret/internal/service"
	"github.com/bernardosecades/sharesecret/pkg/observability"
	observabilityMiddleware "github.com/bernardosecades/sharesecret/pkg/observability/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const timeOutHandlers = 30 * time.Second

func main() {
	ctx := context.Background()

	// ENVIRONMENT VARIABLES
	secretKey, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		panic("SECRET_KEY is not present")
	}
	defaultPassword, ok := os.LookupEnv("DEFAULT_PASSWORD")
	if !ok {
		panic("DEFAULT_PASSWORD is not present")
	}
	mongoDBUri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		panic("DEFAULT_PASSWORD is not present")
	}
	mongoDBName, ok := os.LookupEnv("MONGODB_NAME")
	if !ok {
		panic("MONGODB_NAME is not present")
	}

	// LOGGER
	loggerOutput := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(loggerOutput)

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
	opts := options.Client().ApplyURI(mongoDBUri).SetConnectTimeout(5 * time.Second)
	client, _ := mongo.Connect(opts)

	secretRepo := repository.NewMongoDbSecretRepository(ctx, client, mongoDBName)
	secretService := service.NewSecretService(secretRepo, defaultPassword, secretKey)

	secretHandler := secret.NewHandler(secretService)
	healthHandler := health.NewHandler(mongoDBUri)

	// ROUTER
	router := mux.NewRouter()
	router.HandleFunc("/secret/{id}", secretHandler.RetrieveSecret).Methods(http.MethodGet)
	router.HandleFunc("/secret", secretHandler.CreateSecret).Methods(http.MethodPost)

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
		Addr:              ":8080",
		ReadHeaderTimeout: 60 * time.Second,
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

	routerWithMiddlewares := http.TimeoutHandler(router, timeOutHandlers, "Timeout!")

	http.Handle("/", routerWithMiddlewares)
	log.Info().Msg(fmt.Sprintf("HTTP server listening on port %s", server.Addr))

	// RUN SERVER
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Msg(fmt.Sprintf("HTTP server error: %v", err))
	}
}
