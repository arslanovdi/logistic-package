package test

import (
	"context"
	"testing"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/service"
	"github.com/arslanovdi/logistic-package/logistic-package-api/mocks"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/stretchr/testify/mock"
)

/*
Тест поведения сервиса logistic-package-api
*/

func TestPackageService_Create(t *testing.T) {
	t.Parallel()

	repoMock := mocks.NewRepo(t)

	packageService := service.NewPackageService(repoMock)

	id := uint64(1)

	repoMock.EXPECT().
		Create(mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("*model.Package")).
		Return(&id, nil)

	_, _ = packageService.Create(context.Background(), &model.Package{})
}

func TestPackageService_Delete(t *testing.T) {
	t.Parallel()

	repoMock := mocks.NewRepo(t)

	packageService := service.NewPackageService(repoMock)

	repoMock.EXPECT().
		Delete(mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("uint64")).
		Return(nil)

	_ = packageService.Delete(context.Background(), 1)
}

func TestPackageService_Get(t *testing.T) {
	t.Parallel()

	repoMock := mocks.NewRepo(t)

	packageService := service.NewPackageService(repoMock)

	var id *model.Package

	repoMock.EXPECT().
		Get(mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("uint64")).
		Return(id, nil)

	_, _ = packageService.Get(context.Background(), 1)
}

func TestPackageService_List(t *testing.T) {
	t.Parallel()

	repoMock := mocks.NewRepo(t)

	packageService := service.NewPackageService(repoMock)

	repoMock.EXPECT().
		List(mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64")).
		Return(nil, nil)

	_, _ = packageService.List(context.Background(), 0, 3)
}

func TestPackageService_Update(t *testing.T) {
	t.Parallel()

	repoMock := mocks.NewRepo(t)

	packageService := service.NewPackageService(repoMock)

	repoMock.EXPECT().
		Update(mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("*model.Package")).
		Return(nil)

	_ = packageService.Update(context.Background(), &model.Package{})
}
