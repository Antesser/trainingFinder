# запустить контейнер с базой данных:
#	docker compose up -d

LOCAL_BIN := $(CURDIR)/bin

bin-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/easyp-tech/easyp/cmd/easyp@v0.7.15

generate-proto:
	PATH=$(LOCAL_BIN):$$PATH easyp generate