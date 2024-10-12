package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/bernardosecades/zaslink/pkg/notifier/telegram"
	"github.com/nats-io/nats.go"
)

func main() {
	// LOGGER
	loggerOutput := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(loggerOutput).With().Timestamp().Logger()

	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		logger.Fatal().Msg("Environment variable NATS_URL is not defined")
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	botAPIToken, ok := os.LookupEnv("NOTIFIER_TELEGRAM_BOT_TOKEN")
	if !ok {
		logger.Fatal().Msg("Environment variable NOTIFIER_TELEGRAM_BOT_TOKEN is not defined")
	}
	telegramService, ok := os.LookupEnv("NOTIFIER_TELEGRAM_USER_ID")
	if !ok {
		logger.Fatal().Msg("Environment variable NOTIFIER_TELEGRAM_USER_ID is not defined")
	}

	notifier := telegram.NewNotifier(botAPIToken, telegramService)

	// This will capture messages for: notifications.telegram.*, notifications.email.*, ...
	subject := "notifications.>"

	// Subscribe to the subject
	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		// TODO Add Backoff Logic to make the system more resilient.
		err = notifier.Notify(msg.Subject, string(msg.Data))
		if err != nil {
			logger.Error().Err(err).Msg("Failed to notify Telegram")
			return
		}

		logger.Info().
			Str("subject", msg.Subject).
			Str("data", string(msg.Data)).
			Msg("Received a message")
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to subscribe to NATS subject")
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Keep the program alive until interrupted
	log.Println("Waiting for messages. Press Ctrl+C to exit.")
	<-sigChan

	logger.Info().Msg("Shutting down...")

	// Gracefully close the NATS connection
	if err := nc.Drain(); err != nil {
		logger.Error().Err(err).Msg("Error draining NATS connection")
	}
	logger.Info().Msg("NATS connection closed.")
}
