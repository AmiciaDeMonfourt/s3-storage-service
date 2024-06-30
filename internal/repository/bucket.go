package repository

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Bucket struct {
	client *s3.S3
	object *Object
}

func NewBucket(client *s3.S3, sess *session.Session) *Bucket {
	return &Bucket{
		client: client,
		object: NewObject(client, sess),
	}
}

func (b *Bucket) Object() *Object {
	return b.object
}

func (b *Bucket) List(ctx context.Context) ([]*s3.Bucket, error) {
	out, err := b.client.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	return out.Buckets, err
}

func (b *Bucket) Create(ctx context.Context, name string) (location string, err error) {
	out, err := b.client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{Bucket: aws.String(name)})
	if err != nil {
		return "", err
	}

	err = b.client.WaitUntilBucketExistsWithContext(ctx, &s3.HeadBucketInput{Bucket: aws.String(name)})
	if err != nil {
		return "", err
	}

	return aws.StringValue(out.Location), nil
}

func (b *Bucket) Delete(ctx context.Context, bucket string) error {
	_, err := b.client.DeleteBucketWithContext(ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)})
	if err != nil {
		return err
	}

	err = b.client.WaitUntilBucketNotExistsWithContext(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)})
	if err != nil {
		return err
	}

	return nil
}
