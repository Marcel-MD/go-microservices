package repositories

import (
	"context"
	"file/config"
	"fmt"
	"io"
	"net/url"
	"sync"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/rs/zerolog/log"
)

type IBlobRepository interface {
	Get(ctx context.Context, fileName string) (io.Reader, error)
	Upload(ctx context.Context, fileName string, reader io.Reader) (string, error)
	Delete(ctx context.Context, fileName string) error
}

type blobRepository struct {
	containerUrl azblob.ContainerURL
}

var (
	blobOnce sync.Once
	blobRepo IBlobRepository
)

func GetBlobRepository() IBlobRepository {
	blobOnce.Do(func() {
		log.Info().Msg("Initializing blob service")

		cfg := config.GetConfig()

		credential, err := azblob.NewSharedKeyCredential(cfg.AzureName, cfg.AzureKey)
		if err != nil {
			log.Fatal().Err(err).Msg("Invalid credentials with error")
		}
		pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

		URL, err := url.Parse(fmt.Sprintf("%s/%s", cfg.AzureEndpoint, cfg.AzureContainer))
		if err != nil {
			log.Fatal().Err(err).Msg("Invalid URL with error")
		}

		containerURL := azblob.NewContainerURL(*URL, pipeline)

		blobRepo = &blobRepository{
			containerUrl: containerURL,
		}
	})

	return blobRepo
}

func (r *blobRepository) Get(ctx context.Context, fileName string) (io.Reader, error) {
	blobUrl := r.containerUrl.NewBlockBlobURL(fileName)
	downloadResponse, err := blobUrl.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, err
	}

	return downloadResponse.Body(azblob.RetryReaderOptions{}), nil
}

func (r *blobRepository) Upload(ctx context.Context, fileName string, reader io.Reader) (string, error) {
	blobUrl := r.containerUrl.NewBlockBlobURL(fileName)
	_, err := azblob.UploadStreamToBlockBlob(ctx, reader, blobUrl, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (r *blobRepository) Delete(ctx context.Context, fileName string) error {
	blobUrl := r.containerUrl.NewBlockBlobURL(fileName)
	_, err := blobUrl.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		return err
	}

	return nil
}
