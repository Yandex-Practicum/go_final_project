all: build
build:
	@go build -o main cmd/api/main.go cmd/api/routes.go
run: 
	@go run cmd/api/main.go cmd/api/routes.go

docker-run:

docker-down: