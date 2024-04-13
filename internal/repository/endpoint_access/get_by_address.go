package endpoint_access

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	serviceModel "github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/internal/repository/endpoint_access/convert"
	repoModel "github.com/mistandok/auth/internal/repository/endpoint_access/model"
	"github.com/mistandok/platform_common/pkg/db"
)

// GetByAddressAndRole получить настройку доступа по адресу и роли
func (r *Repo) GetByAddressAndRole(ctx context.Context, address string, role serviceModel.UserRole) (*serviceModel.EndpointAccess, error) {
	queryFormat := `
	SELECT 
	    %s, %s, %s, %s %s, %s %s
	FROM 
	    %s
	WHERE
	    %s = @%s
		AND %s = @%s
    `

	query := fmt.Sprintf(
		queryFormat,
		idColumn, addressColumn, roleColumn, createdAtColumn, createdAtAliasColumn, updatedAtColumn, updatedAtAliasColumn,
		endpointAccessTable,
		addressColumn, addressColumn,
		roleColumn, roleColumn,
	)

	q := db.Query{
		Name:     "endpoint_access_repository.GetByAddress",
		QueryRaw: query,
	}

	args := pgx.NamedArgs{
		addressColumn: address,
		roleColumn:    role,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repoModel.EndpointAccess])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrEndpointAccessNotFound
		}
		return nil, err
	}

	return convert.ToServiceEndpointAccessFromRepoEndpointAccess(&out), nil
}
