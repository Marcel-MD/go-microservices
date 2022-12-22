package main

import (
	"gateway/domain"
	"gateway/infrastructure/services"
	"gateway/presentation"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load .env file.")
	}

	cfg := domain.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse environment variables.")
	}

	service := services.NewUserService(cfg)

	s := presentation.NewServer(cfg, service)

	if err := s.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run server.")
	}
}
