.PHONY: default test
all: default test

gosec:
	go get github.com/securego/gosec/cmd/gosec
sec:
	@gosec ./...
	@echo "[OK] Go security check was completed!"

init:
	export GOPROXY=https://goproxy.cn

lint:
	go mod tidy
	gofumports -w .
	gofumpt -w .
	gofmt -s -w .
	go mod tidy
	go fmt ./...
	revive .
	goimports -w .
	golangci-lint run --enable-all

install: init
	go install -ldflags="-s -w" ./...

test: init
	go test ./...

bench: init
	go test -bench . ./...
