-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NUll,
    password TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at timestamp with time zone default current_timestamp NOT NULL,
    updated_at timestamp with time zone default current_timestamp NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_role_email_name ON "user" (role, email, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
