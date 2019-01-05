DISCORD_BOT_NAME := discord-gce-manager
CHECK_ALL_NAME := check-all

BUILD_LINUX_OPTS := GOOS=linux GOARCH=amd64 CGO_ENABLED=0
DOCKER_IMAGE_NAME := suzukenz/discord-gce-manger

DISCORD_TOKEN := MY_DISCORD_TOKEN
PROJECT_ID := MY_GCP_PROJECTID

all: build

deps:
	dep ensure

build/%:
	go build -o bin/$* cmd/$*/main.go

build-bot:
	@make build/$(DISCORD_BOT_NAME)

build-checker:
	@make build/$(CHECK_ALL_NAME)

build: build-bot build-checker

build-linux:
	$(BUILD_LINUX_OPTS) go build -o bin/linux/$(DISCORD_BOT_NAME) cmd/$(DISCORD_BOT_NAME)/main.go
	$(BUILD_LINUX_OPTS) go build -o bin/linux/$(CHECK_ALL_NAME) cmd/$(CHECK_ALL_NAME)/main.go

run-%:
	go run cmd/$*/main.go -t $(DISCORD_TOKEN) -p $(PROJECT_ID)

docker-build: build-linux
	docker build -t $(DOCKER_IMAGE_NAME) .

docker-run-local:
	docker run \
		-e GOOGLE_APPLICATION_CREDENTIALS=/config/mygcloud/application_default_credentials.json \
		-v $(HOME)/.config/gcloud:/config/mygcloud \
		--rm -ti $(DOCKER_IMAGE_NAME) -t $(DISCORD_TOKEN) -p $(PROJECT_ID)

clean:
	rm -rf bin/*
	rm -rf vendor/*
