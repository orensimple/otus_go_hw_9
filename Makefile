# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
SOURCE_CLIENT_NAME=client.go
BINARY_CLIENT_NAME=./client/client.go
SOURCE_SERVER_NAME=server.go
BINARY_SERVER_NAME=./serever/server.go

all: deps build test
deps:
	$(GOGET) github.com/spf13/cobra
build:
	$(GOBUILD) -o $(BINARY_CLIENT_NAME) -v
	$(GOBUILD) -o $(BINARY_SERVER_NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_CLIENT_NAME)
	$(GOCLEAN)
	rm -f $(BINARY_SERVER_NAME)
run:
	$(GOBUILD) -o $(SOURCE_CLIENT_NAME) -v ./...
	./$(BINARY_CLIENT_NAME) 
	$(GOBUILD) -o $(SOURCE_SERVER_NAME) -v ./...
	./$(BINARY_SERVER_NAME) go-telnet --timeout=21s 127.0.0.1 3303
