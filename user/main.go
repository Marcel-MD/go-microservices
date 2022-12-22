package main

import (
	"net"
	"user/application"
	"user/domain"
	"user/infrastructure"
	"user/infrastructure/repositories"
	"user/pb"
	"user/presentation"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
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

	listener, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen.")
	}

	db := infrastructure.NewDB(cfg)
	repo := repositories.NewUserRepository(db)
	svc := application.NewUserService(repo, cfg)
	server := presentation.NewServer(svc)

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, server)

	log.Info().Msg("Starting server")
	if err := s.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve.")
	}
}
