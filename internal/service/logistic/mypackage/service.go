package mypackage

import (
	"github.com/hablof/omp-bot/internal/app/commands/logistic/packageApi"
	"github.com/hablof/omp-bot/internal/model/logistic"
)

var _ packageApi.PackageService = &PackageService{}

type PackageService struct{}

// Create implements mypackage.PackageService
func (*PackageService) Create(logistic.Package) (uint64, error) {
	panic("unimplemented")
}

// Describe implements mypackage.PackageService
func (*PackageService) Describe(packageID uint64) (*logistic.Package, error) {
	panic("unimplemented")
}

// List implements mypackage.PackageService
func (*PackageService) List(cursor uint64, limit uint64) ([]logistic.Package, error) {
	panic("unimplemented")
}

// Remove implements mypackage.PackageService
func (*PackageService) Remove(packageID uint64) (bool, error) {
	panic("unimplemented")
}

// Update implements mypackage.PackageService
func (*PackageService) Update(packageID uint64, editMap map[string]string) error {
	panic("unimplemented")
}

func NewService() *PackageService {
	return &PackageService{}
}
