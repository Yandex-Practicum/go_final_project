FROM golang:1.23.2 as init-stage

FROM ubuntu:latest

WORKDIR /app

COPY . .

FROM init-stage AS build-stage

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go mod download
RUN go build -o /todoserver

ARG TODOPORT=7540
ENV TODO_PORT=${TODOPORT}
ARG TODOPASSWORD=1234
ENV TODO_PASSWORD=${TODOPASSWORD}
ARG TODODBFILE=scheduler.db
ENV TODO_DBFILE=${TODODBFILE}

EXPOSE ${TODO_PORT}
CMD ["./todoapp"]