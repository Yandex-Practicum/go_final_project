FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY .env .

RUN apk add --no-cache gcc musl-dev

RUN CGO_ENABLED=1 GOOS=linux go build -o go_final_project

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/go_final_project ./

COPY --from=build /app/web ./web

COPY --from=build /app/.env ./

EXPOSE 7540

CMD ["./go_final_project"]