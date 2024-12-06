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

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
)

// StatusServer - http сервер для мониторинга состояния приложения
type StatusServer struct {
	server *http.Server
}

// NewStatusServer - конструктор http сервера для мониторинга состояния приложения
func NewStatusServer(isReady *atomic.Value) *StatusServer {

	cfg := config.GetConfigInstance()

	statusAddr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)

	mux := http.DefaultServeMux

	mux.HandleFunc(cfg.Status.LivenessPath, livenessHandler)
	mux.HandleFunc(cfg.Status.ReadinessPath, readinessHandler(isReady))
	mux.HandleFunc(cfg.Status.VersionPath, versionHandler(&cfg))

	server := &http.Server{
		Addr:              statusAddr,
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 5,
	}

	return &StatusServer{
		server: server,
	}
}

// Start - запуск http сервера
func (s *StatusServer) Start() {

	log := slog.With("func", "StatusServer.Start")

	cfg := config.GetConfigInstance()

	statusAddr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)

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

func versionHandler(cfg *config.Config) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {

		log := slog.With("func", "versionHandler")

		data := map[string]interface{}{
			"name":        cfg.Project.Name,
			"debug":       cfg.Project.Debug,
			"environment": cfg.Project.Environment,
			"version":     cfg.Project.Version,
			"commitHash":  cfg.Project.CommitHash,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err1 := json.NewEncoder(w).Encode(data); err1 != nil {
			log.Error("Service information encoding error", slog.String("error", err1.Error()))
		}
	}
}
