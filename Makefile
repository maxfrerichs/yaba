BINARY_NAME=yaba
SRC=src/main.go

build:
	go build -o $(BINARY_NAME) $(SRC)

clean:
	go clean
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

install: build
	cp $(BINARY_NAME) /usr/local/bin/

uninstall:
	rm -f /usr/local/bin/$(BINARY_NAME)

help:
	@echo "Makefile commands:"
	@echo "  build  - Build the application"
	@echo "  clean  - Clean up build artifacts"
	@echo "  run    - Build and run the application"
	@echo "  help   - Display this help message"

.PHONY: build clean run install uninstall help
