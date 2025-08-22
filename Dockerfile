# Используем официальный образ Golang для компиляции приложения
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем go.mod и go.sum в контейнер
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Компилируем приложение
RUN go build -o go_final_project .

# Используем минимальный образ для запуска приложения
FROM ubuntu:latest

# Устанавливаем зависимости, необходимые для работы приложения
RUN apt-get update && apt-get install -y ca-certificates

# Копируем скомпилированный исполняемый файл и веб-файлы в контейнер
COPY --from=builder /app/go_final_project /usr/local/bin/go_final_project
COPY web /usr/local/web

# Устанавливаем переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=12345

# Открываем порт, на котором будет работать веб-сервер
EXPOSE 7540

# Запускаем приложение
CMD ["go_final_project"]
