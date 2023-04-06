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
func (ps *PackageService) Describe(packageID uint64) (logistic.Package, error) {

	req := &pb.DescribePackageV1Request{
		PackageID: packageID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := ps.grpcclient.DescribePackageV1(ctx, req)
	if err != nil {
		return logistic.Package{}, err
	}

	unit := resp.GetValue()

	return logistic.Package{
		ID:            packageID,
		Title:         unit.GetTitle(),
		Material:      unit.GetMaterial(),
		MaximumVolume: unit.GetMaximumVolume(),
		Reusable:      unit.GetReusable(),
	}, nil
}

// List implements mypackage.PackageService
func (ps *PackageService) List(offset uint64, limit uint64) ([]logistic.Package, error) {
	req := &pb.ListPackagesV1Request{
		Offset: offset,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := ps.grpcclient.ListPackagesV1(ctx, req)
	if err != nil {
		return nil, err
	}

	packages := resp.GetPackages()
	out := make([]logistic.Package, 0, len(packages))

	for _, pack := range packages {
		unit := logistic.Package{
			ID:            pack.GetID(),
			Title:         pack.GetTitle(),
			Material:      pack.GetMaterial(),
			MaximumVolume: pack.GetMaximumVolume(),
			Reusable:      pack.GetReusable(),
		}

		out = append(out, unit)
	}

	return out, nil
}

// Remove implements mypackage.PackageService
func (ps *PackageService) Remove(packageID uint64) (bool, error) {

	req := &pb.RemovePackageV1Request{
		PackageID: packageID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := ps.grpcclient.RemovePackageV1(ctx, req)
	if err != nil {
		return false, err
	}

	return resp.GetSuc(), nil
}

// Update implements mypackage.PackageService
func (ps *PackageService) Update(packageID uint64, editMap map[string]string) (bool, error) {

	// filed with zero-values, which ignored by api-service
	req := &pb.UpdatePackageV1Request{}

	for key, value := range editMap {
		switch key {
		case logistic.Title:
			req.Title = value

		case logistic.Material:
			req.Material = value

		case logistic.MaximumVolume:
			volume, err := strconv.ParseFloat(value, 32)
			if err != nil {
				log.Debug().Err(err).Msg("failed to parse volume")
				return false, packageApi.ErrBadRequest
			}
			req.MaximumVolume = float32(volume)

		case logistic.Reusable:
			reusable, err := strconv.ParseBool(value)
			if err != nil {
				log.Debug().Err(err).Msg("failed to parse reusable")
				return false, packageApi.ErrBadRequest
			}
			req.Reusable = &pb.MaybeBool{Reusable: reusable}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := ps.grpcclient.UpdatePackageV1(ctx, req)
	if err != nil {
		return false, err
	}

	return resp.GetSuc(), nil
}

func NewService(cc grpc.ClientConnInterface) *PackageService {
	return &PackageService{
		grpcclient: pb.NewLogisticPackageApiServiceClient(cc),
	}
}
