package mypackage

import (
	"google.golang.org/grpc"

	pb "github.com/hablof/logistic-package-api/pkg/logistic-package-api"

	"github.com/hablof/omp-bot/internal/app/commands/logistic/packageApi"
	"github.com/hablof/omp-bot/internal/model/logistic"
)

var _ packageApi.PackageService = &PackageService{}

type PackageService struct {
	grpcclient pb.LogisticPackageApiServiceClient
}

// Create implements mypackage.PackageService
func (ps *PackageService) Create(createMap map[string]string) (uint64, error) {
	// ps.grpcclient.CreatePackageV1()
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

func NewService(cc grpc.ClientConnInterface) *PackageService {
	return &PackageService{
		grpcclient: pb.NewLogisticPackageApiServiceClient(cc),
	}
}
