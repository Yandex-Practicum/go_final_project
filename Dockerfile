# Стадия сборки: builder
FROM golang:1.23.5-alpine AS builder

# Устанавливаем необходимые пакеты: git, gcc, musl-dev, sqlite
RUN apk add --no-cache git gcc musl-dev sqlite

WORKDIR /app

# Копируем файлы зависимостей и скачиваем модули
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник; так как main.go находится в каталоге cmd, переходим туда
RUN cd cmd && CGO_ENABLED=1 GOOS=linux go build -o /app/go_final_project main.go

# Финальный образ для приложения (используется сервис app)
FROM alpine:latest
RUN apk add --no-cache sqlite

WORKDIR /app
COPY --from=builder /app/go_final_project .

EXPOSE 7540

ENV SERVER_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

CMD ["/app/go_final_project"]
