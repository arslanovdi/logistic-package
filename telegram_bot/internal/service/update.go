package service

import (
	"context"
	"github.com/arslanovdi/logistic-package/pkg/model"
)

// Update изменяем существующий пакет
func (c *LogisticPackageService) Update(pkg *model.Package) error {

	ctx, cancel := context.WithTimeout(context.Background(), c.ctxTimeout)
	defer cancel()

	err := c.api.Update(ctx, pkg)
	if err != nil {
		return err
	}

	return nil
}
