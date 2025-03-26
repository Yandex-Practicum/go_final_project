FROM golang:1.21-alpine as builder

ADD . /go/src/sheduler

WORKDIR /go/src/sheduler

RUN go build -o sheduler ./cmd/app



FROM alpine:latest

ENV TODO_PORT=7540
ENV TODO_DBFILE=./scheduler.db
ENV TODO_PASSWORD=test

EXPOSE ${TODO_PORT}

ADD ./web /opt/sheduler

COPY --from=builder /go/src/sheduler/sheduler /opt/sheduler

WORKDIR /opt/sheduler
VOLUME ${TODO_DBFILE}
CMD ["./sheduler"]
