FROM golang:1.23.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o todo-app main.go

FROM alpine:3.19

WORKDIR /app

ARG TODO_PORT=7540
ENV TODO_PORT=${TODO_PORT}
ENV TODO_PASSWORD=""
ENV TODO_DBFILE=/data/scheduler.db

COPY --from=builder /app/todo-app ./todo-app
COPY web ./web

EXPOSE ${TODO_PORT}

CMD ["./todo-app"]
