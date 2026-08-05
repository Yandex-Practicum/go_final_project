FROM golang:1.26.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o my_app

FROM --platform=linux/amd64 alpine:latest
WORKDIR /app
COPY --from=builder /app/my_app .
COPY web/ web/
EXPOSE 7540

CMD ["./my_app"]

