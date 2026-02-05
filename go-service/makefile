DB_URL=postgresql://postgres:postgres@localhost:5432/notedbv2?sslmode=disable
MIGRATIONS_PATH=./migrations

.PHONY: migrate-up migrate-down migrate-force migrate-version migrate-create

migrate-up:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) up

migrate-down:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) down 1

migrate-force:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) force $(version)

migrate-version:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) version

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
