GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=termracer
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v

# run all test in this project
# ./... looks for root, subdirectories in this project
test: 
	$(GOTEST) ./... -v

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# docker-build:
# TODO
