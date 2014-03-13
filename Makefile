## simple makefile to log workflow
.PHONY: all test clean build

all: build test
	@echo "## bye."

build:
	@echo "build github.com/sbinet/go-clang"
	@go get -v ./...

test: build
	@go test

## EOF
