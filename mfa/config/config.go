package config

import (
	"sync"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port string `env:"PORT" envDefault:":8083"`

	OtpExpiry time.Duration `env:"OTP_EXPIRY" envDefault:"10m"`

	RabbitMQUrl string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`
	RedisUrl    string `env:"REDIS_URL" envDefault:"redis://:password@redis:6379/0"`
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
