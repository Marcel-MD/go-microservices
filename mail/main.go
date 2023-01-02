package main

import (
	"mail/rabbitmq"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	consumer := rabbitmq.GetConsumer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	<-quit
	log.Warn().Msg("Shutting down server...")

	err := consumer.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to close consumer")
	}

	log.Info().Msg("Server exiting")
}
