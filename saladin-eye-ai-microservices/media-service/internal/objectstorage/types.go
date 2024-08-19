package objectstorage

import "context"

type ObjectStorageIface interface {
	Init() error
	GeneratePresignedUploadUrl(ctx context.Context, path string, durationMinute int) (string, error)
	GeneratePresignedDownloadUrl(ctx context.Context, path string, durationMinute int) (string, error)
	ListObjectsByPrefix(ctx context.Context, prefix string) ([]string, error)
}
