package storages

import (
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func UploadMP3ToMinio(endpoint string, publicEndpoint string, accessKeyID string, secretAccessKey string, bucketName string, secure bool, objectName string, filePath string) (string, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return "", err
	}

	// Ensure the bucket exists
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return "", err
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: ""})
		if err != nil {
			return "", err
		}
	}

	// Upload the MP3 file
	uploadInfo, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: "audio/mpeg",
	})
	if err != nil {
		return "", err
	}

	log.Printf("Successfully uploaded %s of size %d\n", uploadInfo.Key, uploadInfo.Size)

	log.Printf("Object %s is now publicly accessible\n", objectName)

	publicURL := publicEndpoint + "/" + bucketName + "/" + objectName

	return publicURL, nil
}
