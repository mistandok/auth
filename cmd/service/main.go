package main

import (
	"context"
	"flag"
	"github.com/mistandok/auth/internal/config"
	"github.com/mistandok/auth/internal/config/env"
	"github.com/mistandok/auth/internal/repositories/postgresql"
	"github.com/mistandok/auth/internal/user/server_v1"
	"github.com/mistandok/auth/pkg/user_v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", ".env", "path to config file")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("ошибка при получении переменных окружения: %v", err)
	}

	grpcConfig, err := env.NewGRPCCfgSearcher().Get()
	if err != nil {
		log.Fatalf("ошибка при получении конфига grpc: %v", err)
	}

	logConfig, err := env.NewLogCfgSearcher().Get()
	if err != nil {
		log.Fatalf("ошибка при получении уровня логирования: %v", err)
	}

	pgConfig, err := env.NewPgCfgSearcher().Get()
	if err != nil {
		log.Fatalf("ошибка при получении конфига DB: %v", err)
	}

	_, connCloser := postgresql.MustInitPgConnection(ctx, *pgConfig)
	defer connCloser()

	listenConfig := net.ListenConfig{}
	listener, err := listenConfig.Listen(ctx, "tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("ошибка при прослушивании: %v", err)
	}

	logger := setupZeroLog(logConfig)
	userServer := server_v1.NewServer(logger)

	server := grpc.NewServer()
	reflection.Register(server)
	user_v1.RegisterUserV1Server(server, userServer)

	log.Printf("сервер запущен на %v", listener.Addr())

	if err := server.Serve(listener); err != nil {
		log.Fatalf("ошибка сервера: %v", err)
	}
}

func setupZeroLog(logConfig *config.LogConfig) *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: logConfig.TimeFormat}
	logger := zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(logConfig.LogLevel)
	zerolog.TimeFieldFormat = logConfig.TimeFormat

	return &logger
}
