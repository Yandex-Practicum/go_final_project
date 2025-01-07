FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y \
    golang \
    sqlite3 \
    build-essential \
    curl

ENV TODO_PORT=7540
ENV TODO_DBFILE="../scheduler.db"
ENV TODO_PASSWORD="123123"

WORKDIR /app

COPY . /app

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./final-wev-server ./cmd/main.go

EXPOSE ${TODO_PORT}

CMD [ "./final-wev-server" ]