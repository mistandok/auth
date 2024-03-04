package env

import (
	"errors"
	"os"

	"github.com/mistandok/auth/internal/config"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

// GRPCCfgSearcher searcher for grpc config.
type GRPCCfgSearcher struct{}

// NewGRPCCfgSearcher get instance for grpc config searcher.
func NewGRPCCfgSearcher() *GRPCCfgSearcher {
	return &GRPCCfgSearcher{}
}

// Get searcher for grpc config.
func (s *GRPCCfgSearcher) Get() (*config.GRPCConfig, error) {
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
