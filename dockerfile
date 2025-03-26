FROM golang:1.24.1-alpine AS builder

RUN apk update && apk add --no-cache \
    build-base \
    sqlite-dev

WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -v -o server main.go
FROM alpine:3.18

RUN apk add --no-cache \
    sqlite-libs \
    tzdata

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/web ./web
COPY --from=builder /app/.env .
EXPOSE 7540
VOLUME /data
CMD ["./server"]
