FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo-app .
COPY --from=builder /app/web ./web

ENV TODO_PASSWORD=DokkaYandex
ENV TODO_DBFILE=/app/scheduler.db
ENV TODO_PORT=7540

EXPOSE 7540

CMD ["./todo-app"]
