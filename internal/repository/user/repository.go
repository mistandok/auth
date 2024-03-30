package user

import (
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/platform_common/pkg/db"
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
	createdAtAliasColumn   = "createdAt"
	updatedAtAliasColumn   = "updatedAt"
	userEmailKeyConstraint = "user_email_key"
	idColumn               = "id"
)

var _ repository.UserRepository = (*Repo)(nil)

// Repo user repo for crud operation.
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
