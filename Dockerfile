FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    golang-go \
    git \
    build-essential && \
    apt-get clean

WORKDIR /app

COPY . .

RUN go build -o todo-app main.go

EXPOSE 7540

ENV TODO_PORT=7540 \
    TODO_DBFILE=/app/scheduler.db \
    TODO_PASSWORD=secret

CMD ["./todo-app"]