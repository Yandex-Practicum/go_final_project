build:
	go build

tidy:
	go mod tidy

fmt:
	gofumpt -w .
	gci write . --skip-generated -s standard -s default

lint: tidy fmt build
	golangci-lint run