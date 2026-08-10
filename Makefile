.PHONY: build run clean

BINARY_NAME=getUserInventory

build:
	@go build -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

clean:
	@rm -f $(BINARY_NAME)