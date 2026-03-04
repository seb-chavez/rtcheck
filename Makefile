.PHONY: build test clean lint

build:
	go build -o rtcheck .

test:
	go test ./... -v

clean:
	rm -f rtcheck

lint:
	golangci-lint run ./...
