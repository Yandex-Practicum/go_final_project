FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app

FROM ubuntu:latest

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db
ENV TODO_PASSWORD=123

WORKDIR /app
COPY --from=builder /my_app /app/my_app

EXPOSE ${TODO_PORT}

CMD ["/app/my_app"]