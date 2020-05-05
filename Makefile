build:
	go build cmd/antibruteforce/main.go

run:
	go build cmd/antibruteforce/main.go && ./main

test:
	go test ./... -v -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=./...