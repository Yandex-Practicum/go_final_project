FROM golang:latest

ENV TODO_PORT 7540
ENV TODO_PASSWORD 12345
ENV TODO_DBFILE ./scheduler.db

WORKDIR /app_go

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY app ./app
COPY web ./web

EXPOSE 7540
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /my_app

CMD ["/my_app"]