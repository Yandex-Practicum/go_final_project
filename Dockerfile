#docker build -t todo-app .
#docker run -d --name todo-app-container -p 7540:7540 todo-app

# Multi-Stage Build [image size 23 MB]
# Stage 1: Build stage
FROM golang:alpine AS builder
# Set the working directory
WORKDIR /app
# Download & install necessary the gcc, libc-dev and binutils packages
RUN apk add --no-cache build-base
# Copy and download dependencies
COPY go.* ./
RUN go mod download
# Copy the source code
COPY . .
# Build the Go application
RUN go build -o  todo-app ./cmd/main.go

# Stage 2: Final stage
FROM alpine:latest
# Set the working directory
WORKDIR /app
# Download & install necessary Sqlite libs packages
RUN apk add --no-cache libc6-compat sqlite-libs
# Copy the binary from the build stage
COPY --from=builder /app/todo-app ./
COPY --from=builder /app/web ./web
COPY --from=builder /app/internal/config/local.yaml ./internal/config/local.yaml

# Copy the database
COPY scheduler.db .
# Port binding
EXPOSE 7540
# Set the entrypoint command
CMD ["./todo-app", "-config=./internal/config/local.yaml"]

# Uncomment to use One-Stage Build
# # Basic One-Stage [image size 1.7 GB]
# FROM golang:1.22.5
# # Set the working directory
# WORKDIR /app
# # Copy and download dependencies
# COPY go.* ./ 
# RUN go mod download
# # Copy the source code
# COPY . .
# # Build the Go application
# RUN go build -o todo-app ./cmd/main.go 
# #Port binding
# EXPOSE 7540
# # Set the entrypoint command
# CMD ["./todo-app", "-config=./internal/config/local.yaml"]




