package services

import (
	"context"
	"errors"
	"file/models"
	"file/repositories"
	"io"
	"sync"

	"github.com/rs/zerolog/log"
)

type IFileService interface {
	FindAll(ctx context.Context) ([]models.File, error)
	FindByName(ctx context.Context, name string) (models.File, error)
	FindByOwnerId(ctx context.Context, ownerId string) ([]models.File, error)

	Upload(ctx context.Context, reader io.Reader, extension string, ownerId string) (models.File, error)
	Read(ctx context.Context, name string) (io.Reader, error)
	Delete(ctx context.Context, name, ownerId string) error
}

type fileService struct {
	fileRepository repositories.IFileRepository
	blobRepository repositories.IBlobRepository
}

var (
	fileOnce sync.Once
	fileRepo IFileService
)

func GetFileService() IFileService {
	fileOnce.Do(func() {
		log.Info().Msg("Initializing file service")

		fileRepo = &fileService{
			fileRepository: repositories.GetFileRepository(),
			blobRepository: repositories.GetBlobRepository(),
		}
	})

	return fileRepo
}

func (s *fileService) FindAll(ctx context.Context) ([]models.File, error) {
	log.Debug().Msg("Finding all files")

	return s.fileRepository.FindAll(ctx)
}

func (s *fileService) FindByName(ctx context.Context, name string) (models.File, error) {
	log.Debug().Str("name", name).Msg("Finding file")

	return s.fileRepository.FindByName(ctx, name)
}

func (s *fileService) FindByOwnerId(ctx context.Context, ownerId string) ([]models.File, error) {
	log.Debug().Str("ownerId", ownerId).Msg("Finding files by owner id")

	return s.fileRepository.FindByOwnerId(ctx, ownerId)
}

func (s *fileService) Upload(ctx context.Context, reader io.Reader, extension string, ownerId string) (models.File, error) {
	log.Debug().Str("ownerId", ownerId).Msg("Uploading file")

	var file models.File

	name, err := s.blobRepository.Upload(ctx, extension, reader)
	if err != nil {
		return file, err
	}

	file = models.File{
		Name:      name,
		OwnerId:   ownerId,
		Extension: extension,
	}

	err = s.fileRepository.Create(ctx, &file)
	if err != nil {
		return file, err
	}

	return file, nil
}

func (s *fileService) Read(ctx context.Context, name string) (io.Reader, error) {
	log.Debug().Str("name", name).Msg("Reading file")

	return s.blobRepository.Download(ctx, name)
}

func (s *fileService) Delete(ctx context.Context, name, ownerId string) error {
	log.Debug().Str("name", name).Msg("Deleting file")

	file, err := s.fileRepository.FindByName(ctx, name)
	if err != nil {
		return err
	}

	if file.OwnerId != ownerId {
		return errors.New("file does not belong to user")
	}

	err = s.fileRepository.Delete(ctx, &file)
	if err != nil {
		return err
	}

	return s.blobRepository.Delete(ctx, name)
}
