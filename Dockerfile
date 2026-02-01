# Многостадийная сборка для уменьшения размера образа
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/app ./cmd/app/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/bin/app /app/app

# Копируем все файлы проекта, включая .env
COPY --from=builder /app/web /app/web
COPY --from=builder /app/pkg /app/pkg
COPY --from=builder /app/cmd /app/cmd
COPY --from=builder /app/go.mod /app/go.mod
COPY --from=builder /app/go.sum /app/go.sum

# Копируем .env файл (если он есть в проекте)
COPY --from=builder /app/.env /app/.env

# Устанавливаем права на выполнение
RUN chmod +x /app/app

# Порт по умолчанию
EXPOSE 7540

# Запускаем приложение
CMD ["/app/app"]