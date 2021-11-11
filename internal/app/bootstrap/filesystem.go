package bootstrap

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/statistico/statistico-ratings/internal/app/filesystem"
)

func (c Container) FilesystemReader() filesystem.Reader {
	key := c.Config.AwsConfig.Key
	secret := c.Config.AwsConfig.Secret

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Region:      aws.String(c.Config.AwsConfig.Region),
	})

	if err != nil {
		panic(err)
	}

	return filesystem.NewS3Reader(s3.New(sess), c.Config.S3Bucket)
}
