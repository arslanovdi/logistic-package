package service

import (
	"context"

	"github.com/arslanovdi/logistic-package/pkg/model"
)

// Create package
func (c *LogisticPackageService) Create(pkg *model.Package) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.ctxTimeout)
	defer cancel()

	id, err := c.grpc.Create(ctx, pkg)
	if err != nil {
		return 0, err
	}

	return *id, nil
}
