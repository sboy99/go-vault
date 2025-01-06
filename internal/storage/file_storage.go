package storage

import (
	"os"
)

type FileStorage struct {
	BasePath string
}

func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{BasePath: basePath}
}

func (f *FileStorage) Save(filename string, data []byte) error {
	file, err := os.Create(f.filePath(filename))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileStorage) Load(filename string) ([]byte, error) {
	file, err := os.Open(f.filePath(filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make([]byte, 0)
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f *FileStorage) Delete(filename string) error {
	return os.Remove(f.filePath(filename))
}

func (f *FileStorage) filePath(filename string) string {
	return f.BasePath + "/" + filename
}
