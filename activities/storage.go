package activities

import (
	"context"

	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/Ratchaphon1412/worker-llama/worker/drivers/storages"
)

func Storage(ctx context.Context, conf configs.Config, objectName, filePath string) (string, error) {
	// Upload the MP3 file to MinIO
	publicURL, err := storages.UploadMP3ToMinio(conf.MinioEndpoint, conf.MinioPublicURL, conf.MinioUserAccessKey, conf.MinioUserSecretKey, conf.MinioDefaultBucket, conf.MinioSSLEnabled, objectName, filePath)
	if err != nil {
		return "", err
	}

	return publicURL, nil
}
