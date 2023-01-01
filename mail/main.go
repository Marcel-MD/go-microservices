package main

import (
	"context"
	"mail/rabbitmq"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	receiver := rabbitmq.GetReceiver()
	go receiver.Listen(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	<-quit
	log.Warn().Msg("Shutting down server...")

	cancel()

	err := receiver.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to close receiver")
	}

	log.Info().Msg("Server exiting")
}
