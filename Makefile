GO_BUILD = go build
GOFLAGS  = CGO_ENABLED=0
DATABASE_HOST     ?= localhost
DATABASE_PORT     ?= $(shell grep "DB_PORT" .env | cut -d '=' -f2)
DATABASE_NAME 	  ?= $(shell grep "DB_NAME" .env | cut -d '=' -f2)
DATABASE_USERNAME ?= $(shell grep "DB_USERNAME" .env | cut -d '=' -f2)
DATABASE_PASSWORD ?= $(shell grep "DB_PASSWORD" .env | cut -d '=' -f2)
DATABSE_DSN       ?= ${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}

# Version of migrations - this is optionally used on goto command
V?=

# Number of migrations - this is optionally used on up and down commands
N?=

.PHONY: migrate_setup migrate_up migrate_down migrate_goto migrate_drop_db
migrate_setup:
	@if [ -z "$$(which migrate)" ]; then echo "Installing golang-migrate..."; go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; fi

migrate_up: migrate_setup
	@ migrate -database 'postgres://${DATABSE_DSN}?sslmode=disable' -path $$(pwd)/migrations up $(N)

migrate_down: migrate_setup
	@ migrate -database 'postgres://${DATABSE_DSN}?sslmode=disable' -path $$(pwd)/migrations down $(N)

migrate_goto: migrate_setup
	@ migrate -database 'postgres://${DATABSE_DSN}?sslmode=disable' -path $$(pwd)/migrations goto $(V)

migrate_drop_db: migrate_setup
	@ migrate -database 'postgres://${DATABSE_DSN}?sslmode=disable' -path $$(pwd)/migrations drop

.PHONY: swagger
swagger:
	@swag init -g util/swagger.go -o api

## build: Build app binary
.PHONY: build
build:
	$(GOFLAGS) $(GO_BUILD) -a -v -ldflags="-w -s" -o bin/app main.go
