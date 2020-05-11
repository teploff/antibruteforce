build:
	go build cmd/antibruteforce/main.go

build_cli:
	go build -o cli tools/admincli/cmd/admincli/main.go

run:
	go build cmd/antibruteforce/main.go && ./main

docker_run:
	cd deployments/stage &&\
	docker-compose up -d --build &&\
	docker image prune -f

docker_stop:
	cd deployments/stage && \
    docker-compose down && \
    docker system prune --force --volumes

coverage:
	 go test ./... -v -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=./... -count=1

test:
	docker run -d -p 27017:27017 --name mongo_test mongo &&\
 	go test ./... -covermode=atomic -v -race ./... -count=1 &&\
 	docker rm -f mongo_test

lint:
	golangci-lint run --enable-all