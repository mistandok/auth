package white_list

import (
	"context"
	"fmt"
	"time"

	"github.com/mistandok/platform_common/pkg/memory_db"

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
	client memory_db.Client
}

// NewWhiteListRepo ..
func NewWhiteListRepo(client memory_db.Client) *WhiteListRepo {
	return &WhiteListRepo{
		client: client,
	}
}

// Set записать токен в белый список
func (r *WhiteListRepo) Set(ctx context.Context, userID int64, jwtString string, expireIn time.Duration) error {
	_, err := r.client.DB().DoContext(ctx, setCommand, userID, jwtString, exCommand, expireIn.Seconds())
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
