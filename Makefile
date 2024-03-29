GOCMD = go
GOBUILD = $(GOCMD) build
GOINSTALL = $(GOCMD) install
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOMOD = $(GOCMD) mod
GOGET = $(GOCMD) get

BIN_DIR = bin
BINARY_NAME = osmosis

.PHONY: all install test clean get upgrade

all: $(BIN_DIR)/$(BINARY_NAME)

install:
	$(GOINSTALL)

$(BIN_DIR)/$(BINARY_NAME): main.go $(wildcard cmd/*.go) $(wildcard cmd/**/*.go)
	$(GOBUILD) -v -o $@

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	$(GOMOD) tidy
	rm -f $(BIN_DIR)/$(BINARY_NAME)

get:
	$(GOMOD) download

upgrade:
	$(GOGET) -u
