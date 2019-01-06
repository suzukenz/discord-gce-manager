include ./configs.env

DISCORD_BOT_NAME := discord-bot
CHECK_ALL_NAME := checker

BUILD_LINUX_OPTS := GOOS=linux GOARCH=amd64 CGO_ENABLED=0

all: build

deps:
	dep ensure

build/%:
	go build -o bin/$* $*/main.go

build-bot:
	@make build/$(DISCORD_BOT_NAME)

build-checker:
	@make build/$(CHECK_ALL_NAME)

build: build-bot build-checker

build-linux:
	$(BUILD_LINUX_OPTS) go build -o bin/linux/$(DISCORD_BOT_NAME) $(DISCORD_BOT_NAME)/main.go
	$(BUILD_LINUX_OPTS) go build -o bin/linux/$(CHECK_ALL_NAME) $(CHECK_ALL_NAME)/main.go

run-bot:
	go run $(DISCORD_BOT_NAME)/main.go

run-checker:
	go run $(CHECK_ALL_NAME)/main.go

docker-build: build-linux
	docker-compose build

docker-run:
	docker-compose up

docker-run-bot:
	docker-compose run $(DISCORD_BOT_NAME)

docker-run-checker:
	docker-compose run $(CHECK_ALL_NAME)

clean:
	rm -rf bin/*
	rm -rf vendor/*
