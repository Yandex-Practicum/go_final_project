FROM golang:1.23.3
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 7540
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/task_app
CMD ["/app/task_app"]