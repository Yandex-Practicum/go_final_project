FROM golang:1.22.2

WORKDIR /usr/src/app

COPY . .

RUN go mod download 

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /usr/src/app/cmd/todo/main ./cmd/todo/main.go

ENV TODO_PORT=7540 \
    TODO_DBFILE="/usr/src/app/scheduler.db" \
    TODO_WEB_DIR="/usr/src/app/web"

CMD ["/usr/src/app/cmd/todo/main"] 