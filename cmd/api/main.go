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

	"github.com/bernardosecades/sharesecret/internal/api/handler/health"
	"github.com/bernardosecades/sharesecret/internal/api/handler/secret"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bernardosecades/sharesecret/internal/api/middleware"
	"github.com/bernardosecades/sharesecret/internal/repository"
	"github.com/bernardosecades/sharesecret/internal/service"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	if err := prometheus.Register(middleware.TotalRequests); err != nil {
		panic(err)
	}
}

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

	secretRepo := repository.NewMongoDbSecretRepository(ctx, mongoDBUri, mongoDBName)
	secretService := service.NewSecretService(secretRepo, defaultPassword, secretKey)

	// HANDLERS
	secretHandler := secret.NewHandler(secretService)
	healthHandler := health.NewHandler(mongoDBUri)

	// ROUTER
	router := mux.NewRouter()
	router.HandleFunc("/secret/{id}", secretHandler.RetrieveSecret).Methods(http.MethodGet)
	router.HandleFunc("/secret", secretHandler.CreateSecret).Methods(http.MethodPost)

	router.HandleFunc("/healthz", healthHandler.Healthz).Methods(http.MethodGet)
	router.Path("/prometheus").Handler(promhttp.Handler())

	// LOGGER
	loggerOutput := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(loggerOutput)

	// MIDDLEWARE
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Prometheus)

	// SERVER
	server := &http.Server{
		Addr:              ":8080",
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

	http.Handle("/", router)
	log.Info().Msg(fmt.Sprintf("HTTP server listening on port %s", server.Addr))

	// RUN SERVER
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Msg(fmt.Sprintf("HTTP server error: %v", err))
	}
}
