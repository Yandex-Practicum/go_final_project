# Этап сборки
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Собираем исполняемый файл
RUN go build -o todo-go .

# Этап финального образа
FROM ubuntu:latest

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исполняемый файл из этапа сборки
COPY --from=builder /app/todo-go /app/todo-go

# Копируем поддиректорию web
COPY web /app/web

# Открываем порт для веб-сервера
EXPOSE 7540

# Команда для запуска приложения
CMD ["/app/todo-go"]
