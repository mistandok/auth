package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/mistandok/auth/internal/client/db"

	"github.com/jackc/pgx/v5"
	serviceModel "github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/internal/repository/user/convert"
	repoModel "github.com/mistandok/auth/internal/repository/user/model"
)

// Get user from db.
func (u *Repo) Get(ctx context.Context, userID int64) (*serviceModel.User, error) {
	queryFormat := `
	SELECT 
	    %s, %s, %s, %s, %s createdAt, %s updatedAt
	FROM 
	    "%s"
	WHERE
	    %s = @%s
    `

	query := fmt.Sprintf(
		queryFormat,
		idColumn, nameColumn, roleColumn, emailColumn, createdAtColumn, updatedAtColumn,
		userTable,
		idColumn, idColumn,
	)

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	args := pgx.NamedArgs{
		idColumn: userID,
	}

	rows, err := u.db.DB().QueryContext(ctx, q, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repoModel.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return convert.ToServiceUserFromRepoUser(&out), nil
}
