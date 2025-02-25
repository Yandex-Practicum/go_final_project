FROM golang:1.22.2 AS builder

WORKDIR /app
COPY . .

RUN go build -o app

FROM ubuntu:latest

WORKDIR /app
COPY --from=builder /app/app ./
COPY --from=builder /app/web ./web

EXPOSE 7540

CMD ["/app/app"]
