-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS endpoint_access (
    id BIGSERIAL PRIMARY KEY,
    address TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role <> ''),
    created_at timestamp with time zone default current_timestamp NOT NULL,
    updated_at timestamp with time zone default current_timestamp NOT NULL
);

ALTER TABLE "endpoint_access" ADD CONSTRAINT idx_endpoint_access_address_role UNIQUE (address, role);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS endpoint_access;
-- +goose StatementEnd
