# Build stage
FROM golang:1.24 AS build

RUN apt-get update && apt-get install -y gcc

WORKDIR /build

COPY . .

RUN go mod download

ENV CGO_ENABLED=1 GOOS=linux 

RUN go build -o taskmanager .

# Final image
FROM ubuntu:latest

ENV TODO_PORT=7540 TODO_DBFILE="./scheduler.db"

WORKDIR /app

COPY --from=build /build/taskmanager .
COPY ./web ./web

EXPOSE ${TODO_PORT}
ENTRYPOINT ["/app/taskmanager"]