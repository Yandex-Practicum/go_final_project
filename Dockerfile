FROM golang:1.22.1

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV TODO_PORT=7540

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE ${TODO_PORT}

RUN  go build -o task_manager_server

CMD ["/task_manager_server"]