package grpcclient

import (
	"context"
	"time"

	"github.com/hablof/omp-bot/internal/config"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewConn(cfg *config.Config) (*grpc.ClientConn, error) {

	var (
		err        error
		connection *grpc.ClientConn
	)
	maxAttempts := cfg.GrpcAPI.Attempts

	for i := 0; i < maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.GrpcAPI.DialTimeout)*time.Second)
		defer cancel()

		connection, err = grpc.DialContext(
			ctx,
			cfg.GrpcAPI.Target,
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err == nil {
			break
		}

		log.Debug().Err(err).Msg("grpc-server dial attempt failed...")
		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, err
	}

	log.Info().Err(err).Msg("grpc-server dial succeeded")

	return connection, nil
}
