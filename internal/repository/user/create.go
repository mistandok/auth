package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	serviceModel "github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository"
	"github.com/pkg/errors"
)

// Create user in db.
func (u *Repo) Create(ctx context.Context, in *serviceModel.UserForCreate) (serviceModel.UserID, error) {
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

	currentTime := time.Now()

	args := pgx.NamedArgs{
		nameColumn:      in.Name,
		emailColumn:     in.Email,
		roleColumn:      in.Role,
		passwordColumn:  in.Password,
		createdAtColumn: currentTime,
		updatedAtColumn: currentTime,
	}

	rows, err := u.pool.Query(ctx, query, args)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var userID serviceModel.UserID

	userID, err = pgx.CollectOneRow(rows, pgx.RowTo[serviceModel.UserID])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == userEmailKeyConstraint {
			return 0, repository.ErrEmailIsTaken
		}
		return 0, err
	}

	return userID, nil
}
