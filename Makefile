BUILD_LINUX_OPTS := GOOS=linux GOARCH=amd64 CGO_ENABLED=0

# All go source files
SOURCES := $(shell find . -name '*.go' -not -name '*_test.go')

# List of binary cmds to build
CMDS := \
	cmd/discord-bot \
	cmd/scheduled-checker


all: build

$(CMDS): $(SOURCES)
	$(BUILD_LINUX_OPTS) go build -o ./bin/$(shell basename "$@") $@/main.go

build:
	@for d in $(CMDS); do $(MAKE) $$d; done

clean:
	rm -rf vendor/*
	rm -rf bin/*

deps:
	dep ensure

docker-build: build
	docker-compose -f ./docker/docker-compose.yml build
