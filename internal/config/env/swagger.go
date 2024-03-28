package env

import (
	"errors"
	"os"

	"github.com/mistandok/auth/internal/config"
)

const (
	swaggerHostEnvName = "SWAGGER_HOST"
	swaggerPortEnvName = "SWAGGER_PORT"
)

// SwaggerConfigSearcher searcher for http config.
type SwaggerConfigSearcher struct{}

// NewSwaggerConfigSearcher get instance for http config searcher.
func NewSwaggerConfigSearcher() *SwaggerConfigSearcher {
	return &SwaggerConfigSearcher{}
}

// Get searcher for grpc config.
func (s *SwaggerConfigSearcher) Get() (*config.SwaggerConfig, error) {
	host := os.Getenv(swaggerHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http host not found")
	}

	port := os.Getenv(swaggerPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}

	return &config.SwaggerConfig{
		Host: host,
		Port: port,
	}, nil
}
