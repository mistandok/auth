-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NUll,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(100) NOT NULL,
    created_at timestamp with time zone default current_timestamp NOT NULL
    updated_at timestamp with time zone default current_timestamp NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_role_email_name ON "user" (role, email, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
