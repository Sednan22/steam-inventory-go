.PHONY: all build run clean

all: build run

build:
	@go build -o getUserInventory

run: 
	@./getUserInventory

clean:
	@rm -f getUserInventory