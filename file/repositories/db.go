package repositories

import (
	"context"
	"file/config"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client   *mongo.Client
	database *mongo.Database
	ctx      context.Context
	dbOnce   sync.Once
)

func GetDB() *mongo.Database {
	dbOnce.Do(func() {

		log.Info().Msg("Initializing database")

		cfg := config.GetConfig()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}

		client, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.DatabaseUrl))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}

		database = client.Database("file")
	})

	return database
}

func CloseDB() error {
	return client.Disconnect(ctx)
}
