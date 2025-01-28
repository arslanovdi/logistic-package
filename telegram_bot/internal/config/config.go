// Package config работа с конфигурацией сервиса
package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	version    string = "dev"
	commitHash string = "-"
)

var cfg *Config

// Config структура конфигурации сервиса
type Config struct {
	Project  project  `yaml:"project"`
	GRPC     grpc     `yaml:"grpc"`
	Jaeger   jaeger   `yaml:"jaeger"`
	Telegram telegram `yaml:"telegram"`
	Metrics  metrics  `yaml:"metrics"`
	Status   status   `yaml:"status"`
}

type project struct {
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

// metrics - contains all parameters metrics information.
type metrics struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	Path string `yaml:"path"`
}

// status config for service.
type status struct {
	Port          int    `yaml:"port"`
	Host          string `yaml:"host"`
	VersionPath   string `yaml:"versionPath"`
	LivenessPath  string `yaml:"livenessPath"`
	ReadinessPath string `yaml:"readinessPath"`
}

type telegram struct {
	Faker   bool   `yaml:"faker"`
	Timeout int    `yaml:"timeout"`
	Token   string `yaml:"token"`
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

	grpcHost, ok := os.LookupEnv("GRPC_HOST")
	if ok {
		cfg.GRPC.Host = grpcHost
	}
	grpcPort, ok := os.LookupEnv("GRPC_PORT")
	if ok {
		cfg.GRPC.Port = grpcPort
	}

	jaegerHost, ok := os.LookupEnv("JAEGER_HOST")
	if ok {
		cfg.Jaeger.Host = jaegerHost
	}

	return nil
}
