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
	"github.com/mistandok/platform_common/pkg/db"
)

// GetByFilter user from db.
func (u *Repo) GetByFilter(ctx context.Context, filter *serviceModel.UserFilter) (*serviceModel.User, error) {
	args := pgx.NamedArgs{}
	whereColumn := ""

	if filter.ID != nil {
		args[idColumn] = *filter.ID
		whereColumn = idColumn
	} else if filter.Email != nil {
		args[emailColumn] = *filter.Email
		whereColumn = emailColumn
	}

	if len(args) == 0 {
		return nil, repository.ErrIncorrectFilters
	}

	queryFormat := `
	SELECT 
	    %s, %s, %s, %s, %s, %s %s, %s %s
	FROM 
	    "%s"
	WHERE
	    %s = @%s
    `

	query := fmt.Sprintf(
		queryFormat,
		idColumn, nameColumn, roleColumn, emailColumn, passwordColumn, createdAtColumn, createdAtAliasColumn, updatedAtColumn, updatedAtAliasColumn,
		userTable,
		whereColumn, whereColumn,
	)

	q := db.Query{
		Name:     "user_repository.GetByFilter",
		QueryRaw: query,
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
