NAME := app
BUILD_CMD ?= CGO_ENABLED=0 go build -o bin/${NAME} -ldflags '-v -w -s' ./cmd

.PHONY: run
run:
	go run cmd/main.go

.PHONY: test
test:
	go test ./... -cover

.PHONY: build
build:
	echo "building"
	${BUILD_CMD}
	echo "build done"