FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
sqlite3

WORKDIR /todo-list

EXPOSE 7540

COPY web ./web

COPY ./builds/todo_app ./cmd/

RUN apt install -y curl

CMD ["cmd/todo_app"]