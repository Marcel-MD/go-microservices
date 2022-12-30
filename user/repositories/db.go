package repositories

import (
	"sync"
	"user/config"
	"user/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbOnce sync.Once
var database *gorm.DB

func GetDB() *gorm.DB {
	dbOnce.Do(func() {

		log.Info().Msg("Initializing database")

		cfg := config.GetConfig()

		db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}

		db.AutoMigrate(&models.User{})

		database = db
	})

	return database
}

func Paginate(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
