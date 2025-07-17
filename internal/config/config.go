package config

import (
	"os"

	"github.com/liuyifan1996/pkg/utils"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

func LoadConfig() *Config {
	cfg := &Config{}

	data, err := os.ReadFile("configs/app.yaml")
	if err != nil {
		utils.LogFatal("Failed to read config file: %v", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		utils.LogFatal("Failed to parse config: %v", err)
	}

	return cfg
}
