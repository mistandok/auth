package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	serviceModel "github.com/mistandok/auth/internal/model"
)

// Delete user from db.
func (u *Repo) Delete(ctx context.Context, userID serviceModel.UserID) error {
	queryFormat := `
    	DELETE FROM "%s"
		WHERE %s = @%s
    `

	query := fmt.Sprintf(
		queryFormat,
		userTable,
		idColumn, idColumn,
	)

	args := pgx.NamedArgs{
		idColumn: int64(userID),
	}

	_, err := u.pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}
