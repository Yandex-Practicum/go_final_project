
# Первый этап: сборка приложения
FROM golang:1.23-alpine AS builder

# Устанавливаем gcc
RUN apk add build-base 

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o myapp

# Второй этап: финальный образ
FROM alpine:latest
WORKDIR /app/

# Копируем всё содержимое из первого этапа
COPY --from=builder /app/ .

# Устанавливаем переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=./scheduler.db

# Указываем порт, на котором будет работать веб-сервер
EXPOSE 7540

# Команда для запуска приложения
CMD ["./myapp"]