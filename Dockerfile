# Базовый образ с Go
FROM golang:latest AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы проекта внутрь контейнера
COPY . .

# Загружаем зависимости
RUN go mod tidy

# Собираем бинарник (указываем путь к main.go!)
RUN go build -o /app/todo_app ./cmd/todo_app

# Финальный минимальный образ
FROM ubuntu:latest
WORKDIR /app

# Копируем скомпилированный бинарник
COPY --from=builder /app/todo_app /app/todo_app

# Указываем порт
EXPOSE 7540

# Запускаем сервер
CMD ["/app/todo_app"]
