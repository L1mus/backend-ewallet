include ./.env

MIGRATION_PATH=db/migrations
DATABASE_URL=postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

migrate-create:
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

migrate-up:
	@migrate -database $(DATABASE_URL) -path $(MIGRATION_PATH) up

migrate-down:
	@migrate -database $(DATABASE_URL) -path $(MIGRATION_PATH) down 1

migrate-down-all:
	@migrate -database $(DATABASE_URL) -path $(MIGRATION_PATH) down -all

migrate-force:
	@migrate -database $(DATABASE_URL) -path $(MIGRATION_PATH) force $(VERSION)

migrate-db-seed:
	@echo "Run seeding data..."
	go run db/seeds/seeder.go
