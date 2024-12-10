// Package config - contains service config
package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var cfg *Config

// Project - contains all parameters project information.
type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string
	CommitHash  string
	Instance    string `yaml:"instance"`

	StartupTimeout  int `yaml:"startupTimeout"`
	ShutdownTimeout int `yaml:"shutdownTimeout"`
}

// Metrics - contains all parameters metrics information.
type Metrics struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

// Jaeger - contains all parameters jaeger information.
type Jaeger struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Kafka - contains all parameters kafka information.
type Kafka struct {
	Topic          string   `yaml:"topic"`
	GroupID        string   `yaml:"groupId"`
	Brokers        []string `yaml:"brokers"`
	SchemaRegistry string   `yaml:"schemaRegistry"`
}

// Status config for service.
type Status struct {
	Port          int    `yaml:"port"`
	Host          string `yaml:"host"`
	VersionPath   string `yaml:"versionPath"`
	LivenessPath  string `yaml:"livenessPath"`
	ReadinessPath string `yaml:"readinessPath"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project Project `yaml:"project"`
	Metrics Metrics `yaml:"metrics"`
	Jaeger  Jaeger  `yaml:"jaeger"`
	Kafka   Kafka   `yaml:"kafka"`
	Status  Status  `yaml:"status"`
}

// GetConfigInstance returns service config
func GetConfigInstance() *Config {
	if cfg != nil {
		return cfg
	}

	return &Config{}
}

// ReadConfigYML - read configurations from file and init instance Config.
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
