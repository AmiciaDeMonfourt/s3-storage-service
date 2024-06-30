package repository

import (
	"bytes"
	"context"
	"s3/internal/model"
	"s3/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Object struct {
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewObject(client *s3.S3, sess *session.Session) *Object {
	return &Object{
		client:     client,
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
	}
}

func (o *Object) UploadLiteObject(ctx context.Context, req *api.UploadMediaRequest) (*api.UploadMediaResponse, error) {
	objectInfo := &s3manager.UploadInput{
		Bucket: aws.String(req.GetBucket()),
		Key:    aws.String(req.GetObjKey()),
		Body:   bytes.NewReader(req.GetData()),
	}

	out, err := o.uploader.UploadWithContext(ctx, objectInfo)
	if err != nil {
		return nil, err
	}

	response := &api.UploadMediaResponse{Url: out.Location}

	return response, nil
}

func (o *Object) UploadLargeFile() {

}

func (o *Object) Download(ctx context.Context, req *api.DownloadMediaRequest) (*api.DownloadMediaResponse, error) {
	buff := &aws.WriteAtBuffer{}

	objectInfo := &s3.GetObjectInput{
		Bucket: aws.String(req.GetBucket()),
		Key:    aws.String(req.GetKey()),
	}

	_, err := o.downloader.DownloadWithContext(ctx, buff, objectInfo)
	if err != nil {
		return nil, err
	}

	return &api.DownloadMediaResponse{
		Data: buff.Bytes(),
	}, nil
}

func (o *Object) Delete() {}

func (o *Object) List(ctx context.Context, bucket string) ([]*model.ObjectInfo, error) {
	out, err := o.client.ListObjectsWithContext(ctx, &s3.ListObjectsInput{Bucket: aws.String(bucket)})
	if err != nil {
		return nil, err
	}

	// form response with objects info
	objects := make([]*model.ObjectInfo, 0)
	for _, item := range out.Contents {
		one := &model.ObjectInfo{
			Key:  aws.StringValue(item.Key),
			Size: aws.Int64Value(item.Size),
		}

		objects = append(objects, one)
	}

	return objects, nil
}
