-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
    id BIGSERIAL PRIMARY KEY,
    name text NOT NULL,
    email text UNIQUE NOT NUll,
    password text NOT NULL,
    role text NOT NULL,
    created_at timestamp with time zone default current_timestamp NOT NULL,
    updated_at timestamp with time zone default current_timestamp NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_role_email_name ON "user" (role, email, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
