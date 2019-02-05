GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GODOWNLOAD = $(GOCMD) download

BIN_DIR = bin
SRC_DIR = src
BINARY_NAME = osmosis

.PHONY: all deps test clean

all: $(BIN_DIR)/$(BINARY_NAME)

$(BIN_DIR)/$(BINARY_NAME): main.go $(wildcard $(SRC_DIR)/**/*.go)
	$(GOBUILD) -v -o $@

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BIN_DIR)/$(BINARY_NAME)

deps:
	$(GODOWNLOAD)
