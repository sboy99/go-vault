package storage

import (
	"os"

	"github.com/sboy99/go-vault/config"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage() *LocalStorage {
	cfg := config.GetConfig()
	return &LocalStorage{BasePath: cfg.Storage.Dest}
}

func (ls *LocalStorage) Save(filename string, data []byte) error {
	file, err := os.Create(ls.filePath(filename))
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

func (ls *LocalStorage) Load(filename string) ([]byte, error) {
	file, err := os.Open(ls.filePath(filename))
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

func (ls *LocalStorage) Delete(filename string) error {
	return os.Remove(ls.filePath(filename))
}

func (ls *LocalStorage) filePath(filename string) string {
	return ls.BasePath + "/" + filename
}
