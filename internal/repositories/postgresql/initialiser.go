package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/mistandok/auth/internal/config"
	"log"
)

func MustInitPgConnection(ctx context.Context, cfg config.PgConfig) (*pgx.Conn, func()) {
	conn, err := pgx.Connect(ctx, getPgUrl(cfg))
	if err != nil {
		log.Fatalf("ошибка при подключении к DB: %v", err)
	}
	closer := func() {
		conn.Close(ctx)
	}

	return conn, closer
}

func getPgUrl(cfg config.PgConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)
}
