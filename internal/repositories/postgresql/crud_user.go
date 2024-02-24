package postgresql

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/auth/internal/repositories"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"time"
)

type CRUDUser struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

func NewCRUDUser(pool *pgxpool.Pool, logger *zerolog.Logger) *CRUDUser {
	return &CRUDUser{
		pool:   pool,
		logger: logger,
	}
}

func (u *CRUDUser) Create(ctx context.Context, in *repositories.CRUDUserCreateIn) (*repositories.CRUDUserCreateOut, error) {
	query := `
	INSERT INTO "user" (name, email, password, role, created_at, updated_at)
	VALUES (@name, @email, @password, @role, @createdAt, @updatedAt)
	RETURNING id
	`

	currentTime := time.Now()

	args := pgx.NamedArgs{
		"name":      in.Name,
		"email":     in.Email,
		"role":      in.Role,
		"password":  in.Password,
		"createdAt": currentTime,
		"updatedAt": currentTime,
	}

	rows, err := u.pool.Query(ctx, query, args)
	if err != nil {
		return &repositories.CRUDUserCreateOut{}, errors.WithStack(err)
	}
	defer rows.Close()

	out, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repositories.CRUDUserCreateOut])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "user_email_key" {
			return nil, repositories.ErrEmailIsTaken
		}
		return nil, errors.WithStack(err)
	}

	return &out, nil
}

func (u *CRUDUser) Update(ctx context.Context, in *repositories.CRUDUserUpdateIn) error {
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
	params = append(params, in.Id)

	sql, _, err := query.ToSql()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = u.pool.Exec(ctx, sql, params...)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *CRUDUser) Get(ctx context.Context, in *repositories.CRUDUserGetIn) (*repositories.CRUDUserGetOut, error) {
	query := `
	SELECT 
	    id, name, email, role, created_at createdAt, updated_at updatedAt
	FROM 
	    "user"
	WHERE
	    id = @id
    `

	args := pgx.NamedArgs{
		"id": in.Id,
	}

	rows, err := u.pool.Query(ctx, query, args)
	if err != nil {
		return &repositories.CRUDUserGetOut{}, errors.WithStack(err)
	}
	defer rows.Close()

	out, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repositories.CRUDUserGetOut])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repositories.ErrUserNotFound
		}
		return nil, errors.WithStack(err)
	}

	return &out, nil
}

func (u *CRUDUser) Delete(ctx context.Context, in *repositories.CRUDUserDeleteIn) error {
	query := `
    	DELETE FROM "user"
		WHERE id = @id
    `

	args := pgx.NamedArgs{
		"id": in.Id,
	}

	_, err := u.pool.Exec(ctx, query, args)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
