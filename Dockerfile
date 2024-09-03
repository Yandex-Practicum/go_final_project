FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum.

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /todo_app main.go

EXPOSE 7540

CMD ["/todo_app"] 