package objectstorage

type ObjectStorageIface interface {
	Init() error
	GeneratePresignedUploadUrl(path string, durationMinute int) (string, error)
	GeneratePresignedDownloadUrl(path string, durationMinute int) (string, error)
	ListObjectsByPrefix(prefix string) ([]string, error)
}
