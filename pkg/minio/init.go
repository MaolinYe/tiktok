package minio

import (
	"tiktok/pkg/viper"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	client                    *minio.Client
	config                    = viper.InitConfig()
	endPoint                  = config.GetString("minio.endpoint")
	accessKey                 = config.GetString("minio.accessKey")
	secretKey                 = config.GetString("minio.secretKey")
	useSSL                    = config.GetBool("minio.useSSL")
	VideoBucketName           = config.GetString("minio.videoBucketName")
	CoverBucketName           = config.GetString("minio.coverBucketName")
	AvatarBucketName          = config.GetString("minio.avatarBucketName")
	BackgroundImageBucketName = config.GetString("minio.backgroundBucketName")
	expireTime                = config.GetUint32("minio.expireTime")
)

func init() {
	cli, err := minio.New(endPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}
	client = cli
}
