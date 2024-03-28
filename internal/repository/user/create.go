package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mistandok/platform_common/pkg/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	serviceModel "github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository"
)

// Create user in db.
func (u *Repo) Create(ctx context.Context, in *serviceModel.UserForCreate) (int64, error) {
	queryFormat := `
	INSERT INTO "%s" (%s, %s, %s, %s, %s, %s)
	VALUES (@%s, @%s, @%s, @%s, @%s, @%s)
	RETURNING id
	`
	query := fmt.Sprintf(
		queryFormat,
		userTable, nameColumn, emailColumn, roleColumn, passwordColumn, createdAtColumn, updatedAtColumn,
		nameColumn, emailColumn, roleColumn, passwordColumn, createdAtColumn, updatedAtColumn,
	)

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	currentTime := time.Now()

	args := pgx.NamedArgs{
		nameColumn:      in.Name,
		emailColumn:     in.Email,
		roleColumn:      in.Role,
		passwordColumn:  in.Password,
		createdAtColumn: currentTime,
		updatedAtColumn: currentTime,
	}

	rows, err := u.db.DB().QueryContext(ctx, q, args)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var userID int64

	userID, err = pgx.CollectOneRow(rows, pgx.RowTo[int64])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == userEmailKeyConstraint {
			return 0, repository.ErrEmailIsTaken
		}
		return 0, err
	}

	return userID, nil
}
