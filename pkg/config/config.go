package config

type Arcitecture string

type Config struct {
	Arcitecture Arcitecture `yaml:"arcitecture"`
	Folder      string      `yaml:"folder"`
}

const (
	Hexagonal Arcitecture = "Hexagonal"
	Module    Arcitecture = "Module"
	MVC       Arcitecture = "MVC"
)
