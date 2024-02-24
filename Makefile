define setup_env
	$(eval ENV_FILE := ./deploy/env/.env.$(1))
	@echo "- setup env $(ENV_FILE)"
	$(eval include ./deploy/env/.env.$(1))
	$(eval export)
endef

setup-local-env:
	$(call setup_env,local)

setup-prod-env:
	$(call setup_env,prod)

LOCAL_BIN:=$(CURDIR)/bin

MIGRATION_DIR=$(MIGRATION_DIR)
MIGRATION_DSN="host=$(PG_HOST) port=$(PG_PORT) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@v0.18.0
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

lint:
	GOBIN=$(LOCAL_BIN) $(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

fix-imports:
	GOBIN=$(LOCAL_BIN) $(LOCAL_BIN)/goimports -w .

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-user-api

generate-user-api:
	mkdir -p pkg/user_v1
	protoc --proto_path api/user_v1 \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/user_v1/user.proto

migration-status:
	GOBIN=$(LOCAL_BIN) $(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} status -v

migration-up:
	GOBIN=$(LOCAL_BIN) $(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} up -v

migration-down:
	GOBIN=$(LOCAL_BIN) $(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} down -v

local-migration-status: setup-local-env migration-status

local-migration-up: setup-local-env migration-up

local-migration-down: setup-local-env migration-down

prod-migration-status: setup-prod-env migration-status

prod-migration-up: setup-prod-env migration-up

prod-migration-down: setup-prod-env migration-down

local-down-app:
	docker-compose --env-file deploy/env/.env.local -f deploy/docker-compose.local.yaml down -v

local-start-app:
	docker-compose --env-file deploy/env/.env.local -f deploy/docker-compose.local.yaml up -d --build

prod-down-app:
	docker-compose --env-file deploy/env/.env.prod -f deploy/docker-compose.prod.yaml down -v

prod-start-app:
	docker-compose --env-file deploy/env/.env.prod -f deploy/docker-compose.prod.yaml up -d --build

create-new-migration:
	GOBIN=$(LOCAL_BIN) $(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} create $(migration_name) sql