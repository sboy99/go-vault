package storage

import "github.com/sboy99/go-vault/config"

type ICloudStorage interface {
	Upload(filename string, data []byte) error
	Download(filename string) ([]byte, error)
	Delete(filename string) error
}

type CloudStorage struct {
	Type            config.CloudEnum
	cloudStorageMap map[config.CloudEnum]ICloudStorage
}

func NewCloudStorage() *CloudStorage {
	cfg := config.GetConfig()
	return &CloudStorage{
		Type: cfg.Storage.Cloud.Type,
		cloudStorageMap: map[config.CloudEnum]ICloudStorage{
			config.AWS: NewAWSCloudStorage(),
			config.GCP: nil,
		},
	}
}

func (c *CloudStorage) Save(filename string, data []byte) error {
	return c.getCloudStorage().Upload(filename, data)
}

func (c *CloudStorage) Load(filename string) ([]byte, error) {
	return c.getCloudStorage().Download(filename)
}

func (c *CloudStorage) Delete(filename string) error {
	return nil
}

func (c *CloudStorage) getCloudStorage() ICloudStorage {
	return c.cloudStorageMap[c.Type]
}
