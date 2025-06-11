# Этап 1: сборка
FROM golang:1.23.4 AS builder

WORKDIR /app

# Установка git (нужен для go mod download)
RUN apt-get update && apt-get install -y git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка бинарного файла
RUN GOOS=linux GOARCH=amd64 go build -o scheduler .

# Этап 2: финальный минимальный образ
FROM ubuntu:latest

WORKDIR /app
RUN apt update && apt install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/scheduler .
COPY web ./web
COPY scheduler.db .

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db

CMD ["./scheduler"]
