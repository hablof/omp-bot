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
	Kafka   Kafka   `yaml:"kafka"`
}

type App struct {
	Debug bool `yaml:"debug"`
}

type GrpcAPI struct {
	Target      string `yaml:"target"`
	Attempts    int    `yaml:"attempts"`
	DialTimeout uint64 `yaml:"dialtimeout"` // in seconds
}

type Tgbot struct {
	Debug          bool `yaml:"debug"`
	PaginationStep int  `yaml:"paginationstep"`
}

type Kafka struct {
	// Capacity uint64   `yaml:"capacity"`
	Topic string `yaml:"topic"`
	// GroupID  string   `yaml:"groupId"`
	Brokers []string `yaml:"brokers"`
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
