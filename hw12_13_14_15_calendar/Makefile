BIN := "./bin/calendar"
BIN_SCH := "./bin/calendar_scheduler"
BIN_SND := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCH) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SND) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/test.toml &
	$(BIN_SCH) -config ./configs/test_scheduler.toml &
	$(BIN_SND) -config ./configs/test_sender.toml

build-img:
	docker-compose -f ./deployments/docker-compose.yaml build --build-arg LDFLAGS="${LDFLAGS}"

run-img: build-img
	docker-compose -f ./deployments/docker-compose.yaml up

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.8

lint: install-lint-deps
	golangci-lint run ./...

generate:
	rm -rf ./pkg
	mkdir -p ./pkg
	protoc --go_out=./pkg --go-grpc_out=./pkg/ --proto_path=./protobuf ./protobuf/*.proto

.PHONY: build run build-img run-img version test lint
