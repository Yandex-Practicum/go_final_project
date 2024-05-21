# Компиляция через образ golang
FROM golang:1.22 AS GO
WORKDIR /src/todo-app
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o todo-app-exec
# Копирование в образ ubuntu исполняемого файла, web директории, .env файла и файлов для тестирования
FROM ubuntu:latest
WORKDIR /app/todo-app
COPY --from=GO /src/todo-app/.env /src/todo-app/todo-app-exec /src/todo-app/go.mod /src/todo-app/go.sum ./
COPY --from=GO /src/todo-app/tests ./tests
COPY --from=GO /src/todo-app/web ./web
EXPOSE 7540
CMD ["./todo-app-exec"]