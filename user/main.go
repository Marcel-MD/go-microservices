package main

import (
	"os"
	"os/signal"
	"syscall"
	"user/repositories"
	"user/rpc"
	"user/services"

	"github.com/rs/zerolog/log"
)

func main() {

	srv, listener := rpc.GetServer()

	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	<-quit
	log.Warn().Msg("Shutting down server...")

	srv.GracefulStop()

	if err := repositories.CloseDB(); err != nil {
		log.Fatal().Err(err).Msg("Failed to close db connection")
	}

	ms := services.GetMailService()
	if err := ms.Close(); err != nil {
		log.Fatal().Err(err).Msg("Failed to close mail service")
	}

	log.Info().Msg("Server exiting")
}
