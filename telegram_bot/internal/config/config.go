// Package config работа с конфигурацией сервиса
package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

var cfg *Config

// Config структура конфигурации сервиса
type Config struct {
	Project  Project  `yaml:"project"`
	GRPC     grpc     `yaml:"grpc"`
	Jaeger   jaeger   `yaml:"jaeger"`
	Telegram telegram `yaml:"telegram"`
}

type Project struct {
	Name            string `yaml:"name"`
	Debug           bool   `yaml:"debug"`
	Environment     string `yaml:"environment"`
	Instance        string `yaml:"instance"`
	StartupTimeout  int    `yaml:"startupTimeout"`
	ShutdownTimeout int    `yaml:"shutdownTimeout"`
}

type grpc struct {
	Host       string        `yaml:"host"`
	Port       string        `yaml:"port"`
	CtxTimeout time.Duration `yaml:"ctxTimeout"`
}

type jaeger struct {
	Service string `yaml:"service"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

type telegram struct {
	Token string `yaml:"token"`
}

// GetConfigInstance returns service config
func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

// ReadConfigYML reads config from file
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err1 := os.Open(filepath.Clean(filePath))
	if err1 != nil {
		return fmt.Errorf("config.ReadConfigYML: %w", err1)
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return fmt.Errorf("config.ReadConfigYML: %w", err)
	}

	return nil
}
