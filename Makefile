# запустить контейнер с базой данных:
#	docker compose up -d

LOCAL_BIN := $(CURDIR)/bin

bin-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/easyp-tech/easyp/cmd/easyp@v0.7.15
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.25.0

generate-proto:
	PATH=$(LOCAL_BIN):$$PATH easyp generate


migration-down:
	$(LOCAL_BIN)/goose $(opts) -dir ./migrations postgres "host=${} user=${} password =${} port=${} dbname=${}"

migration-up:
	$(LOCAL_BIN)/goose $(opts) -dir ./migrations postgres "host=${} user=${} password =${} port=${} dbname=${}"

migration:
	mkdir -p $(migrations_dir)
	$(LOCAL_BIN)/goose -dir ./migrations create $(shell bash -c 'read -p "Migration name: " migration_name; echo $$migration_name') sql