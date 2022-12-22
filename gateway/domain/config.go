package domain

import "time"

type Config struct {
	Port          string        `env:"PORT" envDefault:":8080"`
	ApiSecret     string        `env:"API_SECRET" envDefault:"SecretSecretSecret"`
	TokenLifespan time.Duration `env:"TOKEN_LIFESPAN" envDefault:"24h"`

	UserServiceUrl string `env:"USER_SERVICE_URL" envDefault:"localhost:8081"`
}
