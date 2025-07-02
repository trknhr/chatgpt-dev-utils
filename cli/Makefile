APP_NAME := chatgpt-dev-utils
BIN_DIR := $(GOPATH)/bin
GO_FILES := main.go

.PHONY: all build install clean run test

all: build

build:
	go build -o $(APP_NAME) $(GO_FILES)

install:
	go install .

run:
	go run $(GO_FILES)

debug:
	dlv debug --headless --api-version=2 --listen=127.0.0.1:43000 .

clean:
	rm -f $(APP_NAME)

test:
	go test ./...
