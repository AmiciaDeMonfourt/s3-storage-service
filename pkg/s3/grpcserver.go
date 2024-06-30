package grpcserver

import (
	"context"
	"s3/internal/s3"
	"s3/pkg/api"
)

type GRPCServer struct {
	s3store *s3.Storage
	api.UnimplementedS3StorageServer
}

func NewGRPSServer(s3store *s3.Storage) *GRPCServer {
	return &GRPCServer{
		s3store: s3store,
	}
}

func (s *GRPCServer) UploadMedia(ctx context.Context, in *api.UploadMediaRequest) (*api.UploadMediaResponse, error) {
	// wc := s.s3store.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	return nil, nil
}
