package postgresql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/auth/internal/repositories"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	userTable              = "user"
	nameColumn             = "name"
	emailColumn            = "email"
	passwordColumn         = "password"
	roleColumn             = "role"
	createdAtColumn        = "created_at"
	updatedAtColumn        = "updated_at"
	userEmailKeyConstraint = "user_email_key"
	idColumn               = "id"
)

// UserRepo user repo for crud operation.
type UserRepo struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

// NewUserRepo  get new repo instance.
func NewUserRepo(pool *pgxpool.Pool, logger *zerolog.Logger) *UserRepo {
	return &UserRepo{
		pool:   pool,
		logger: logger,
	}
}

// Create user in db.
func (u *UserRepo) Create(ctx context.Context, in *repositories.UserCreateIn) (*repositories.UserCreateOut, error) {
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
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	out, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repositories.UserCreateOut])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == userEmailKeyConstraint {
			return nil, repositories.ErrEmailIsTaken
		}
		return nil, errors.WithStack(err)
	}

	return &out, nil
}

// Update user in db.
func (u *UserRepo) Update(ctx context.Context, in *repositories.UserUpdateIn) error {
	queryFormat := `
		UPDATE "%s"
		SET %s
		WHERE %s=@%s
    `

	setParams := make([]string, 0)
	setFormat := "%s=@%s"
	namedArgs := make(pgx.NamedArgs)

	if in.Name != nil {
		setParams = append(setParams, fmt.Sprintf(setFormat, nameColumn, nameColumn))
		namedArgs[nameColumn] = in.Name
	}

	if in.Role != nil {
		setParams = append(setParams, fmt.Sprintf(setFormat, roleColumn, roleColumn))
		namedArgs[roleColumn] = in.Role
	}

	if in.Email != nil {
		setParams = append(setParams, fmt.Sprintf(setFormat, emailColumn, emailColumn))
		namedArgs[emailColumn] = in.Email
	}

	isNeedUpdate := len(namedArgs) > 0

	if !isNeedUpdate {
		return nil
	}

	namedArgs[idColumn] = in.ID

	setParams = append(setParams, fmt.Sprintf(setFormat, updatedAtColumn, updatedAtColumn))
	namedArgs[updatedAtColumn] = time.Now()

	setParamsStr := strings.Join(setParams, ", ")

	query := fmt.Sprintf(queryFormat, userTable, setParamsStr, idColumn, idColumn)

	_, err := u.pool.Exec(ctx, query, namedArgs)

	return err
}

// Get user from db.
func (u *UserRepo) Get(ctx context.Context, in *repositories.UserGetIn) (*repositories.UserGetOut, error) {
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
		idColumn: in.ID,
	}

	rows, err := u.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repositories.UserGetOut])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repositories.ErrUserNotFound
		}
		return nil, err
	}

	return &out, nil
}

// Delete user from db.
func (u *UserRepo) Delete(ctx context.Context, in *repositories.UserDeleteIn) error {
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
		idColumn: in.ID,
	}

	_, err := u.pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}
