# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINDIR=bins
BINARY_NAME_CLI=terradrift-cli
BINARY_NAME_SERVER=terradrift-server

all: test build
build:
	$(GOBUILD) -o $(BINDIR)/$(BINARY_NAME_CLI) -v -ldflags="-X main.version=${VERSION}" ./terradrift-cli
	$(GOBUILD) -o $(BINDIR)/$(BINARY_NAME_SERVER) -v -ldflags="-X main.version=${VERSION}" ./terradrift-server

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINDIR)/

deps:
	$(GOCMD) mod download
