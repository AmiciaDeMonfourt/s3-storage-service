package s3

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage struct {
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func New() *Storage {
	custResolverFn := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		return endpoints.ResolvedEndpoint{
			URL: "https://storage.yandexcloud.net",
		}, nil
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:           aws.String("ru-central1"),
			EndpointResolver: endpoints.ResolverFunc(custResolverFn),
		},
	}))

	return &Storage{
		client:     s3.New(sess),
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
	}
}

func (s *Storage) ListBuckets() (*s3.ListBucketsOutput, error) {
	return s.client.ListBuckets(nil)
}

func (s *Storage) CreateBucket(name string) (*s3.CreateBucketOutput, error) {
	bucket, err := s.client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	err = s.client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(name)})

	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func (s *Storage) ListBucketItems(bucket string) (*s3.ListObjectsV2Output, error) {
	return s.client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
}

func (s *Storage) UploadItem(bucket, localpath, key string) (url string, err error) {
	file, err := os.Open(localpath)
	if err != nil {
		return "", fmt.Errorf("unable to open the file %s, why: %s", key, err.Error())
	}
	defer file.Close()

	uploadingInfo, err := s.uploader.UploadWithContext(context.TODO(), &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return "", fmt.Errorf("unable to upload item: %s, why: %s", key, err.Error())
	}

	return uploadingInfo.Location, nil
}

func (s *Storage) DownloadItem(bucket, item string) error {
	file, err := os.Create(item)
	if err != nil {
		return fmt.Errorf("unable to open the file: %s", err.Error())
	}
	defer file.Close()

	_, err = s.downloader.DownloadWithContext(
		context.TODO(),
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		},
	)

	if err != nil {
		return fmt.Errorf("unable to download item: %s, why: %s", item, err.Error())
	}

	return nil
}

func (s *Storage) DeleteItem(bucket, item string) error {
	_, err := s.client.DeleteObjectWithContext(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})

	if err != nil {
		return fmt.Errorf("unable to delete item: %s\nwhy: %s", item, err.Error())
	}

	err = s.client.WaitUntilObjectNotExistsWithContext(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetPresignURL(bucket, item string) (string, error) {
	req, resp := s.client.PutObjectRequest(&s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})

	err := req.Send()
	if err != nil {
		return "", fmt.Errorf("resp: %s\nerror: %s", resp, err.Error())
	}

	url, err := req.Presign(1 * time.Hour)
	if err != nil {
		return "", err
	}

	return url, nil
}
