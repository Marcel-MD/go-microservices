package config

import (
	"sync"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port string `env:"PORT" envDefault:":8084"`

	DatabaseUrl string `env:"DATABASE_URL" envDefault:"mongodb://root:password@mongo:27017"`

	AzureEndpoint  string `env:"AZURE_ENDPOINT" envDefault:"http://localhost:10000"`
	AzureName      string `env:"AZURE_NAME" envDefault:"devstoreaccount1"`
	AzureKey       string `env:"AZURE_KEY" envDefault:"Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="`
	AzureContainer string `env:"AZURE_CONTAINER" envDefault:"files"`
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
