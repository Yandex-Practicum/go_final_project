# Сборка
FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /todo

# Финальный образ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /todo .
COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=""

VOLUME /data
EXPOSE ${TODO_PORT}

CMD ["/app/todo"]