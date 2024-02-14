package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mistandok/auth/internal/servers/user"
	"github.com/mistandok/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

const envName = ".env"

func main() {
	envs := mustGetEnvs(envName)
	grpcPort := mustFetchGRPCPort(envs)
	listener := mustGetListener(grpcPort)

	server := grpc.NewServer()
	reflection.Register(server)
	user_v1.RegisterUserV1Server(server, &user.Server{})

	log.Printf("сервер слушает на %v", listener.Addr())

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
