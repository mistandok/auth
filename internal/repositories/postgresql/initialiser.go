package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/auth/internal/config"
)

// MustInitPgConnection get gp connection pool or log fatal.
func MustInitPgConnection(ctx context.Context, cfg config.PgConfig) (*pgxpool.Pool, func()) {
	pgxConfig, err := pgxpool.ParseConfig(getPgURL(cfg))
	if err != nil {
		log.Fatalf("ошибка при формировании конфига для pgxpool: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatalf("ошибка при подключении к DB: %v", err)
	}
	closer := func() {
		pool.Close()
	}

	return pool, closer
}

func getPgURL(cfg config.PgConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)
}