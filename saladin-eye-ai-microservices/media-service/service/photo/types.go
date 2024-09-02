package photo

import "context"

type ObjectFile struct {
	Name        string
	DownloadUrl string
}

type PhotoServiceIface interface {
	GenerateUploadPresignedUrl(ctx context.Context, deviceId, idempotentKey string) (string, error)
	ListDate(ctx context.Context, deviceId string) ([]string, error)
	ListHourByDate(ctx context.Context, deviceId, date string) ([]string, error)
	ListObjectsByDateHour(ctx context.Context, deviceId string, date string, hour int32) ([]ObjectFile, error)
}
