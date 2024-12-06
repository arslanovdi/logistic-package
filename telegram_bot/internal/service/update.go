package service

import (
	"context"
	"fmt"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/model"
)

// Update изменяем существующий пакет
func (c *LogisticPackageService) Update(pkg model.Package) error {

	ctx, cancel := context.WithTimeout(context.Background(), c.ctxTimeout)
	defer cancel()

	err := c.api.Update(ctx, pkg)
	if err != nil {
		return fmt.Errorf("service.Update: %w", err)
	}

	return nil
}
