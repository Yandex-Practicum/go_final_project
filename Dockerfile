FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o todo-app .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/todo-app .

COPY web ./web

ENV TODO_PORT=7540

ENV TODO_DBFILE=/app/data/scheduler.db

EXPOSE 7540

CMD ["./todo-app"]