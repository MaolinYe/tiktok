package minio

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

func UploadFileByPath(bucketName, objectName, path, contentType string) (int64, error) {
	if len(bucketName) <= 0 || len(objectName) <= 0 || len(path) <= 0 {
		return -1, errors.New("invalid argument")
	}
	uploadInfo, err := client.FPutObject(context.Background(), bucketName, objectName, path, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return -1, err
	}
	return uploadInfo.Size, nil
}

func UploadFileByIO(bucketName, objectName string, reader io.Reader, size int64, contentType string) (int64, error) {
	if len(bucketName) <= 0 || len(objectName) <= 0 {
		return -1, errors.New("invalid argument")
	}
	uploadInfo, err := client.PutObject(context.Background(), bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return -1, err
	}
	return uploadInfo.Size, nil
}

func GetFileTemporaryURL(bucketName, objectName string) (string, error) {
	if len(bucketName) <= 0 || len(objectName) <= 0 {
		return "", errors.New("invalid argument")
	}
	expiry := time.Second * time.Duration(expireTime)
	presignedURL, err := client.PresignedGetObject(context.Background(), bucketName, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
