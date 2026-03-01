package domain

import "context"

type StorageProvider interface {
	InitializeInfrastructure(ctx context.Context) error
	GenerateUploadURL(ctx context.Context, objectKey string, contentType string) (string, error)
	GenerateDownloadURL(ctx context.Context, objectKey string, inline bool) (string, error)
	MoveObject(ctx context.Context, sourceKey, destKey string) error
	DeleteObject(ctx context.Context, objectKey string) error
	Ping(ctx context.Context) error
}

var (
	ErrObjectNotFound = "document not found"
	ErrProviderFailed = "storage provider operation failed"
)
