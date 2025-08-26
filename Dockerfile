FROM golang:1.23
WORKDIR /app
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/mydatabase.db
COPY go.mod go.sum ./
RUN go mod download
COPY . .  
RUN apt-get update && apt-get install -y libsqlite3-dev
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /my_app ./cmd/service
EXPOSE 7540
CMD ["/my_app"]


