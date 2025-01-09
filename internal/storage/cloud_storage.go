package storage

type CloudStorage struct {
	// TODO: Add creds and other cloud storage specific fields
}

func NewCloudStorage() *CloudStorage {
	return &CloudStorage{}
}

func (c *CloudStorage) Save(filename string, data []byte) error {
	return nil
}

func (c *CloudStorage) Load(filename string) ([]byte, error) {
	return nil, nil
}

func (c *CloudStorage) Delete(filename string) error {
	return nil
}
