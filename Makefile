GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

BIN_DIR = bin
CMD_DIR = commands
BINARY_NAME = osmosis

.PHONY: all deps test clean

all: $(BIN_DIR)/$(BINARY_NAME)

$(BIN_DIR)/$(BINARY_NAME): main.go $(wildcard $(CMD_DIR)/*.go)
	$(GOBUILD) -v -o $@

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BIN_DIR)/$(BINARY_NAME)

get:
	$(GOGET)

upgrade:
	$(GOGET) -u
