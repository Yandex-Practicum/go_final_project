# Используем образ golang для компиляции
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Компилируем приложение для Linux
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o scheduler .

# Создаём финальный образ на основе ubuntu
FROM ubuntu:22.04

WORKDIR /app

# Устанавливаем зависимости для SQLite
RUN apt-get update && apt-get install -y sqlite3

# Копируем скомпилированный бинарник и директорию web
COPY --from=builder /app/scheduler .
COPY web ./web

# Указываем порт
EXPOSE 7540

# Команда для запуска приложения
CMD ["./scheduler"]