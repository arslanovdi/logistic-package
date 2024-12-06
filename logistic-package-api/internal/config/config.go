// Package config - contains service config
package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// Build information -ldflags .
const (
	version    string = "dev"
	commitHash string = "-"
)

var cfg *Config

// GetConfigInstance returns service config
func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

// Database - contains all parameters database connection.
type Database struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Migrations   string `yaml:"migrations"`
	Name         string `yaml:"name"`
	SslMode      string `yaml:"sslmode"`
	Driver       string `yaml:"driver"`
	QueryTimeout int    `yaml:"queryTimeout"`
}

// Grpc - contains parameter address grpc.
type Grpc struct {
	Port              int    `yaml:"port"`
	MaxConnectionIdle int64  `yaml:"maxConnectionIdle"`
	Timeout           int64  `yaml:"timeout"`
	MaxConnectionAge  int64  `yaml:"maxConnectionAge"`
	Host              string `yaml:"host"`
}

// Rest - contains parameter rest json connection.
type Rest struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

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
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	Path string `yaml:"path"`
}

// Jaeger - contains all parameters metrics information.
type Jaeger struct {
	Service string `yaml:"service"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

// Kafka - contains all parameters kafka information.
type Kafka struct {
	Capacity       int      `yaml:"capacity"`
	Topic          string   `yaml:"topic"`
	GroupID        string   `yaml:"groupId"`
	FlushTimeout   int      `yaml:"flushTimeout"`
	Brokers        []string `yaml:"brokers"`
	SchemaRegistry string   `yaml:"schemaregistry"`
}

type Outbox struct {
	BatchSize     int `yaml:"batchSize"`
	Ticker        int `yaml:"ticker"`
	ProducerCount int `yaml:"producerCount"`
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
	Project  Project  `yaml:"project"`
	Grpc     Grpc     `yaml:"grpc"`
	Rest     Rest     `yaml:"rest"`
	Database Database `yaml:"database"`
	Metrics  Metrics  `yaml:"metrics"`
	Jaeger   Jaeger   `yaml:"jaeger"`
	Kafka    Kafka    `yaml:"kafka"`
	Status   Status   `yaml:"status"`
	Outbox   Outbox   `yaml:"outbox"`
}

// ReadConfigYML - read configurations from file and init instance Config.
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
	if err2 := decoder.Decode(&cfg); err2 != nil {
		return fmt.Errorf("config.ReadConfigYML: %w", err2)
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}
