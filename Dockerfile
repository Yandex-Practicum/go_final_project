FROM golang:alpine as builder
ADD . /go/src/sheduler
WORKDIR /go/src/sheduler
RUN apk add git && go get . && go build -o sheduler ./cmd/app


FROM alpine:latest

ENV TODO_PORT=7540
ENV TODO_DBFILE=./scheduler.db
ENV TODO_PASSWORD=test

EXPOSE ${TODO_PORT}

RUN mkdir /opt/sheduler
COPY --from=builder /go/src/sheduler/sheduler /opt/sheduler
ADD ./web /opt/sheduler

WORKDIR /opt/sheduler
VOLUME ${TODO_DBFILE}
CMD ["./vpn_bot"]
