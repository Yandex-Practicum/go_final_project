FROM golang:1.22 AS builder
WORKDIR /build

COPY cmd /build/cmd
COPY config /build/config
COPY internal /build/internal
COPY infrastructure /build/infrastructure
COPY web /build/web
COPY go.mod go.sum /build/

RUN go mod download

RUN GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app ./cmd/main.go

FROM alpine:3.16 as server

COPY --from=builder /build/app /app/
COPY --from=builder /build/web /app/web

WORKDIR /app

CMD ["./app"]
