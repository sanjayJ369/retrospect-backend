DB_URL=postgresql://root:root@localhost:5432/retrospect?sslmode=disable

postgres: 
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root retrospect

dropdb: 
	docker exec -it postgres17 dropdb retrospect

migrateup: 
	migrate -path db/migration -database "$(DB_URL)" --verbose up

migrateup1: 
	migrate -path db/migration -database "$(DB_URL)" --verbose up 1

migratedown: 
	migrate -path db/migration -database "$(DB_URL)" --verbose down

migratedown1: 
	migrate -path db/migration -database "$(DB_URL)" --verbose down 1

dump_schema:
	@mkdir -p db/schema
	docker exec postgres17 pg_dump -s -U root retrospect > db/schema/schema.sql
	@echo "Database schema dumped to db/schema/schema.sql"

sqlc: 
	sqlc generate

db_docs: 
	dbdocs build docs/db.dbml

test:
	go test -v -cover -short ./...

server: 
	go run ./main.go

mock:
	mockgen -package mockDB -destination ./db/mock/store.go github.com/sanjayj369/retrospect-backend/db/sqlc Store
	mockgen -package mockmail -destination mail/mock/sender.go github.com/sanjayj369/retrospect-backend/mail EmailSender

.PHONY: postgres createdb dropdb migrateup migratedown dump_schema sqlc test server mock db_docs