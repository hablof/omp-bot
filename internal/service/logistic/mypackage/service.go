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
	"github.com/hablof/omp-bot/internal/model"
	"github.com/hablof/omp-bot/internal/model/logistic"
)

const defaultTimeout = 10 * time.Second

var _ packageApi.PackageService = &PackageService{}

type CacheDict interface {
	SetDescription(unit logistic.Package) error
	ReadDescription(id uint64) (*logistic.Package, error)
	RemoveDescription(id uint64) error
}

type CacheEventSender interface {
	SendCacheEvent(event model.CacheEvent) error
}

type PackageService struct {
	grpcclient pb.LogisticPackageApiServiceClient

	// cache policy:
	// on Create:   set new description;
	// on Describe: try read cache, if failed - set new description;
	// on Remove:   remove cache entry;
	// on Update:   remove cache entry;
	// on List:     nothing;
	cache            CacheDict
	cacheEventSender CacheEventSender // kafka
}

func (ps *PackageService) asyncSendEvent(packageID uint64, eventType model.CacheEventType, success bool) {

	timestamp := time.Now()

	go func(packageID uint64, eventType model.CacheEventType, success bool, timestamp time.Time) {
		err := ps.cacheEventSender.SendCacheEvent(model.CacheEvent{
			PackageID: packageID,
			EventType: eventType,
			Success:   success,
			Timestamp: timestamp,
		})
		if err != nil {
			log.Debug().Err(err).Msg("failed send cache event info")
		}
	}(packageID, eventType, success, timestamp)
}

// Create implements mypackage.PackageService
func (ps *PackageService) Create(createMap map[string]string) (uint64, error) {

	// call general create func,
	// returns with
	// 1) database entry id
	// 2) model obj to set cache
	// 3) error
	id, unit, err := ps.create(createMap)
	if err != nil {
		return 0, err
	}

	if err := ps.cache.SetDescription(unit); err != nil {
		log.Debug().Err(err).Msg("failed set cache description")
		ps.asyncSendEvent(id, model.SetDescription, false)
	} else {
		ps.asyncSendEvent(id, model.SetDescription, true)
	}

	return id, nil
}

// General create function.
// Works throuht grpc api.
func (ps *PackageService) create(createMap map[string]string) (uint64, logistic.Package, error) {

	unit := logistic.Package{}
	req := pb.CreatePackageV1Request{}
	for key, value := range createMap {
		switch key {
		case logistic.Title:
			req.Title = value
			unit.Title = value

		case logistic.Material:
			req.Material = value
			unit.Material = value

		case logistic.MaximumVolume:
			volume, err := strconv.ParseFloat(value, 32)
			if err != nil {
				log.Debug().Err(err).Msg("failed to parse volume")
				return 0, logistic.Package{}, packageApi.ErrBadArgument{Argument: "MaximumVolume"}
			}
			req.MaximumVolume = float32(volume)
			unit.MaximumVolume = float32(volume)

		case logistic.Reusable:

			if strings.ToLower(value) == "да" {
				req.Reusable = true
				unit.Reusable = true
			} else if strings.ToLower(value) == "нет" {
				req.Reusable = false
				unit.Reusable = false
			} else {
				log.Debug().Msg("failed to parse reusable")
				return 0, logistic.Package{}, packageApi.ErrBadArgument{Argument: "Reusable"}
			}

		default:
			return 0, logistic.Package{}, packageApi.ErrBadArgument{Argument: key}
		}
	}

	if err := req.Validate(); err != nil {
		log.Debug().Err(err).Msg("PackageService.Create req validation failed")

		if err, ok := err.(pb.CreatePackageV1RequestValidationError); ok {
			return 0, logistic.Package{}, packageApi.ErrBadArgument{Argument: err.Field()}
		}

		return 0, logistic.Package{}, packageApi.ErrBadArgument{Argument: "unable to fetch invalid field"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := ps.grpcclient.CreatePackageV1(ctx, &req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.Internal:
				return 0, logistic.Package{}, packageApi.ErrInternal
			}
		}

		return 0, logistic.Package{}, err
	}

	unit.ID = resp.GetID()

	return resp.GetID(), unit, nil
}

// Describe implements mypackage.PackageService
func (ps *PackageService) Describe(packageID uint64) (logistic.Package, error) {

	if unit, err := ps.cache.ReadDescription(packageID); err != nil {
		log.Debug().Msg("cache miss")
		ps.asyncSendEvent(packageID, model.ReadDescription, false)
	} else {
		log.Debug().Msg("cache hit")
		ps.asyncSendEvent(packageID, model.ReadDescription, true)

		return *unit, nil
	}

	unit, err := ps.describe(packageID)
	if err != nil {
		return logistic.Package{}, err
	}

	if err := ps.cache.SetDescription(unit); err != nil {
		log.Debug().Err(err).Msg("failed set cache description")
		ps.asyncSendEvent(packageID, model.SetDescription, false)
	} else {
		ps.asyncSendEvent(packageID, model.SetDescription, true)
	}

	return unit, nil
}

// General describe function.
// Works throuht grpc api.
func (ps *PackageService) describe(packageID uint64) (logistic.Package, error) {
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

	pbPackage := resp.GetValue()
	unit := logistic.Package{
		ID:            packageID,
		Title:         pbPackage.GetTitle(),
		Material:      pbPackage.GetMaterial(),
		MaximumVolume: pbPackage.GetMaximumVolume(),
		Reusable:      pbPackage.GetReusable(),
	}

	return unit, nil
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

	if resp.GetSuc() {
		if err := ps.cache.RemoveDescription(packageID); err != nil {
			ps.asyncSendEvent(packageID, model.RemoveDescription, false)
		} else {
			ps.asyncSendEvent(packageID, model.RemoveDescription, true)
		}
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

	if resp.GetSuc() {
		if err := ps.cache.RemoveDescription(packageID); err != nil {
			ps.asyncSendEvent(packageID, model.RemoveDescription, false)
		} else {
			ps.asyncSendEvent(packageID, model.RemoveDescription, true)
		}
	}

	return resp.GetSuc(), nil
}

func NewService(cc grpc.ClientConnInterface, ch CacheDict, ces CacheEventSender) *PackageService {
	return &PackageService{
		grpcclient:       pb.NewLogisticPackageApiServiceClient(cc),
		cache:            ch,
		cacheEventSender: ces,
	}
}
