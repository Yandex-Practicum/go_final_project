# Базовый образ для сборки Go-приложения
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы в рабочую директорию
COPY . .

# Сборка приложения
RUN go build -o main .

# Базовый образ для запуска приложения
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исполняемый файл из предыдущего этапа
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

# Устанавливаем переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db
ENV TODO_PASSWORD=123

# Указываем порт для прослушивания; делаем его динамическим с помощью переменной окружения
EXPOSE ${TODO_PORT}

# Запускаем приложение
CMD ["./main"]