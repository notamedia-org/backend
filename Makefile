.PHONY: all build test lint

all: build

build:
	go build \
	-mod=vendor \
	-o bin/beats \
	-ldflags "-X main.version=$(VERSION)" \
	cmd/beats/main.go

test:
	go test -v ./...

bench:
	go test -bench=. -benchmem ./...

lint: linter
	bin/golangci-lint run

linter: bin/golangci-lint

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.54.2

