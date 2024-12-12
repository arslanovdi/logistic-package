package service

import (
	"context"

	"github.com/arslanovdi/logistic-package/pkg/model"
)

// Get package
func (c *LogisticPackageService) Get(id uint64) (model.Package, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.ctxTimeout)
	defer cancel()

	pkg, err := c.grpc.Get(ctx, id)
	if err != nil {
		return model.Package{}, err
	}

	return *pkg, nil
}
