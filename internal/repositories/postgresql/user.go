package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
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
	query := `
	INSERT INTO "%s" (%s, %s, %s, %s, %s, %s)
	VALUES (@%s, @%s, @%s, @%s, @%s, @%s)
	RETURNING id
	`
	query = fmt.Sprintf(
		query,
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
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	params := make([]any, 0)

	query := psql.Update("\"user\"")

	if in.Name != nil {
		query = query.Set("name", "?")
		params = append(params, *in.Name)
	}

	if in.Role != nil {
		query = query.Set("role", "?")
		params = append(params, *in.Role)
	}

	if in.Email != nil {
		query = query.Set("email", "?")
		params = append(params, *in.Email)
	}

	query = query.Set("updated_at", "?")
	params = append(params, time.Now())

	query = query.Where(squirrel.Eq{"id": "?"})
	params = append(params, in.ID)

	query = query.Suffix("RETURNING \"id\"")

	sql, _, err := query.ToSql()
	if err != nil {
		return errors.WithStack(err)
	}

	row := u.pool.QueryRow(ctx, sql, params...)
	var userID int64
	if err = row.Scan(&userID); errors.Is(err, pgx.ErrNoRows) {
		return repositories.ErrUserNotFound
	}

	return err
}

// Get user from db.
func (u *UserRepo) Get(ctx context.Context, in *repositories.UserGetIn) (*repositories.UserGetOut, error) {
	query := `
	SELECT 
	    id, name, email, role, created_at createdAt, updated_at updatedAt
	FROM 
	    "user"
	WHERE
	    id = @id
    `

	args := pgx.NamedArgs{
		"id": in.ID,
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
	query := `
    	DELETE FROM "user"
		WHERE id = @id
    `

	args := pgx.NamedArgs{
		"id": in.ID,
	}

	_, err := u.pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}
