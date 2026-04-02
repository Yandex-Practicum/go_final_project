FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o scheduler ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/scheduler .

COPY --from=builder /app/web ./web

EXPOSE 7540

CMD ["./scheduler"]