FROM golang:1.21.0

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main cmd/api/main.go cmd/api/routes.go 

EXPOSE 7540

CMD ["/main"]