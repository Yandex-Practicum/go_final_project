## Build state
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/scheduler ./cmd/main.go


## Final state
FROM alpine:3.21

WORKDIR /shced

COPY --from=builder /app/bin/scheduler .

EXPOSE 7540

CMD ["./scheduler"]
