.PHONY: all build run test coverage start-docker

all: build test run

build:
	go build ./cmd/pismo

run:
	go run ./cmd/pismo

debug:
	dlv debug ./cmd/pismo

test:
	go install github.com/mfridman/tparse@latest | go mod tidy
	go test -json -cover ./... | tparse -all -pass

coverage:
	go test -count=1 -coverprofile cover.out ./... -covermode=count
	go tool cover -html=cover.out

start-docker:
	docker compose up -d

