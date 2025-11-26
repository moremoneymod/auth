LOCAL_BIN:=$(CURDIR)/bin

include .env
export

install-deps:
	set GOBIN=$(LOCAL_BIN) && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	set GOBIN=$(LOCAL_BIN) && go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	set GOBIN=$(LOCAL_BIN) && go install github.com/pressly/goose/v3/cmd/goose@latest

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


local-migration-status:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

local-migration-up:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v



generate-auth-api:
	if not exist pkg\auth_v1 mkdir pkg\auth_v1
	protoc --proto_path api/auth_v1 --proto_path="C:\protoc-33.0-win\include" --go_out=pkg/auth_v1 --go_opt=paths=source_relative --plugin=protoc-gen-go=bin/protoc-gen-go.exe --go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative --plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc.exe api/auth_v1/auth.proto


generate-access-api:
	if not exist pkg\access_v1 mkdir pkg\access_v1
	protoc --proto_path api/access_v1 --proto_path="C:\protoc-33.0-win\include" --go_out=pkg/access_v1 --go_opt=paths=source_relative --plugin=protoc-gen-go=bin/protoc-gen-go.exe --go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative --plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc.exe api/access_v1/access.proto


ghz-login-load-test:
	ghz \
		--proto api/auth_v1/auth.proto \
		--call auth_v1.AuthV1.Login \
		--data '{"username":"test", "password":"qwerty123"}' \
		--rps 1000 \
		--total 30000 \
		--insecure \
		localhost:50051




generate: generate-auth-api generate-access-api

.PHONY: install-deps get-deps generate generate-auth-api