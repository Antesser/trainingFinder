include .env
export


LOCAL_BIN := $(CURDIR)/bin
PROTOC_VERSION := 31.1
PROTOC_ZIP := protoc-$(PROTOC_VERSION)-linux-x86_64.zip
migrations_dir:= $(CURDIR)/migrations

ifeq ($(PLATFORM),Darwin)
    PROTOC_ZIP = protoc-$(PROTOC_VERSION)-osx-x86_64.zip
endif
    PROTOC_URL := https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP)

bin-deps:
	curl -sSL $(PROTOC_URL) -o /tmp/$(PROTOC_ZIP)
	unzip -o /tmp/$(PROTOC_ZIP) -d /tmp/protoc
	chmod u+w /tmp/protoc/bin/protoc
	cp /tmp/protoc/bin/protoc $(LOCAL_BIN)/
	cp -r /tmp/protoc/include $(LOCAL_BIN)/include
	rm -rf /tmp/$(PROTOC_ZIP) /tmp/protoc
	GOBIN=$(LOCAL_BIN) go install github.com/easyp-tech/easyp/cmd/easyp@v0.7.15
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.25.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.2.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

generate-proto:
	PATH=$(LOCAL_BIN):$$PATH easyp generate

migration-up:
	$(LOCAL_BIN)/goose $(opts) -allow-missing -dir ./migrations postgres "host=$$DB_HOST port=$$DB_PORT user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=disable" up

migration-down:
	$(LOCAL_BIN)/goose $(opts) -dir ./migrations postgres "host=$$DB_HOST port=$$DB_PORT user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=disable" down

migration:
	mkdir -p $(migrations_dir)
	$(LOCAL_BIN)/goose -dir ./migrations create $(shell bash -c 'read -p "Migration name: " migration_name; echo $$migration_name') sql

vet:
	go vet ./...

start-db:
	docker compose up -d
run:
	go run cmd/server/main.go