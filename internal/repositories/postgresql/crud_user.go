package postgresql

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/auth/internal/repositories"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// CRUDUserRepo user repo for crud operation.
type CRUDUserRepo struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

// NewCRUDUserRepo  get new repo instance.
func NewCRUDUserRepo(pool *pgxpool.Pool, logger *zerolog.Logger) *CRUDUserRepo {
	return &CRUDUserRepo{
		pool:   pool,
		logger: logger,
	}
}

// Create user in db.
func (u *CRUDUserRepo) Create(ctx context.Context, in *repositories.CRUDUserCreateIn) (*repositories.CRUDUserCreateOut, error) {
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
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	out, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repositories.CRUDUserCreateOut])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.ConstraintName == "user_email_key":
				return nil, repositories.ErrEmailIsTaken
			case pgErr.ConstraintName == "user_name_key":
				return nil, repositories.ErrNameIsTaken
			}
		}
		return nil, errors.WithStack(err)
	}

	return &out, nil
}

// Update user in db.
func (u *CRUDUserRepo) Update(ctx context.Context, in *repositories.CRUDUserUpdateIn) error {
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
func (u *CRUDUserRepo) Get(ctx context.Context, in *repositories.CRUDUserGetIn) (*repositories.CRUDUserGetOut, error) {
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

// Delete user from db.
func (u *CRUDUserRepo) Delete(ctx context.Context, in *repositories.CRUDUserDeleteIn) error {
	query := `
    	DELETE FROM "user"
		WHERE id = @id
    `

	args := pgx.NamedArgs{
		"id": in.ID,
	}

	_, err := u.pool.Exec(ctx, query, args)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
