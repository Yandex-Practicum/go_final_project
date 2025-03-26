FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download

COPY . .
COPY db/scheduler.db /app/db/scheduler.db

RUN CGO_ENABLED=0 GOOS=linux go build -o go_final_project .

FROM ubuntu:latest

WORKDIR /root/

COPY --from=builder /app/go_final_project .
COPY --from=builder /app/.env .
COPY --from=builder /app/db/scheduler.db /root/db/scheduler.db

EXPOSE 7540

CMD ["./go_final_project"]