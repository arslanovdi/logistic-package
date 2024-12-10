package service

import (
	"context"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/general"
)

// List возвращаем packages с позиции offset, количество - limit
func (c *LogisticPackageService) List(offset, limit uint64) ([]model.Package, error) {

	ctx, cancel := context.WithTimeout(context.Background(), c.ctxTimeout)
	defer cancel()

	packages, err := c.api.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	if uint64(len(packages)) < limit {
		return packages, general.ErrEndOfList
	}

	return packages, nil
}
