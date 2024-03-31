package endpoint_access

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	serviceModel "github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/platform_common/pkg/db"
	"time"
)

func (r *Repo) Create(ctx context.Context, endpointAccess *serviceModel.EndpointAccess) (int64, error) {
	queryFormat := `
	INSERT INTO "%s" (%s, %s, %s, %s)
	VALUES (@%s, @%s, @%s, @%s)
	RETURNING id
	`

	query := fmt.Sprintf(
		queryFormat,
		endpointAccessTable, addressColumn, roleColumn, createdAtColumn, updatedAtColumn,
		addressColumn, roleColumn, createdAtColumn, updatedAtColumn,
	)

	q := db.Query{
		Name:     "endpoint_access_repository.Create",
		QueryRaw: query,
	}

	currentTime := time.Now()

	args := pgx.NamedArgs{
		addressColumn:   endpointAccess.Address,
		roleColumn:      endpointAccess.Role,
		createdAtColumn: currentTime,
		updatedAtColumn: currentTime,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var endpointAccessID int64

	endpointAccessID, err = pgx.CollectOneRow(rows, pgx.RowTo[int64])
	if err != nil {
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.ConstraintName == endpointAccessAddressKeyConstraint {
				return 0, repository.ErrEndpointAccessExists
			}
			return 0, err
		}
	}

	return endpointAccessID, nil
}
