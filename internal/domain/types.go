package domain

import "time"

// ConnParams holds connection settings for dump/restore engines.
type ConnParams struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
	SSLMode  string
}

// DumpResult holds metadata about a successful dump stream setup.
type DumpResult struct {
	ServerVersion string
	ServerMajor   int
	BinaryPath    string
}

// ObjectInfo describes a stored backup artifact.
type ObjectInfo struct {
	Key          string
	Size         int64
	SHA256       string
	LastModified time.Time
}

// Settings holds the subset of configuration written by the setup wizard.
type Settings struct {
	DatabaseType    DatabaseType
	DBHost          string
	DBPort          int
	DBName          string
	DBUsername      string
	DBPassword      string
	DBSSLMode       string
	StorageType     StorageType
	StorageDest     string
	CloudType       string
	AWSRegion       string
	AWSBucketName   string
	AWSAccessKeyID  string
	AWSAccessSecret string
	AWSEndpoint     string
}
