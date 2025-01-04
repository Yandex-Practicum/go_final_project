FROM golang:1.21-bullseye AS builder

WORKDIR /app
COPY . .

RUN apt-get update && apt-get install -y gcc g++ libc6-dev build-essential sqlite3 libsqlite3-dev
RUN go mod tidy && CGO_ENABLED=1 go build -o main ./cmd/main.go

# Переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=1234
ENV JWT_SECRET=z6ewwwYIU4Tc25wQTxtGHiT3e6UZzgEV8tputDNJouY=

# Открываем порт
EXPOSE 7540

CMD ["./main"]
