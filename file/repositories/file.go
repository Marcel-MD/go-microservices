package repositories

import (
	"context"
	"file/models"
	"sync"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IFileRepository interface {
	FindByName(ctx context.Context, name string) (models.File, error)
	FindAll(ctx context.Context) ([]models.File, error)
	FindByOwnerId(ctx context.Context, ownerId string) ([]models.File, error)
	Create(ctx context.Context, file *models.File) error
	Delete(ctx context.Context, file *models.File) error
}

type fileRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

var fileOnce sync.Once
var fileRepo IFileRepository

func GetFileRepository() IFileRepository {
	fileOnce.Do(func() {
		log.Info().Msg("Initializing file repository")

		repo := fileRepository{
			db: GetDB(),
		}

		repo.coll = repo.db.Collection("files")

		fileRepo = &repo
	})

	return fileRepo
}

func (r *fileRepository) FindByName(ctx context.Context, name string) (models.File, error) {
	var file models.File
	err := r.coll.FindOne(ctx, bson.M{"name": name}).Decode(&file)
	if err != nil {
		return file, err
	}

	return file, nil
}

func (r *fileRepository) FindAll(ctx context.Context) ([]models.File, error) {
	var files []models.File
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return files, err
	}

	err = cursor.All(ctx, &files)
	if err != nil {
		return files, err
	}

	return files, nil
}

func (r *fileRepository) FindByOwnerId(ctx context.Context, ownerId string) ([]models.File, error) {
	var files []models.File
	cursor, err := r.coll.Find(ctx, bson.M{"owner_id": ownerId})
	if err != nil {
		return files, err
	}

	err = cursor.All(ctx, &files)
	if err != nil {
		return files, err
	}

	return files, nil
}

func (r *fileRepository) Create(ctx context.Context, file *models.File) error {
	_, err := r.coll.InsertOne(ctx, file)
	if err != nil {
		return err
	}

	return nil
}

func (r *fileRepository) Delete(ctx context.Context, file *models.File) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"name": file.Name})
	if err != nil {
		return err
	}

	return nil
}
