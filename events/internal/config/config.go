// Package config - contains service config
package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	version    string = "dev"
	commitHash string = "-"
)

var cfg *Config

// project - contains all parameters project information.
type project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string
	CommitHash  string
	Instance    string `yaml:"instance"`

	StartupTimeout  int `yaml:"startupTimeout"`
	ShutdownTimeout int `yaml:"shutdownTimeout"`
}

// metrics - contains all parameters metrics information.
type metrics struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

// jaeger - contains all parameters jaeger information.
type jaeger struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// kafka - contains all parameters kafka information.
type kafka struct {
	Topic          string   `yaml:"topic"`
	GroupID        string   `yaml:"groupId"`
	Brokers        []string `yaml:"brokers"`
	SchemaRegistry string   `yaml:"schemaRegistry"`
}

// status config for service.
type status struct {
	Port          int    `yaml:"port"`
	Host          string `yaml:"host"`
	VersionPath   string `yaml:"versionPath"`
	LivenessPath  string `yaml:"livenessPath"`
	ReadinessPath string `yaml:"readinessPath"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project project `yaml:"project"`
	Metrics metrics `yaml:"metrics"`
	Jaeger  jaeger  `yaml:"jaeger"`
	Kafka   kafka   `yaml:"kafka"`
	Status  status  `yaml:"status"`
}

// GetConfigInstance returns service config
func GetConfigInstance() *Config {
	if cfg != nil {
		return cfg
	}

	return &Config{}
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	log := slog.With("func", "ReadConfigYML")

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}

	defer func() {
		err1 := file.Close()
		if err1 != nil {
			log.Error("failed to close file", slog.String("error", err1.Error()))
		}
	}()

	decoder := yaml.NewDecoder(file)
	if err2 := decoder.Decode(&cfg); err2 != nil {
		return err2
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}
