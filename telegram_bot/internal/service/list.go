package service

import (
	"context"
	"fmt"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/model"
)

// List возвращаем packages с позиции offset, количество - limit
func (c *LogisticPackageService) List(offset uint64, limit uint64) ([]model.Package, error) {

	ctx, cancel := context.WithTimeout(context.Background(), c.ctxTimeout)
	defer cancel()

	packages, err := c.api.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("service.List: %w", err)
	}

	if uint64(len(packages)) < limit {
		return packages, model.ErrEndOfList
	}

	return packages, nil
}
