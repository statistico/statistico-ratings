package filesystem

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"io"
)

type Reader interface {
	Reader(filename string) (io.ReadCloser, error)
}

type s3Reader struct {
	client  s3iface.S3API
	bucket  string
}

func (s *s3Reader) Reader(filename string) (io.ReadCloser, error) {
	input := s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	}

	object, err := s.client.GetObject(&input)

	if err != nil {
		return nil, err
	}

	return object.Body, nil
}

func NewS3Reader(c s3iface.S3API, b string) Reader {
	return &s3Reader{
		client: c,
		bucket: b,
	}
}
