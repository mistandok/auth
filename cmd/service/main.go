package main

import (
	"context"
	"flag"
	userImpl "github.com/mistandok/auth/internal/api/user"
	postgresql2 "github.com/mistandok/auth/internal/repository/postgresql"
	"github.com/mistandok/auth/internal/repository/user"
	userService "github.com/mistandok/auth/internal/service/user"
	"log"
	"net"
	"os"

	"github.com/mistandok/auth/internal/config"
	"github.com/mistandok/auth/internal/config/env"
	"github.com/mistandok/auth/pkg/user_v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	pool, connCloser := postgresql2.MustInitPgConnection(ctx, *pgConfig)
	defer connCloser()

	listenConfig := net.ListenConfig{}
	listener, err := listenConfig.Listen(ctx, "tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("ошибка при прослушивании: %v", err)
	}

	logger := setupZeroLog(logConfig)

	userRepo := user.NewRepo(pool, logger)
	userServ := userService.NewService(logger, userRepo)
	userImplServer := userImpl.NewImplementation(userServ)

	server := grpc.NewServer()
	reflection.Register(server)
	user_v1.RegisterUserV1Server(server, userImplServer)

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
