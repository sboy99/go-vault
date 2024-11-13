package config

import (
	"path/filepath"

	"github.com/sboy99/go-nester/pkg/utils"
	"gopkg.in/yaml.v2"
)

// -----------------------------------TYPES----------------------------------------- //

type Arcitecture string

type Config struct {
	Arcitecture Arcitecture `yaml:"arcitecture"`
	Folder      string      `yaml:"folder"`
}

// -----------------------------------VARIABLES----------------------------------------- //

const _CONFIG_FILE string = "nester.yaml"
const (
	HEXAGONAL Arcitecture = "Hexagonal"
	MODULE    Arcitecture = "Module"
	MVC       Arcitecture = "MVC"
)

var _config *Config

// -----------------------------------PUBLIC----------------------------------------- //

func NewConfig() *Config {
	return &Config{}
}

func GetConfig() *Config {
	return _config
}

func LoadConfig() (*Config, error) {
	config := GetConfig()
	err := config.Load()
	if err != nil {
		return nil, err
	}
	return config, nil
}

// -----------------------------------STRUCT_METHODS----------------------------------------- //

func (c *Config) Save() error {
	rootPath, _ := filepath.Abs(".")
	configPath := filepath.Join(rootPath, _CONFIG_FILE)
	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	err = utils.WriteFile(configPath, string(yamlBytes))
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Load() error {
	rootPath, _ := filepath.Abs(".")
	configPath := filepath.Join(rootPath, _CONFIG_FILE)
	yamlBytes, err := utils.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(yamlBytes), c)
	return err
}

// -----------------------------------PRIVATE----------------------------------------- //

func init() {
	_config = &Config{}
}
