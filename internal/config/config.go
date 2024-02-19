package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"net"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

type GRPCConfig struct {
	Host string
	Port string
}

func (cfg *GRPCConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}

type LogConfig struct {
	LogLevel   zerolog.Level
	TimeFormat string
}

type GRPCConfigSearcher interface {
	Get() (*GRPCConfig, error)
}

type LogConfigSearcher interface {
	Get() (*LogConfig, error)
}
