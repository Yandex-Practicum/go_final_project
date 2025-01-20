# ubuntu:latest is the latest version of ubuntu #
FROM ubuntu:latest 

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV TODO_PORT=7540

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y golang
RUN apt-get update && apt-get install -y ca-certificates

RUN go mod download

EXPOSE ${TODO_PORT}

RUN  go build -o server_tasks

CMD ["./server_tasks"]