postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123 -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root bank

dropdb:
	docker exec -it postgres dropdb bank

rundb:
	docker run postgres16

connect_bank:
	docker exec -it postgres psql -U root bank

connect_root:
	docker exec -it postgres psql -U root root

migrateup:
	migrate -path db/migration -database "postgresql://root:123@localhost:5432/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:123@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v  -cover ./... 

.PHONY: postgres createdb dropdb connect_bank connect_root 
	rundb migrateup migratedown sqlc test