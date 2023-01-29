package main

import (
	"file/repositories"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	<-quit
	log.Warn().Msg("Shutting down server...")

	if err := repositories.CloseDB(); err != nil {
		log.Fatal().Err(err).Msg("Failed to close db connection")
	}

	log.Info().Msg("Server exiting")
}
