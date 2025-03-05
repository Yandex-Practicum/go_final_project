FROM golang:1.23.2-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./

COPY web/  ./bin/web
COPY sqlite/scheduler_creator.sql ./bin/sqlite

#Переменные окружения для ОС
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

#Переменная окружения для адреса базы данных
ENV TODO_DBFILE=""

#Создаёт исполняемый файл в среде, реализованной в контейнере
RUN go build -o /app/bin/todorun ./cmd/todo/main.go

#Запускает исполняемый файл
CMD ["./bin/todorun"]
