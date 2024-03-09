package user

import (
	"context"
	"fmt"
	"github.com/mistandok/platform_common/pkg/db"

	"github.com/jackc/pgx/v5"
)

// Delete user from db.
func (u *Repo) Delete(ctx context.Context, userID int64) error {
	queryFormat := `
    	DELETE FROM "%s"
		WHERE %s = @%s
    `

	query := fmt.Sprintf(
		queryFormat,
		userTable,
		idColumn, idColumn,
	)

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	args := pgx.NamedArgs{
		idColumn: userID,
	}

	_, err := u.db.DB().ExecContext(ctx, q, args)
	if err != nil {
		return err
	}

	return nil
}
