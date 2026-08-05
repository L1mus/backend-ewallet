include ./.env

#postgres://user:password@host:port/dbname?query
MIGRATION_PATH=./schema/migrations
DB_URL=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL)

create_table:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

run-migrations:
	migrate -database $(DB_URL) -path $(MIGRATION_PATH) up

down-migration:
	migrate -database $(DB_URL) -path $(MIGRATION_PATH) down 1

down-all-migration:
	migrate -database $(DB_URL) -path $(MIGRATION_PATH) down -all

force-migration:
	migrate -database $(DB_URL) -path $(MIGRATION_PATH) force $(VERSION)