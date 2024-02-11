postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123 -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root root

dropdb:
	docker exec -it postgres dropdb root

rundb:
	docker run postgres16

connect_bank:
	docker exec -it postgres psql -U root bank

connect_root:
	docker exec -it postgres psql -U root root

migrateup:
	migrate -path db/migration -database "postgresql://root:123@localhost:5432/root?sslmode=disable" -verbose up

aws_migrateup:
	migrate -path db/migration -database "postgresql://root:123456syd@simple-bankdb.cxa46u4qwwyq.us-east-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:123@localhost:5432/root?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:123@localhost:5432/root?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:123@localhost:5432/root?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v  -cover ./... 

server:
	go run main.go

proto:
	rm -f pb/*.go 
	rm -f doc/swagger/*.swagger.json
	export PATH="$PATH:$(go env GOPATH)/bin"
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto

.PHONY: proto

evans:
	evans --host localhost --port 9090  -r repl

.PHONY: evans

.PHONY: postgres createdb dropdb connect_bank connect_root 
	rundb migrateup migratedown sqlc test server migrateup1
	migratedown1 aws_migrateup 