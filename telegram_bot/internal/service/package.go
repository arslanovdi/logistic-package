// Package service слой бизнес-логики
package service

import (
	"context"
	"time"

	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/config"
)

// Client интерфейс grpc клиента
type Client interface {
	// Create package
	Create(ctx context.Context, pkg *model.Package) (*uint64, error)
	// Delete package
	Delete(ctx context.Context, id uint64) error
	// Get package
	Get(ctx context.Context, id uint64) (*model.Package, error)
	// List packages
	List(ctx context.Context, offset, limit uint64) ([]model.Package, error)
	// Update package
	Update(ctx context.Context, pkg *model.Package) error
	// Close grpc client
	Close()
}

// LogisticPackageService слой бизнес-логики
type LogisticPackageService struct {
	grpc       Client        // gRPC client
	ctxTimeout time.Duration // Таймаут контекста gRPC запросов
}

// NewPackageService конструктор
func NewPackageService(grpc Client) *LogisticPackageService {
	cfg := config.GetConfigInstance()
	srv := &LogisticPackageService{
		grpc:       grpc,
		ctxTimeout: cfg.GRPC.CtxTimeout,
	}
	return srv
}
