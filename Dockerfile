FROM ubuntu:latest

WORKDIR /app

COPY scheduler .
COPY web ./web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

CMD ["/app/scheduler"]