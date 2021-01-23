.PHONY: build
build:
	go build -v ./src/mediaserver.go

.DEFAULT_GOAL: build
