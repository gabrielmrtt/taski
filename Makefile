dev:
	@air --build.cmd="make build" --build.bin="bin/taski"

build:
	@go build -o bin/taski cmd/api/main.go

new-migration:
	@migrate create -ext sql -dir internal/core/database/migrations -seq $(name)

migrate-up:
ifeq ($(env),test)
	@echo "Upgrading migrations for test"
	@go run cmd/migrate/main.go test up
else
	@echo "Upgrading migrations for default"
	@go run cmd/migrate/main.go default up
endif

migrate-down:
ifeq ($(env),test)
	@echo "Downgrading migrations for test"
	@go run cmd/migrate/main.go test down $(step)
else
	@echo "Downgrading migrations for default"
	@go run cmd/migrate/main.go default down $(step)
endif

seed:
ifeq ($(env),test)
	@echo "Seeding for test"
	@go run cmd/seed/main.go test
else
	@echo "Seeding for default"
	@go run cmd/seed/main.go
endif

swagger:
	@swag init -g cmd/api/main.go -d . -o docs --parseDependency --parseInternal