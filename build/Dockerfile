
# Первый этап: сборка приложения
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp

# Второй этап: финальный образ
FROM alpine:latest
WORKDIR /app/

# Копируем всё содержимое из первого этапа
COPY --from=builder /app/ .

# Устанавливаем переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=./scheduler.db

# Указываем порт, на котором будет работать веб-сервер
EXPOSE ${TODO_PORT}

# Команда для запуска приложения
CMD ["./myapp"]