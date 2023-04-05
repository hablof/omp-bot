package mypackage

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc"

	pb "github.com/hablof/logistic-package-api/pkg/logistic-package-api"
	"github.com/rs/zerolog/log"

	"github.com/hablof/omp-bot/internal/app/commands/logistic/packageApi"
	"github.com/hablof/omp-bot/internal/model/logistic"
)

const defaultTimeout = 10 * time.Second

var _ packageApi.PackageService = &PackageService{}

type PackageService struct {
	grpcclient pb.LogisticPackageApiServiceClient
}

// Create implements mypackage.PackageService
func (ps *PackageService) Create(createMap map[string]string) (uint64, error) {

	req := pb.CreatePackageV1Request{}
	for key, value := range createMap {
		switch key {
		case logistic.Title:
			req.Title = value

		case logistic.Material:
			req.Material = value

		case logistic.MaximumVolume:
			volume, err := strconv.ParseFloat(value, 32)
			if err != nil {
				log.Debug().Err(err).Msg("failed to parse volume")
				return 0, packageApi.ErrBadRequest
			}
			req.MaximumVolume = float32(volume)

		case logistic.Reusable:
			reusable, err := strconv.ParseBool(value)
			if err != nil {
				log.Debug().Err(err).Msg("failed to parse reusable")
				return 0, packageApi.ErrBadRequest
			}
			req.Reusable = reusable

		default:
			return 0, packageApi.ErrBadRequest
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := ps.grpcclient.CreatePackageV1(ctx, &req)
	if err != nil {
		return 0, err
	}

	return resp.GetID(), nil
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
