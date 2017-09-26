FROM golang:1.9-alpine

WORKDIR /go/src/github.com/astronomerio/clickstream-ingestion-api

RUN apk update && apk add git make

COPY . .

RUN go build -o server main.go

FROM alpine:3.4

COPY --from=0 /go/src/github.com/astronomerio/clickstream-ingestion-api/server /server
COPY --from=0 /go/src/github.com/astronomerio/clickstream-ingestion-api/config.yml /config.yml

ENV GIN_MODE=release
EXPOSE 8080 8081

ENTRYPOINT ./server



