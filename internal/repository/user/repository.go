package user

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/auth/internal/repository"
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

var _ repository.UserRepository = (*Repo)(nil)

// Repo user repo for crud operation.
type Repo struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

// NewRepo  get new repo instance.
func NewRepo(pool *pgxpool.Pool, logger *zerolog.Logger) *Repo {
	return &Repo{
		pool:   pool,
		logger: logger,
	}
}
