package objectstorage

import "fmt"

func New(provider string) (ObjectStorageIface, error) {
	switch provider {
	case ProviderCloudflareR2:
		objs := &CloudflareR2{}
		objs.Init()
		return objs, nil
	default:
		return nil, fmt.Errorf("unknown object storage provider: %s", provider)
	}
}
