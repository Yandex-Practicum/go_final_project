FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

COPY *.go ./

COPY web/ ./web

COPY .env ./

RUN go mod tidy && go mod download

RUN go build -o /myapp

EXPOSE 7540

CMD [ "/myapp" ]

