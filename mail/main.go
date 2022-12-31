package main

import (
	"context"
	"mail/kafka"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	reader := kafka.GetReader()
	ctx, cancel := context.WithCancel(context.Background())

	go reader.ReadMessages(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	<-quit
	log.Warn().Msg("Shutting down server...")

	cancel()

	log.Info().Msg("Server exiting")
}
