package mypackage

import (
	"context"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
				return 0, packageApi.ErrBadArgument{Argument: "MaximumVolume"}
			}
			req.MaximumVolume = float32(volume)

		case logistic.Reusable:
			// reusable, err := strconv.ParseBool(value)
			if strings.ToLower(value) == "да" {
				req.Reusable = true
			} else if strings.ToLower(value) == "нет" {
				req.Reusable = false
			} else {
				log.Debug().Msg("failed to parse reusable")
				return 0, packageApi.ErrBadArgument{Argument: "Reusable"}
			}

		default:
			return 0, packageApi.ErrBadArgument{Argument: key}
		}
	}

	if err := req.Validate(); err != nil {
		log.Debug().Err(err).Msg("PackageService.Create req validation failed")

		if err, ok := err.(pb.CreatePackageV1RequestValidationError); ok {
			return 0, packageApi.ErrBadArgument{Argument: err.Field()}
		}

		return 0, packageApi.ErrBadArgument{Argument: "unable to fetch invalid field"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := ps.grpcclient.CreatePackageV1(ctx, &req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.Internal:
				return 0, packageApi.ErrInternal
			}
		}

		return 0, err
	}

	return resp.GetID(), nil
}

// Describe implements mypackage.PackageService
func (ps *PackageService) Describe(packageID uint64) (logistic.Package, error) {

	req := &pb.DescribePackageV1Request{
		PackageID: packageID,
	}

	if err := req.Validate(); err != nil {
		log.Debug().Err(err).Msg("PackageService.Describe req validation failed")

		if err, ok := err.(pb.DescribePackageV1RequestValidationError); ok {
			return logistic.Package{}, packageApi.ErrBadArgument{Argument: err.Field()}
		}

		return logistic.Package{}, packageApi.ErrBadArgument{Argument: "unable to fetch invalid field"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := ps.grpcclient.DescribePackageV1(ctx, req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.Internal:
				return logistic.Package{}, packageApi.ErrInternal

			case codes.NotFound:
				return logistic.Package{}, packageApi.ErrNotFound
			}
		}

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
		Limit:  limit,
	}

	if err := req.Validate(); err != nil {
		log.Debug().Err(err).Msg("PackageService.List req validation failed")

		if err, ok := err.(pb.ListPackagesV1RequestValidationError); ok {
			return nil, packageApi.ErrBadArgument{Argument: err.Field()}
		}

		return nil, packageApi.ErrBadArgument{Argument: "unable to fetch invalid field"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := ps.grpcclient.ListPackagesV1(ctx, req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.Internal:
				return nil, packageApi.ErrInternal
			}
		}

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

	if err := req.Validate(); err != nil {
		log.Debug().Err(err).Msg("PackageService.Remove req validation failed")

		if err, ok := err.(pb.RemovePackageV1RequestValidationError); ok {
			return false, packageApi.ErrBadArgument{Argument: err.Field()}
		}

		return false, packageApi.ErrBadArgument{Argument: "unable to fetch invalid field"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := ps.grpcclient.RemovePackageV1(ctx, req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.Internal:
				return false, packageApi.ErrInternal

			case codes.NotFound:
				return false, packageApi.ErrNotFound
			}
		}

		return false, err
	}

	return resp.GetSuc(), nil
}

// Update implements mypackage.PackageService
func (ps *PackageService) Update(packageID uint64, editMap map[string]string) (bool, error) {

	// filed with zero-values, which ignored by api-service
	req := &pb.UpdatePackageV1Request{
		PackageID:     packageID,
		Title:         "",
		Material:      "",
		MaximumVolume: 0,
		Reusable:      nil,
	}

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
				return false, packageApi.ErrBadArgument{Argument: "MaximumVolume"}
			}
			req.MaximumVolume = float32(volume)

		case logistic.Reusable:
			// reusable, err := strconv.ParseBool(value)
			if strings.ToLower(value) == "да" {
				req.Reusable = &pb.MaybeBool{Reusable: true}
			} else if strings.ToLower(value) == "нет" {
				req.Reusable = &pb.MaybeBool{Reusable: false}
			} else {
				log.Debug().Msg("failed to parse reusable")
				return false, packageApi.ErrBadArgument{Argument: "Reusable"}
			}
		}
	}

	if err := req.Validate(); err != nil {
		log.Debug().Err(err).Msg("PackageService.Update req validation failed")

		if err, ok := err.(pb.UpdatePackageV1RequestValidationError); ok {
			return false, packageApi.ErrBadArgument{Argument: err.Field()}
		}

		return false, packageApi.ErrBadArgument{Argument: "unable to fetch invalid field"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := ps.grpcclient.UpdatePackageV1(ctx, req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.Internal:
				return false, packageApi.ErrInternal

			case codes.NotFound:
				return false, packageApi.ErrNotFound
			}
		}

		return false, err
	}

	return resp.GetSuc(), nil
}

func NewService(cc grpc.ClientConnInterface) *PackageService {
	return &PackageService{
		grpcclient: pb.NewLogisticPackageApiServiceClient(cc),
	}
}
