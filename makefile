.PHONY: build
build:
	go build -v ./cmd/main.go
.DEFAULT_GOAL := build

.PHONY: lint
lint:
	golangci-lint run -c ./.golangci.yml > lint.txt

.PHONY: format
format:
	go fmt ./...