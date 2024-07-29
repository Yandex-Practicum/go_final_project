FROM golang:1.21.0

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /alexeyjudin ./cmd/main.go

CMD ["/alexeyjudin"]
