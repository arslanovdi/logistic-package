package service

import (
	"context"
	"github.com/arslanovdi/logistic-package/pkg/model"
)

// Get возвращаем package с id: cursor
func (c *LogisticPackageService) Get(cursor uint64) (model.Package, error) {

	ctx, cancel := context.WithTimeout(context.Background(), c.ctxTimeout)
	defer cancel()

	pkg, err := c.api.Get(ctx, cursor)
	if err != nil {
		return model.Package{}, err
	}

	return *pkg, nil
}
