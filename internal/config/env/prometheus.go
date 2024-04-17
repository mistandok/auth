package env

import (
	"errors"
	"os"

	"github.com/mistandok/auth/internal/config"
)

const (
	prometheusHostEnvName       = "PROMETHEUS_SERVER_HOST"
	prometheusPortEnvName       = "PROMETHEUS_SERVER_PORT"
	prometheusPublicPortEnvName = "PROMETHEUS_PORT"
)

// PrometheusCfgSearcher searcher for http config.
type PrometheusCfgSearcher struct{}

// NewPrometheusCfgSearcher get instance for http config searcher.
func NewPrometheusCfgSearcher() *PrometheusCfgSearcher {
	return &PrometheusCfgSearcher{}
}

// Get searcher for grpc config.
func (s *PrometheusCfgSearcher) Get() (*config.PrometheusConfig, error) {
	host := os.Getenv(prometheusHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("prometheus host not found")
	}

	port := os.Getenv(prometheusPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("prometheus port not found")
	}

	publicPort := os.Getenv(prometheusPublicPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("prometheus public port not found")
	}

	return &config.PrometheusConfig{
		Host:       host,
		Port:       port,
		PublicHost: host,
		PublicPort: publicPort,
	}, nil
}
