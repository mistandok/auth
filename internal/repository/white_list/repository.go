package white_list

import (
	"context"
	"fmt"
	"github.com/mistandok/auth/internal/temp_redis"
	"time"

	"github.com/mistandok/auth/internal/repository"
)

const (
	setCommand = "SET"
	exCommand  = "EX"
	getCommand = "GET"
)

var _ repository.WhiteListRepository = (*WhiteListRepo)(nil)

// WhiteListRepo ..
type WhiteListRepo struct {
	client temp_redis.Client
}

// NewWhiteListRepo ..
func NewWhiteListRepo(client temp_redis.Client) *WhiteListRepo {
	return &WhiteListRepo{
		client: client,
	}
}

// Set записать токен в белый список
func (r *WhiteListRepo) Set(ctx context.Context, userID int64, jwtString string, expireIn time.Duration) error {
	db := r.client.DB()
	_, err := db.String(db.DoContext(ctx, setCommand, userID, jwtString, exCommand, expireIn.Seconds()))
	if err != nil {
		return fmt.Errorf("ошибка при попытке сохранить запись в WhiteListRepo: %w", err)
	}

	return err
}

// Get получить токен из белого списка.
func (r *WhiteListRepo) Get(ctx context.Context, userID int64) (string, error) {
	db := r.client.DB()
	reply, err := db.String(db.DoContext(ctx, getCommand, userID))
	if err != nil {
		return "", fmt.Errorf("ошибка при попытке поулчить запись из WhiteListRepo: %w", err)
	}

	return reply, nil
}
