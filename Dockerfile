FROM ubuntu:latest AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y gcc libc-dev curl

RUN curl -fsSL https://go.dev/dl/go1.21.13.linux-amd64.tar.gz | tar -C /usr/local -xz

ENV PATH="/usr/local/go/bin:$PATH"

COPY . .

ENV CGO_ENABLED=1

RUN go build -o go_final_project .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/go_final_project .

COPY web ./web

COPY .env .

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=12345

EXPOSE 7540

RUN chmod +x go_final_project

CMD ["./go_final_project"]
