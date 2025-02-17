FROM golang:1.23 AS builder
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN apt-get update && apt-get install -y \
    gcc \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*
COPY . .
RUN go build -o my_app ./cmd/server/main.go

FROM debian:latest
WORKDIR /app
COPY --from=builder app/my_app .
COPY ./web /app/web
EXPOSE $TODO_PORT
CMD ["/app/my_app"]