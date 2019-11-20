# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
SOURCE_CLIENT_NAME=client/client
BINARY_CLIENT_NAME=client/client.go
SOURCE_SERVER_NAME=server/server
BINARY_SERVER_NAME=server/server.go

all: deps build test
deps:
	$(GOGET) github.com/spf13/cobra
install:
	$(GOBUILD) -o $(SOURCE_CLIENT_NAME) $(BINARY_CLIENT_NAME)
	$(GOBUILD) -o $(SOURCE_SERVER_NAME) $(BINARY_SERVER_NAME)
clean:
	$(GOCLEAN)
	rm -f $(BINARY_CLIENT_NAME)
	$(GOCLEAN)
	rm -f $(BINARY_SERVER_NAME)
run_server:
	$(GOBUILD) -o $(SOURCE_SERVER_NAME) $(BINARY_SERVER_NAME)
	./$(SOURCE_SERVER_NAME)
run_client:
	$(GOBUILD) -o $(SOURCE_CLIENT_NAME) $(BINARY_CLIENT_NAME)
	./$(SOURCE_CLIENT_NAME) go-telnet --timeout=21s 127.0.0.1 3303
