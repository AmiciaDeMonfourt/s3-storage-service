package s3

import (
	"s3/internal/repository"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage interface {
	Bucket()
}

type Storage struct {
	bucket *repository.Bucket
}

func New() *Storage {
	custResolverFn := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		return endpoints.ResolvedEndpoint{
			URL: "https://storage.yandexcloud.net",
		}, nil
	}

	config := aws.Config{
		Region:           aws.String("ru-central1"),
		EndpointResolver: endpoints.ResolverFunc(custResolverFn),
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{Config: config}))

	client := s3.New(sess)

	return &Storage{
		bucket: repository.NewBucket(client, sess),
	}
}

func (s *Storage) Bucket() *repository.Bucket {
	return s.bucket
}
