# 1. Этап сборки
FROM golang:1.23-bullseye AS builder

WORKDIR /todo-list-project

RUN apt-get update && apt-get install -y gcc libc6-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка бинарного файла в директории проекта
RUN CC=gcc CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /todo-list-project/todo-list ./cmd

# 2. Финальный минимальный образ
FROM debian:bullseye

WORKDIR /app

# Копируем бинарный файл
COPY --from=builder /todo-list-project/todo-list .

# Копируем конфиг
COPY --from=builder /todo-list-project/config ./config

# Копируем .env (если он существует в контексте билда)
COPY .env .env

# Запуск приложения
CMD ["/app/todo-list"]
