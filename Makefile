DB_URL=postgresql://booketo_admin:virtualrelay@localhost:5432/booketo_db?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

mysql:
	docker run --name mysql8 -p 3306:3306  -e MYSQL_ROOT_PASSWORD=secret -d mysql:8

createdb:
	docker exec -it postgres12 createdb --username=booketo_admin --owner=booketo_admin booketo_db

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path src/db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path src/db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path src/db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path src/db/migration -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	rm -f src/doc/*.sql
	dbml2sql --postgres -o src/doc/schema.sql src/doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store


.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 db_docs db_schema sqlc test server mock proto evans redis
