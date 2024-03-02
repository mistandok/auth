package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/auth/internal/client/db"
	"github.com/mistandok/auth/internal/client/db/prettier"
	"github.com/rs/zerolog"
)

type key string

const (
	TxKey key = "tx" // TxKey ключ, по которому из контекста можно достать транзакцию
)

type pg struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

// NewDB ..
func NewDB(dbc *pgxpool.Pool, logger *zerolog.Logger) db.DB {
	return &pg{
		pool:   dbc,
		logger: logger,
	}
}

// ExecContext ..
func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, p.logger, q)

	tx, ok := ContextTx(ctx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}

	return p.pool.Exec(ctx, q.QueryRaw, args...)
}

// QueryContext ..
func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	logQuery(ctx, p.logger, q)

	tx, ok := ContextTx(ctx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	return p.pool.Query(ctx, q.QueryRaw, args...)
}

// QueryRowContext ..
func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, p.logger, q)

	tx, ok := ContextTx(ctx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.pool.QueryRow(ctx, q.QueryRaw, args...)
}

// CopyFromContext ..
func (p *pg) CopyFromContext(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	tx, ok := ContextTx(ctx)
	if ok {
		return tx.CopyFrom(
			ctx,
			tableName,
			columnNames,
			rowSrc,
		)
	}

	return p.pool.CopyFrom(
		ctx,
		tableName,
		columnNames,
		rowSrc,
	)
}

// BeginTx ..
func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.pool.BeginTx(ctx, txOptions)
}

// Ping ..
func (p *pg) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

// Close ..
func (p *pg) Close() {
	p.pool.Close()
}

// MakeContextTx добавляет транзакцию в контекст
func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

// ContextTx извлекает транзакцию из контекста
func ContextTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		return nil, false
	}

	return tx, true
}

func logQuery(_ context.Context, logger *zerolog.Logger, q db.Query) {
	if logLevel := logger.GetLevel(); logLevel == zerolog.DebugLevel {
		prettyQuery := prettier.Pretty(q.QueryRaw)
		logger.Debug().Str("sql", q.Name).Str("query", prettyQuery).Msg("лог запроса")
	}
}
