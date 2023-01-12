package config

import (
	"sync"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port          string        `env:"PORT" envDefault:":8081"`
	ApiSecret     string        `env:"API_SECRET" envDefault:"SecretSecretSecret"`
	TokenLifespan time.Duration `env:"TOKEN_LIFESPAN" envDefault:"24h"`

	RabbitMQUrl string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`
	DatabaseUrl string `env:"DATABASE_URL" envDefault:"postgres://postgres:password@postgres:5432/users"`
}

var (
	once sync.Once
	cfg  Config
)

func GetConfig() Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to load .env file.")
		}

		if err := env.Parse(&cfg); err != nil {
			log.Fatal().Err(err).Msg("Failed to parse environment variables.")
		}
	})

	return cfg
}
