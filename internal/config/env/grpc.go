package env

import (
	"github.com/mistandok/auth/internal/config"
	"os"

	"github.com/pkg/errors"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type GRPCConfigSearcher struct{}

func NewGRPCConfigSearcher() *GRPCConfigSearcher {
	return &GRPCConfigSearcher{}
}

func (s *GRPCConfigSearcher) Get() (*config.GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	return &config.GRPCConfig{
		Host: host,
		Port: port,
	}, nil
}
