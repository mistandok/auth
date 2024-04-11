package temp_redis

import "context"

type Client interface {
	DB() DB
	Close() error
}

// DB интерфейс для работы с БД
type DB interface {
	QueryExecutor
	Close() error
}

type QueryExecutor interface {
	DoContext(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error)
}
