FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o todo-app main.go

FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /app/todo-app ./todo-app
COPY web ./web
EXPOSE 7540
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=""
CMD ["./todo-app"]
