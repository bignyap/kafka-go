include .env
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: test
test:
	@go test -v ./...

# .PHONY: migrate-create
# migration:
# 	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: goose-up
goose-up:
	@goose -dir $(MIGRATIONS_PATH) $(DB_DRIVER) "user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

.PHONY: producer-up
producer-up:
	@go run .\cmd\producer\main.go

.PHONY: consumer-up
consumer-up:
	@go run .\cmd\consumer\main.go

# .PHONY: seed
# seed: 
# 	@go run cmd/migrate/seed/main.go

# .PHONY: gen-docs
# gen-docs:
# 	@swag init -g ./api/main.go -d cmd,internal && swag fmt