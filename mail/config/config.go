package config

import (
	"sync"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port string `env:"PORT" envDefault:":8082"`

	KafkaUrl       string `env:"KAFKA_URL" envDefault:"kafka:9092"`
	KafkaTopic     string `env:"KAFKA_TOPIC" envDefault:"mails"`
	KafkaPartition int    `env:"KAFKA_PARTITION" envDefault:"0"`

	SmtpHost     string `env:"SMTP_HOST" envDefault:"smtp.gmail.com"`
	SmtpPort     string `env:"SMTP_PORT" envDefault:"587"`
	SmtpEmail    string `env:"SMTP_EMAIL" envDefault:""`
	SmtpPassword string `env:"SMTP_PASSWORD" envDefault:""`
	SenderName   string `env:"SENDER_NAME" envDefault:""`
}

var once sync.Once
var cfg Config

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
