lint:
	go fmt ./...
	go vet ./...
	golint ./...

build:
	go build -o ./cmd/kids2 ./cmd/

run:
	NODE_CONFIG_PATH=./cmd/nodes.properties ./cmd/kids2

staticcheck:
	staticcheck ./...

all: lint staticcheck build run
