GOFILES = $(shell find ./cmd -name '*.go' -not -path './vendor/*')
# GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

BINARY_PATH = "bin/discord-gce-manager"
IMAGE_NAME = "suzukenz/discord-gce-manger"

DISCORD_TOKEN = "MY DISCORD TOKEN"
PROJECT_ID = "MY GCP PROJECTID"

run:
	go run cmd/main.go -t $(DISCORD_TOKEN) -p $(PROJECT_ID)

run-bin: build
	$(BINARY_PATH) -t $(DISCORD_TOKEN) -p $(PROJECT_ID)

build:
	go build -o $(BINARY_PATH) $(GOFILES)

build-linux: 
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_PATH) $(GOFILES)

docker-build: build-linux
	docker build -t $(IMAGE_NAME) .

docker-run-local:
	docker run \
		-e GOOGLE_APPLICATION_CREDENTIALS=/config/mygcloud/application_default_credentials.json \
		-v $(HOME)/.config/gcloud:/config/mygcloud \
		--rm -ti $(IMAGE_NAME) -t $(DISCORD_TOKEN) -p $(PROJECT_ID)
