package config

import (
	"sync"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port          string        `env:"PORT" envDefault:":8084"`
	ApiSecret     string        `env:"API_SECRET" envDefault:"SecretSecretSecret"`
	TokenLifespan time.Duration `env:"TOKEN_LIFESPAN" envDefault:"24h"`

	AllowOrigin string `env:"ALLOW_ORIGIN" envDefault:"*"`

	DatabaseUrl string `env:"DATABASE_URL" envDefault:"mongodb://root:password@mongo:27017"`

	AzureBlobConnectionString string `env:"AZURE_BLOB_CONNECTION_STRING" envDefault:"AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;DefaultEndpointsProtocol=http;BlobEndpoint=http://azurite:10000/devstoreaccount1;QueueEndpoint=http://azurite:10001/devstoreaccount1;TableEndpoint=http://azurite:10002/devstoreaccount1;"`
	AzureBlobContainerName    string `env:"AZURE_BLOB_CONTAINER_NAME" envDefault:"files"`
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
