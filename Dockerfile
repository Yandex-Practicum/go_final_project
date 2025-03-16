FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./my_app 

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/my_app /app/my_app
COPY . .

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=1234

EXPOSE ${TODO_PORT}
CMD ["/app/my_app"]
