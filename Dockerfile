# Этап 1. Сборка приложения
FROM golang:1.26 AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение для Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# Этап 2. Создание образа для запуска
FROM ubuntu:latest

WORKDIR /app

# Копируем собранное приложение
COPY --from=builder /app/app .

# Копируем веб-интерфейс
COPY --from=builder /app/web ./web

# Открываем порт приложения
EXPOSE 7540

# Запускаем приложение
CMD ["./app"]