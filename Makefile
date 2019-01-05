GOFILES = $(shell find ./cmd -name '*.go' -not -path './vendor/*')
# GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

BINARY_PATH = "bin/discord-gce-manager"

DISCORD_TOKEN = "MY DISCORD TOKEN"
PROJECT_ID = "MY GCP PROJECTID"

run:
	go run cmd/main.go -t $(DISCORD_TOKEN) -p $(PROJECT_ID)

run-bin: build
	$(BINARY_PATH) -t $(DISCORD_TOKEN) -p $(PROJECT_ID)

build:
	go build -o $(BINARY_PATH) $(GOFILES)

build-linux: 
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_PATH) .
