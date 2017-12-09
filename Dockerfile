FROM astronomerio/alpine-librdkafka-golang:1.9-0.11.0-r0 as builder
WORKDIR /go/src/github.com/astronomerio/event-api
RUN apk update && apk add git
COPY . .
RUN make build

FROM alpine:3.4
COPY --from=builder /go/src/github.com/astronomerio/event-api/event-api /event-api
RUN apk update && apk add ca-certificates curl && update-ca-certificates
ENV GIN_MODE=release
EXPOSE 8080 8081
ENTRYPOINT ["/event-api"]
