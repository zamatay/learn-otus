BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
MIGRATION_PATH="./migrations"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

run: build
	$(BIN) -config ./configs/config.yaml

dc-build:
	/usr/bin/docker compose -f ./deployments/docker-compose.yaml -p deployments build

up:
	/usr/bin/docker compose -f ./deployments/docker-compose.yaml -p deployments up -d

down:
	/usr/bin/docker compose -f ./deployments/docker-compose.yaml -p deployments down

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

migration:
	goose -dir $(MIGRATION_PATH) postgres "host=0.0.0.0 user=postgres password=postgres dbname=calendar sslmode=disable" up

generate:
	@echo "Генерация скриптов protobuf."
	protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=. --grpc-gateway_opt generate_unbound_methods=true --openapiv2_out . proto/api.proto
container_up:
	@echo "Запуск контейнеров"
	docker-compose -f deployments/docker-compose.yaml up
container_down:
	@echo "Остановке контейнеров"
	docker-compose -f deployments/docker-compose.yaml down
cmt:
	git commit --allow-empty -m"Trigger pipelines"

integration-test:
	set -e ;\
    test_status_code=0 ;\
    docker-compose -f ./deployments/docker-compose-test.yaml up --build -d || test_status_code=$$? ;\
    docker-compose -f ./deployments/docker-compose-test.yaml down ;\
    exit $$test_status_code ;\