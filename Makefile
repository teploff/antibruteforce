build:
	go build cmd/antibruteforce/main.go

build_cli:
	go build -o cli tools/admincli/cmd/admincli/main.go

run:
	go build cmd/antibruteforce/main.go && ./main

coverage:
	go test ./... -v -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=github.com/teploff/antibruteforce/domain/entity,github.com/teploff/antibruteforce/internal/implementation/repository/bucket,github.com/teploff/antibruteforce/internal/implementation/repository/ip

test:
	go test ./... -v -race ./...

lint:
	golangci-lint run --enable-all