include .env
MIGRATIONS_PATH = C:/Users/User/Desktop/Go-program/GolangProject/auth/cmd/migrate/migrations/

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=${MIGRATIONS_PATH} -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: gen-docs
gen-docs:
	@if not exist "docs" mkdir docs
	@swag init -d ./cmd/api -o ./docs --parseDependency