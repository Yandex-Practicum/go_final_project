# Используем многостадийную сборку
FROM golang:1.25.1-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache gcc musl-dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod go.sum ./
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Компилируем приложение
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

# Устанавливаем зависимости для SQLite
RUN apk --no-cache add ca-certificates

# Создаем пользователя для безопасности
RUN addgroup -S app && adduser -S app -G app

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарник из builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

# Создаем директорию для базы данных
RUN mkdir -p /data && chown -R app:app /data

# Переключаемся на пользователя app
USER app

# Используем переменную для порта 
EXPOSE ${TODO_PORT:-7540}

# Переменные окружения по умолчанию
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

# Запускаем приложение
CMD ["./main"]