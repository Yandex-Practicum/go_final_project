FROM golang:1.22
WORKDIR /go_final_project
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app
COPY scheduler.db ./
CMD ["/my_app"] 