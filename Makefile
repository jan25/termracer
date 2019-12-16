GOCMD=go
CMD_DIR=./cmd
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=termracer
BINARY_UNIX=$(BINARY_NAME)_unix

.PHONY: all
all: test build

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(CMD_DIR)

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(CMD_DIR)

# run all test in this project
# ./... looks for root, subdirectories in this project
.PHONY: test
test: 
	$(GOTEST) ./... -v

.PHONY: run
run: build
	./$(BINARY_NAME)

.PHONY: debug
debug: build
	./$(BINARY_NAME) -debug=true

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# docker-build:
# TODO
