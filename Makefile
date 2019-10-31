NAME := github
CGO_ENABLED = 0
GO := go
BUILD_TARGET = build
BUILDFLAGS = -ldflags

build:
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOARCH=amd64 $(GO) $(BUILD_TARGET) -o bin/$(NAME)

linux:
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOARCH=amd64 GOOS=linux $(GO) $(BUILD_TARGET) -o bin/linux/$(NAME)

image:
	docker build -t surenpi/github-proxy:test .

push-image:
	docker push surenpi/github-proxy:test