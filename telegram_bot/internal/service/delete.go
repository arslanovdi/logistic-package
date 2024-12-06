package service

import (
	"context"
	"fmt"
)

// Delete удаляем пакет с id: cursor
func (c *LogisticPackageService) Delete(cursor uint64) error {

	ctx, cancel := context.WithTimeout(context.Background(), c.ctxTimeout)
	defer cancel()

	err := c.api.Delete(ctx, cursor)
	if err != nil {
		return fmt.Errorf("service.Delete: %w", err)
	}

	return nil
}
