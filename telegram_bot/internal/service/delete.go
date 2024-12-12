package service

import (
	"context"
)

// Delete package
func (c *LogisticPackageService) Delete(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.ctxTimeout)
	defer cancel()

	err := c.grpc.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
