MIGRATE=migrate
DB_URL=postgresql://samualhalder:samualpass@localhost:5433/social?sslmode=disable
MIGRATION_DIR=cmd/migrate/migrations

up:
	$(MIGRATE) -path=$(MIGRATION_DIR) -database="$(DB_URL)" up

down:
	$(MIGRATE) -path=$(MIGRATION_DIR) -database="$(DB_URL)" down 1

force-clean:
	$(MIGRATE) -path=$(MIGRATION_DIR) -database="$(DB_URL)" force 1

version:
	$(MIGRATE) -path=$(MIGRATION_DIR) -database="$(DB_URL)" version

create:
	$(MIGRATE) create -seq -ext sql -dir $(MIGRATION_DIR) $(name)

seed:
	go run cmd/migrate/seed/main.go
