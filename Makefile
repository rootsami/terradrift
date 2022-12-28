# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME_CLI=terradrift-cli
BINARY_NAME_SERVER=terradrift-server

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME_CLI) -v -ldflags="-X main.version=${VERSION}" ./terradrift-cli
	$(GOBUILD) -o $(BINARY_NAME_SERVER) -v -ldflags="-X main.version=${VERSION}" ./terradrift-server

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f ./terradrift-cli/$(BINARY_NAME_CLI)
	rm -f ./terradrift-server/$(BINARY_NAME_SERVER)

run:
	$(GOBUILD) -o $(BINARY_NAME_CLI) -v ./cmd/terradrift-cli
	./$(BINARY_NAME_CLI)

deps:
	$(GOGET) github.com/...
