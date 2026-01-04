# Makefile pour pdbc

BINARY_NAME=pdbc
INSTALL_PATH=/usr/local/bin

.PHONY: all build run clean install

all: build

build:
	go build -o $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	go clean

install: build
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

uninstall:
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
