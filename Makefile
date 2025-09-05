dev:
	@air --build.cmd="make build" --build.bin="bin/taski"

build:
	@go build -o bin/taski cmd/api/main.go

new-migration:
	@migrate create -ext sql -dir internal/core/database/postgres/migrations -seq $(name)

migrate-up:
	@migrate -path internal/core/database/postgres/migrations -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable up

migrate-down:
	@migrate -path internal/core/database/postgres/migrations -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable down
