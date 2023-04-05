package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App     App     `yaml:"app"`
	GrpcAPI GrpcAPI `yaml:"grpcapi"`
	Tgbot   Tgbot   `yaml:"tgbot"`
}

type App struct {
	Debug bool `yaml:"debug"`
}

type GrpcAPI struct {
	Target   string `yaml:"target"`
	Attempts int    `yaml:"attempts"`
}

type Tgbot struct {
	Debug bool `yaml:"debug"`
}

func ReadConfigYML(filePath string) (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
