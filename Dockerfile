FROM golang:1.21.0

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /alexeyjudin ./cmd/main.go

CMD ["/alexeyjudin"]
