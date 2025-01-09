package storage

import "github.com/sboy99/go-vault/config"

// IStorage defines the interface for a storage backend.
type IStorage interface {
	Save(filename string, data []byte) error
	Load(filename string) ([]byte, error)
	Delete(filename string) error
}

type Storage struct {
	storageMap map[config.StorageEnum]IStorage
}

func NewStorage() *Storage {
	return &Storage{
		storageMap: map[config.StorageEnum]IStorage{
			config.LOCAL: NewLocalStorage(),
			config.CLOUD: NewCloudStorage(),
		},
	}
}

func (s *Storage) Save(storageType config.StorageEnum, filename string, data []byte) error {
	storage := s.getStorage(storageType)
	return storage.Save(filename, data)
}

func (s *Storage) Load(storageType config.StorageEnum, filename string) ([]byte, error) {
	storage := s.getStorage(storageType)
	return storage.Load(filename)
}

func (s *Storage) Delete(storageType config.StorageEnum, filename string) error {
	storage := s.getStorage(storageType)
	return storage.Delete(filename)
}

func (s *Storage) getStorage(storageType config.StorageEnum) IStorage {
	return s.storageMap[storageType]
}
