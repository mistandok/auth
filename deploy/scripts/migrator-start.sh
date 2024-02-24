#!/bin/sh
wait_database()
{
  HOST=$1
  PORT=$2

  echo "Waiting for database..."
  echo "$HOST:$PORT"
  while ! nc -z $HOST $PORT; do
    sleep 1
  done

  echo "database started"
}

export MIGRATION_DSN="host=${PG_HOST} port=${PG_PORT} dbname=${POSTGRES_DB} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} sslmode=disable"

wait_database $PG_HOST $PG_PORT
goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v