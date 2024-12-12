// Package repo - работа с событиями в БД
package repo

import (
	"context"

	"github.com/arslanovdi/logistic-package/pkg/model"
)

// EventRepo - интерфейс работы с БД событий.
type EventRepo interface {
	// Lock заблокировать в БД n записей
	Lock(ctx context.Context, n int) ([]model.PackageEvent, error)
	// Unlock разблокировать в БД n записей
	Unlock(ctx context.Context, eventID []int64) error
	// Remove удалить из БД n записей
	Remove(ctx context.Context, eventIDs []int64) error
	// UnlockAll разблокировать в БД все записи
	// используется при инициализации ретранслятора
	UnlockAll(ctx context.Context) error
}
