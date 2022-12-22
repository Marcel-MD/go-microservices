package infrastructure

import (
	"user/domain"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg domain.Config) *gorm.DB {
	log.Info().Msg("Initializing database")

	dsn := cfg.DatabaseUrl

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	db.AutoMigrate(&domain.User{})

	return db
}
