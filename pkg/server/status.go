package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

// StatusServer - http сервер для мониторинга состояния приложения
type StatusServer struct {
	server *http.Server
	config StatusConfig
	info   ProjectInfo
}

type StatusConfig struct {
	Host          string
	Port          int
	LivenessPath  string
	ReadinessPath string
	VersionPath   string
}

type ProjectInfo struct {
	Name        string
	Debug       bool
	Environment string
	Version     string
	CommitHash  string
	Instance    string
}

// NewStatusServer - конструктор http сервера для мониторинга состояния приложения
func NewStatusServer(isReady *atomic.Value, cfg StatusConfig, projectInfo ProjectInfo) *StatusServer {

	statusAddr := fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)

	mux := http.DefaultServeMux

	mux.HandleFunc(cfg.LivenessPath, livenessHandler)
	mux.HandleFunc(cfg.ReadinessPath, readinessHandler(isReady))
	mux.HandleFunc(cfg.VersionPath, versionHandler(projectInfo))

	server := &http.Server{
		Addr:              statusAddr,
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 5,
	}

	return &StatusServer{
		server: server,
		config: cfg,
		info:   projectInfo,
	}
}

// Start - запуск http сервера
func (s *StatusServer) Start() {

	log := slog.With("func", "StatusServer.Start")

	statusAddr := fmt.Sprintf("%s:%v", s.config.Host, s.config.Port)

	go func() {
		log.Info("Status server is running", slog.String("address", statusAddr))
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed running status server", slog.String("error", err.Error()))

			os.Exit(1) // приложение завершается с ошибкой, при ошибке запуска сервера
		}
	}()
}

// Stop - остановка http сервера
func (s *StatusServer) Stop(ctx context.Context) {

	log := slog.With("func", "StatusServer.Stop")

	if err1 := s.server.Shutdown(ctx); err1 != nil {
		log.Error("StatusServer.Shutdown", slog.String("error", err1.Error()))
	} else {
		log.Info("StatusServer shut down correctly")
	}
}

func livenessHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func versionHandler(projectInfo ProjectInfo) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {

		log := slog.With("func", "versionHandler")

		data := map[string]interface{}{
			"name":        projectInfo.Name,
			"debug":       projectInfo.Debug,
			"environment": projectInfo.Environment,
			"version":     projectInfo.Version,
			"commitHash":  projectInfo.CommitHash,
			"instance":    projectInfo.Instance,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err1 := json.NewEncoder(w).Encode(data); err1 != nil {
			log.Error("Service information encoding error", slog.String("error", err1.Error()))
		}
	}
}
