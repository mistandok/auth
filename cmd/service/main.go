package main

import (
	"context"
	"flag"
	"log"

	"github.com/mistandok/auth/internal/app"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", ".env", "path to config file")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx, configPath)
	if err != nil {
		log.Fatalf("ошибка при инициализации приложения: %s", err.Error())
	}

	if err := application.Run(); err != nil {
		log.Fatalf("ошибка во время работы приложения: %s", err.Error())
	}
}
