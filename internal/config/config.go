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
	Redis   Redis   `yaml:"redis"`
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
	// GroupID  string   `yaml:"groupId"`
	TgCommandTopic  string   `yaml:"tgCommandTopic"`
	CacheEventTopic string   `yaml:"cacheEventTopic"`
	Brokers         []string `yaml:"brokers"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
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
