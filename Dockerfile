FROM golang:1.22.3 as builder

WORKDIR /app

COPY . .

RUN go build -o main .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_PASSWORD=your_password

CMD ["./main"]
