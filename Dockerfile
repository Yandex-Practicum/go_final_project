FROM debian:latest

WORKDIR /app

COPY ./web ./web
COPY final .
COPY scheduler.db .

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db
ENV TODO_PASSWORD=852369_qeT

EXPOSE $TODO_PORT

CMD [ "./final" ]

