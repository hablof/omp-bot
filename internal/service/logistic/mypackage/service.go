package mypackage

import (
	"fmt"

	"github.com/hablof/omp-bot/internal/model/logistic"
)

type PackageService interface {
	Describe(packageID uint64) (*logistic.Package, error)
	List(cursor uint64, limit uint64) ([]logistic.Package, error)
	Create(logistic.Package) (uint64, error)
	Update(packageID uint64, mypackage logistic.Package) error
	Remove(packageID uint64) (bool, error)
}

var _ PackageService = &DummyPackageService{}

type DummyPackageService struct{}

// Create implements PackageService
func (*DummyPackageService) Create(unit logistic.Package) (uint64, error) {

	logistic.AllEntities = append(logistic.AllEntities, &unit)

	return uint64(len(logistic.AllEntities)), nil
}

// Describe implements PackageService
func (*DummyPackageService) Describe(packageID uint64) (*logistic.Package, error) {

	if !logistic.CheckInbounds(packageID) {
		return nil, fmt.Errorf("index out of bounds")
	}

	return logistic.AllEntities[packageID-1], nil
}

// List implements PackageService
func (*DummyPackageService) List(cursor uint64, limit uint64) ([]logistic.Package, error) {

	output := make([]logistic.Package, 0, limit)

	if !logistic.CheckInbounds(cursor) {
		return nil, fmt.Errorf("index out of bounds")
	}

	var rightBorder uint64
	if cursor-1+limit > uint64(len(logistic.AllEntities)) {
		rightBorder = uint64(len(logistic.AllEntities))
	} else {
		rightBorder = cursor - 1 + limit
	}

	for i := cursor - 1; i < rightBorder; i++ {
		output = append(output, *logistic.AllEntities[i])
	}

	return output, nil
}

// Remove implements PackageService
func (*DummyPackageService) Remove(packageID uint64) (bool, error) {
	if !logistic.CheckInbounds(packageID) {
		return false, fmt.Errorf("index out of bounds")
	}

	logistic.AllEntities = append(logistic.AllEntities[:packageID-1], logistic.AllEntities[packageID:]...)
	return true, nil
}

// Update implements PackageService
func (*DummyPackageService) Update(packageID uint64, mypackage logistic.Package) error {
	if !logistic.CheckInbounds(packageID) {
		return fmt.Errorf("index out of bounds")
	}

	logistic.AllEntities[packageID-1] = &mypackage
	return nil
}

func NewService() *DummyPackageService {
	return &DummyPackageService{}
}
