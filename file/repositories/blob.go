package repositories

import (
	"context"
	"file/config"
	"io"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type IBlobRepository interface {
	Download(ctx context.Context, name string) (io.Reader, error)
	Upload(ctx context.Context, extension string, reader io.Reader) (string, error)
	Delete(ctx context.Context, name string) error
}

type blobRepository struct {
	client    *azblob.Client
	container string
}

var (
	blobOnce sync.Once
	blobRepo IBlobRepository
)

func GetBlobRepository() IBlobRepository {
	blobOnce.Do(func() {
		log.Info().Msg("Initializing blob repository")

		cfg := config.GetConfig()

		client, err := azblob.NewClientFromConnectionString(cfg.AzureBlobConnectionString, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize blob repository")
		}

		_, err = client.CreateContainer(context.Background(), cfg.AzureBlobContainerName, nil)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to create container, it may already exist")
		}

		blobRepo = &blobRepository{
			client:    client,
			container: cfg.AzureBlobContainerName,
		}
	})

	return blobRepo
}

func (b *blobRepository) Upload(ctx context.Context, extension string, reader io.Reader) (string, error) {
	log.Debug().Msg("Uploading blob")

	name := uuid.New().String() + extension

	_, err := b.client.UploadStream(ctx, b.container, name, reader, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload blob")
		return "", err
	}

	return name, nil
}

func (b *blobRepository) Download(ctx context.Context, name string) (io.Reader, error) {
	log.Debug().Msg("Getting blob")

	rsp, err := b.client.DownloadStream(ctx, b.container, name, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get blob")
		return nil, err
	}

	return rsp.Body, nil
}

func (b *blobRepository) Delete(ctx context.Context, name string) error {
	log.Debug().Msg("Deleting blob")

	_, err := b.client.DeleteBlob(ctx, b.container, name, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete blob")
		return err
	}

	return nil
}
