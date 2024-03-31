package endpoint_access

import (
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/platform_common/pkg/db"
	"github.com/rs/zerolog"
)

const (
	endpointAccessTable                = "endpoint_access"
	idColumn                           = "id"
	addressColumn                      = "address"
	roleColumn                         = "role"
	createdAtColumn                    = "created_at"
	updatedAtColumn                    = "updated_at"
	createdAtAliasColumn               = "createdAt"
	updatedAtAliasColumn               = "updatedAt"
	endpointAccessAddressKeyConstraint = "idx_endpoint_access_address_role"
)

var _ repository.EndpointAccessRepository = (*Repo)(nil)

type Repo struct {
	logger *zerolog.Logger
	db     db.Client
}

// NewRepo  get new repo instance.
func NewRepo(logger *zerolog.Logger, client db.Client) *Repo {
	return &Repo{
		logger: logger,
		db:     client,
	}
}
