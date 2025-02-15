FROM golang:1.23-alpine AS builder


WORKDIR /app


COPY . .


RUN go mod download


RUN go build -o scheduler .


FROM alpine:latest


WORKDIR /app


COPY --from=builder /app/scheduler .


COPY ./web ./web


ENV TODO_PORT=7540


EXPOSE 7540


CMD ["./scheduler"]