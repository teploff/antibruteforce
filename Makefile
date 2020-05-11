build:
	go build cmd/antibruteforce/main.go

build_cli:
	go build -o cli tools/admincli/cmd/admincli/main.go

run:
	go build cmd/antibruteforce/main.go && ./main

coverage:
	 go test ./... -v -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=./... -count=1

test:
	go test ./... -v -race ./... -count=1

lint:
	golangci-lint run --enable-all