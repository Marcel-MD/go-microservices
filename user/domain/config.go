package domain

import "time"

type Config struct {
	Port          string        `env:"PORT" envDefault:":8081"`
	ApiSecret     string        `env:"API_SECRET" envDefault:"SecretSecretSecret"`
	TokenLifespan time.Duration `env:"TOKEN_LIFESPAN" envDefault:"24h"`

	DatabaseUrl string `env:"DATABASE_URL" envDefault:"postgres://postgres:password@localhost:5432/users"`
}
