// Package config работа с конфигурацией сервиса
package config

import (
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
	Metrics  Metrics  `yaml:"metrics"`
	Status   Status   `yaml:"status"`
}

type Project struct {
	Name            string `yaml:"name"`
	Debug           bool   `yaml:"debug"`
	Environment     string `yaml:"environment"`
	Version         string
	CommitHash      string
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
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Metrics - contains all parameters metrics information.
type Metrics struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	Path string `yaml:"path"`
}

// Status config for service.
type Status struct {
	Port          int    `yaml:"port"`
	Host          string `yaml:"host"`
	VersionPath   string `yaml:"versionPath"`
	LivenessPath  string `yaml:"livenessPath"`
	ReadinessPath string `yaml:"readinessPath"`
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
func ReadConfigYML(filePath string) (err error) {
	if cfg != nil {
		return nil
	}

	file, err1 := os.Open(filepath.Clean(filePath))
	if err1 != nil {
		return err1
	}

	decoder := yaml.NewDecoder(file)
	if err2 := decoder.Decode(&cfg); err2 != nil {
		return err2
	}

	err3 := file.Close()
	if err3 != nil {
		return err3
	}

	return nil
}
