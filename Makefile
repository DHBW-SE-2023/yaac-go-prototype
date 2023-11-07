BINARY_PATH=build/yaac
SOURCE_PATH=cmd/yaac/

all: build test

build: $(SOURCE_PATH)/*.go
	go build -o ./$(BINARY_PATH) ./$(SOURCE_PATH)

test: $(SOURCE_PATH)/*.go
	go test -v ./test/

run:
	make build
	./$(BINARY_PATH)

clean:
	go clean
	rm ./$(BINARY_PATH)

.PHONY: all, build, test, run, clean