FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

COPY *.go ./

COPY web/ ./web

COPY .env ./

RUN go mod tidy && go mod download

RUN go build -o /grigorpower

EXPOSE 7540

CMD [ "/grigorpower" ]

