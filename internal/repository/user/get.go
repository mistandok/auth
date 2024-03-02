package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	serviceModel "github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/internal/repository/user/convert"
	repoModel "github.com/mistandok/auth/internal/repository/user/model"
)

// Get user from db.
func (u *Repo) Get(ctx context.Context, userID serviceModel.UserID) (*serviceModel.User, error) {
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

	args := pgx.NamedArgs{
		idColumn: int64(userID),
	}

	rows, err := u.pool.Query(ctx, query, args)
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