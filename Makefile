# Makefile pour pdbc

BINARY_NAME=pdbc
INSTALL_PATH=/usr/local/bin

.PHONY: all build build-prod clean install

all: build

build:
	go build -o $(BINARY_NAME)

build-prod:
	go build -ldflags="-s -w -X main.version=$(shell git describe --tags --always --dirty)" -trimpath -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	go clean

install: build-prod
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

uninstall:
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
