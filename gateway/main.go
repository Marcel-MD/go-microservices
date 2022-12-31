package main

import (
	"context"
	"gateway/http"
	"gateway/services"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {

	srv := http.GetServer()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	<-quit
	log.Warn().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	us := services.GetUserService()
	us.Close()

	<-ctx.Done()

	log.Info().Msg("Server exiting")
}
