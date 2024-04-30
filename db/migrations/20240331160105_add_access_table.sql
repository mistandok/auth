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

INSERT INTO endpoint_access (address, role, created_at, updated_at)
VALUES
    ('/chat_v1.ChatV1/Create', 'USER', current_date, current_date),
    ('/chat_v1.ChatV1/SendMessage', 'USER', current_date, current_date),
    ('/chat_v1.ChatV1/Delete', 'USER', current_date, current_date),
    ('/chat_v1.ChatV1/ConnectChat', 'USER', current_date, current_date)
ON CONFLICT (address, role) DO UPDATE
SET role = EXCLUDED.role;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS endpoint_access;
-- +goose StatementEnd
