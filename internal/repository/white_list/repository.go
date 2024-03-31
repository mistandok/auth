package white_list

import (
	"context"
	"fmt"
	"github.com/mistandok/auth/internal/repository"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	setCommand = "SET"
	exCommand  = "EX"
	getCommand = "GET"
)

var _ repository.WhiteListRepository = (*WhiteListRepo)(nil)

// WhiteListRepo ..
type WhiteListRepo struct {
	pool *redis.Pool
}

// NewWhiteListRepo ..
func NewWhiteListRepo(pool *redis.Pool) *WhiteListRepo {
	return &WhiteListRepo{
		pool: pool,
	}
}

// Set записать токен в белый список
func (r *WhiteListRepo) Set(ctx context.Context, userID int64, jwtString string, expireIn time.Duration) error {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return fmt.Errorf("ошибка при получении соединения к WhiteListRepo из пула соединений: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	_, err = redis.String(conn.Do(setCommand, userID, jwtString, exCommand, expireIn.Seconds()))
	if err != nil {
		return fmt.Errorf("ошибка при попытке сохранить запись в WhiteListRepo: %w", err)
	}

	return err
}

// Get получить токен из белого списка.
func (r *WhiteListRepo) Get(ctx context.Context, userID int64) (string, error) {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении соединения к WhiteListRepo из пула соединений: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	reply, err := redis.String(conn.Do(getCommand, userID))
	if err != nil {
		return "", fmt.Errorf("ошибка при попытке поулчить запись из WhiteListRepo: %w", err)
	}

	return reply, nil
}
