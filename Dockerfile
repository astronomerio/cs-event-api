FROM astronomerio/alpine-librdkafka-golang:1.9-0.11.0-r0

WORKDIR /go/src/github.com/astronomerio/clickstream-ingestion-api

COPY . .

RUN go build -tags static -o server main.go

FROM alpine:3.4 

COPY --from=0 /go/src/github.com/astronomerio/clickstream-ingestion-api/server /server
COPY --from=0 /go/src/github.com/astronomerio/clickstream-ingestion-api/config.yml /config.yml

ENV GIN_MODE=release
EXPOSE 8080 8081

ENTRYPOINT ./server



