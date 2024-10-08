FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Базовый образ для запуска приложения
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db
ENV TODO_PASSWORD=123

EXPOSE ${TODO_PORT}

CMD ["./main"]