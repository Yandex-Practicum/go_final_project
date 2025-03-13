FROM golang:1.23.1-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache gcc musl-dev sqlite

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app/server cmd/server/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates sqlite

COPY .env /app/.env

COPY --from=builder /app/server /app/server

COPY --from=builder /app/web ./web/

RUN chmod +x /app/server

EXPOSE 7540

CMD ["/app/server"]