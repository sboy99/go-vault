package strategies

type StorageEnum string

const (
	FILE  StorageEnum = "file"
	CLOUD StorageEnum = "cloud"
)

// IStorage defines the interface for a storage backend.
type IStorage interface {
	Save(filename string, data []byte) error
	Load(filename string) ([]byte, error)
	Delete(filename string) error
}

type Storage struct {
	storageMap map[StorageEnum]IStorage
}

func NewStorage(path string) *Storage {
	return &Storage{
		storageMap: map[StorageEnum]IStorage{
			FILE:  NewFileStorage(path),
			CLOUD: nil,
		},
	}
}

func (s *Storage) Save(storageType StorageEnum, filename string, data []byte) error {
	storage := s.GetStorage(storageType)
	return storage.Save(filename, data)
}

func (s *Storage) Load(storageType StorageEnum, filename string) ([]byte, error) {
	storage := s.GetStorage(storageType)
	return storage.Load(filename)
}

func (s *Storage) Delete(storageType StorageEnum, filename string) error {
	storage := s.GetStorage(storageType)
	return storage.Delete(filename)
}

func (s *Storage) GetStorage(storageType StorageEnum) IStorage {
	return s.storageMap[storageType]
}
