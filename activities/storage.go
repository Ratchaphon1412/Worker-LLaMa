package activities

import (
	"context"

	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/Ratchaphon1412/worker-llama/utils"
)

func Storage(ctx context.Context, conf configs.Config, objectName, filePath string) error {
	// Upload the MP3 file to MinIO
	err := utils.UploadMP3ToMinio(conf.MinioEndpoint, conf.MinioUserAccessKey, conf.MinioUserSecretKey, conf.MinioDefaultBucket, conf.MinioSSLEnabled, objectName, filePath)
	if err != nil {
		return err
	}

	return nil
}
