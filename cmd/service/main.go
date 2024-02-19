package main

import (
	"fmt"
	"github.com/mistandok/auth/internal/user/server_v1"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/mistandok/auth/pkg/user_v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const envName = ".env.local"

func main() {
	envs := mustGetEnvs(envName)
	grpcPort := mustFetchGRPCPort(envs)
	logLevel := mustFetchLogLevel(envs)
	listener := mustGetListener(grpcPort)

	logger := setupZeroLog(logLevel, time.RFC3339)
	userServer := server_v1.NewServer(logger)

	server := grpc.NewServer()
	reflection.Register(server)
	user_v1.RegisterUserV1Server(server, userServer)

	log.Printf("сервер запущен на %v", listener.Addr())

	if err := server.Serve(listener); err != nil {
		log.Fatalf("ошибка сервера: %v", err)
	}
}

func mustGetEnvs(env string) map[string]string {
	envs, err := godotenv.Read(env)
	if err != nil {
		log.Fatalf("не удалось прочитать переменные окружения: %v", err)
	}

	return envs
}

func mustFetchGRPCPort(envs map[string]string) int {
	name := "GRPC_PORT"
	portStr, ok := envs[name]
	if !ok {
		log.Fatalf("не задана переменная окружения: %s", name)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("некорректное значение для переменной окружения %s, err: %v", name, err)
	}

	return port
}

func mustGetListener(port int) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("ошибка при прослушивании: %v", err)
	}

	return lis
}

func mustFetchLogLevel(envs map[string]string) zerolog.Level {
	name := "LOG_LEVEL"
	logLevelStr, ok := envs[name]
	if !ok {
		log.Fatalf("не задана переменная окружения: %s", name)
	}
	logLevelInt, err := strconv.Atoi(logLevelStr)
	if err != nil {
		log.Fatalf("некорректное значение для переменной окружения %s, err: %v", name, err)
	}
	return zerolog.Level(logLevelInt)
}

func setupZeroLog(logLevel zerolog.Level, timeFormat string) *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat}
	logger := zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = timeFormat

	return &logger
}
