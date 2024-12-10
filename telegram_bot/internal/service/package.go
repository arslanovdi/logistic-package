// Package service слой бизнес-логики
package service

import (
	"context"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/config"
	"time"
)

// Client интерфейс grpc клиента
type Client interface {
	Create(ctx context.Context, pkg *model.Package) (*uint64, error)
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*model.Package, error)
	List(ctx context.Context, offset, limit uint64) ([]model.Package, error)
	Update(ctx context.Context, pkg *model.Package) error
	Close()
}

// LogisticPackageService слой бизнес-логики
type LogisticPackageService struct {
	api        Client        // gRPC клиент
	ctxTimeout time.Duration // Таймаут контекста gRPC запросов
}

// NewPackageService инициализирует слой бизнес-логики
func NewPackageService(grpc Client) *LogisticPackageService {
	cfg := config.GetConfigInstance()
	srv := &LogisticPackageService{
		api:        grpc,
		ctxTimeout: cfg.GRPC.CtxTimeout,
	}
	return srv
}
