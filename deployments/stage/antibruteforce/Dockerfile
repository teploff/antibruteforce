FROM golang:1.14 as builder
LABEL mainater="Alexander Teplov teploff.aa@gmail.com"
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOSE=linux GO111MODULE=on go build -mod=vendor -a -installsuffix nocgo -o antibruteforce /app/cmd/antibruteforce/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/ ./
EXPOSE 8087 8088
CMD ["./antibruteforce", "--config=./init/config_stage.yaml"]