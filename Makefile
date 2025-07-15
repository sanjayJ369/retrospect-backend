postgres: 
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root retrospect

dropdb: 
	docker exec -it postgres17 dropdb retrospect

migrateup: 
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/retrospect?sslmode=disable" --verbose up

migratedown: 
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/retrospect?sslmode=disable" --verbose down

dump_schema:
	@mkdir -p db/schema
	docker exec postgres17 pg_dump -s -U root retrospect > db/schema/schema.sql
	@echo "Database schema dumped to db/schema/schema.sql"

sqlc: 
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown dump_schema sqlc