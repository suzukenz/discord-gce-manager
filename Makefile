APP_NAME := discord-gce-manager

BUILD_LINUX_OPTS := GOOS=linux GOARCH=amd64 CGO_ENABLED=0
DOCKER_IMAGE_NAME := suzukenz/discord-gce-manger

PROJECT_ID := MY_GCP_PROJECTID
DISCORD_TOKEN := MY_DISCORD_TOKEN
DISCORD_WEBHOOK := MY_DISCORD_WEBHOOKURL

all: build

deps:
	dep ensure

build:
	go build -o bin/$(APP_NAME) cmd/$*/main.go

build-linux:
	$(BUILD_LINUX_OPTS) go build -o bin/linux/$(APP_NAME) .

run:
	go run ./main.go -p $(PROJECT_ID) -t $(DISCORD_TOKEN) -d $(DISCORD_WEBHOOK)

docker-build: build-linux
	docker build \
		-t $(DOCKER_IMAGE_NAME) .

docker-run:
	docker run \
		-e GOOGLE_APPLICATION_CREDENTIALS=/config/mygcloud/application_default_credentials.json \
		-v $(HOME)/.config/gcloud:/config/mygcloud \
		-p 8080:8080 \
		--rm -ti $(DOCKER_IMAGE_NAME) -p $(PROJECT_ID) -t $(DISCORD_TOKEN) -d $(DISCORD_WEBHOOK)

deploy: deploy-app deploy-cron

deploy-app: docker-build
	gcloud app deploy

deploy-cron:
	gcloud app deploy cron.yaml

clean:
	rm -rf bin/*
	rm -rf vendor/*
