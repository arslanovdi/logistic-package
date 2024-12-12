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

// GetConfigInstance returns service config
func GetConfigInstance() *Config {
	if cfg != nil {
		return cfg
	}

	return &Config{}
}

// database - contains all parameters database connection.
type database struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Migrations   string `yaml:"migrations"`
	Name         string `yaml:"name"`
	Ssl          string `yaml:"ssl"`
	Driver       string `yaml:"driver"`
	QueryTimeout int    `yaml:"queryTimeout"`
}

// grpc - contains parameter address grpc.
type grpc struct {
	Port              int    `yaml:"port"`
	MaxConnectionIdle int64  `yaml:"maxConnectionIdle"`
	Timeout           int64  `yaml:"timeout"`
	MaxConnectionAge  int64  `yaml:"maxConnectionAge"`
	Host              string `yaml:"host"`
}

// rest - contains parameter rest json connection.
type rest struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

// project - contains all parameters project information.
type project struct {
	Debug           bool   `yaml:"debug"`
	Name            string `yaml:"name"`
	Environment     string `yaml:"environment"`
	Version         string
	CommitHash      string
	Instance        string `yaml:"instance"`
	StartupTimeout  int    `yaml:"startupTimeout"`
	ShutdownTimeout int    `yaml:"shutdownTimeout"`
}

// metrics - contains all parameters metrics information.
type metrics struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	Path string `yaml:"path"`
}

// jaeger - contains all parameters jaeger information.
type jaeger struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// kafka - contains all parameters kafka information.
type kafka struct {
	Capacity       int      `yaml:"capacity"`
	Topic          string   `yaml:"topic"`
	GroupID        string   `yaml:"groupId"`
	FlushTimeout   int      `yaml:"flushTimeout"`
	Brokers        []string `yaml:"brokers"`
	SchemaRegistry string   `yaml:"schemaRegistry"`
}

type outbox struct {
	BatchSize     int `yaml:"batchSize"`
	Ticker        int `yaml:"ticker"`
	ProducerCount int `yaml:"producerCount"`
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
	Project  project  `yaml:"project"`
	Grpc     grpc     `yaml:"grpc"`
	Rest     rest     `yaml:"rest"`
	Database database `yaml:"database"`
	Metrics  metrics  `yaml:"metrics"`
	Jaeger   jaeger   `yaml:"jaeger"`
	Kafka    kafka    `yaml:"kafka"`
	Status   status   `yaml:"status"`
	Outbox   outbox   `yaml:"outbox"`
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
