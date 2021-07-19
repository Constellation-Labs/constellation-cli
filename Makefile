GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=cl_cli

all: test build
build:
		echo "Building updater"
		$(GOBUILD) -a -v -o constellation-updater ./cmd/updater
		echo "Building cli tool"
		$(GOBUILD) -a -v -o constellation-cli ./cmd/cli
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)

# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o constellation-updater -v ./cmd/updater && $(GOBUILD) -o constellation-cli -v ./cmd/cli
