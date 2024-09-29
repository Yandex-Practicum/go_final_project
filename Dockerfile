FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
sqlite3

WORKDIR /todo-list

EXPOSE 7540

COPY web ./web

COPY ./builds/todo_app ./cmd/

RUN mkdir database

#CMD ["cmd/todo_app"]