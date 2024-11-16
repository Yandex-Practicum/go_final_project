FROM ubuntu:latest

WORKDIR /app

COPY ./web ./web
COPY final_project ./


RUN apt-get update -y
RUN apt-get install sqlite3 -y
RUN apt-get install golang -y


ENV TODO_PORT=7540
ENV TODO_DBFILE=./app/scheduler.db
ENV TODO_PASSWORD=852369_qeT

EXPOSE 7540

CMD [ "./final_project" ]

