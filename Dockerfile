FROM astronomerio/alpine-librdkafka-golang:1.9-0.11.0-r0

WORKDIR /go/src/github.com/astronomerio/clickstream-ingestion-api

RUN apk update && apk add git

COPY . .

RUN make build

FROM alpine:3.4

COPY --from=0 /go/src/github.com/astronomerio/clickstream-ingestion-api/server /server
RUN apk update && apk add ca-certificates && update-ca-certificates

ENV GIN_MODE=release

EXPOSE 8080 8081

ENTRYPOINT ["/server"]

