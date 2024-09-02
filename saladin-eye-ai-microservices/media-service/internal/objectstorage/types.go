package objectstorage

import "context"

type ObjectStorageIface interface {
	Init() error
	GeneratePresignedUploadUrl(ctx context.Context, path string, durationMinute int) (string, error)
	GeneratePresignedDownloadUrl(ctx context.Context, path string, durationMinute int) (string, error)
	ListObjectsByPrefix(ctx context.Context, prefix string) ([]string, error)
	ListDate(ctx context.Context, deviceId string) ([]string, error)
	ListHourByDate(ctx context.Context, deviceId, date string) ([]string, error)
}
