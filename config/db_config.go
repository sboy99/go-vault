package config

type DBConfig struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{}
}
