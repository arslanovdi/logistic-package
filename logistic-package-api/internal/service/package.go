// Package service - слой бизнес-логики
package service

import (
	"context"
	"fmt"
	"github.com/arslanovdi/logistic-package/pkg/model"
)

// Repo interface for work with database
type Repo interface {
	Create(ctx context.Context, pkg *model.Package) (*uint64, error)
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*model.Package, error)
	List(ctx context.Context, offset, limit uint64) ([]model.Package, error)
	Update(ctx context.Context, pkg *model.Package) error
}

// PackageService is service for Package
type PackageService struct {
	repo Repo
}

// Create - создание нового пакета
func (p *PackageService) Create(ctx context.Context, pkg *model.Package) (*uint64, error) {

	id, err := p.repo.Create(ctx, pkg)
	if err != nil {
		return nil, fmt.Errorf("service.PackageService.Create: %w", err)
	}

	return id, nil
}

// Delete - удаление пакета
func (p *PackageService) Delete(ctx context.Context, id uint64) error {

	err := p.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service.PackageService.Delete: %w", err)
	}

	return nil
}

// Get - получение пакета
func (p *PackageService) Get(ctx context.Context, id uint64) (*model.Package, error) {

	pkg, err := p.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.PackageService.Get: %w", err)
	}
	return pkg, nil
}

// List - получение списка пакетов
func (p *PackageService) List(ctx context.Context, offset, limit uint64) ([]model.Package, error) {

	packages, err := p.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("service.PackageService.List: %w", err)
	}
	return packages, nil
}

// Update - изменение пакета
func (p *PackageService) Update(ctx context.Context, pkg *model.Package) error {

	err := p.repo.Update(ctx, pkg)
	if err != nil {
		return fmt.Errorf("service.PackageService.Update: %w", err)
	}
	return nil
}

// NewPackageService - конструктор
func NewPackageService(repo Repo) *PackageService {
	return &PackageService{
		repo: repo,
	}
}
