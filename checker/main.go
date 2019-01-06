package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/suzukenz/discord-gce-manager/internal"
)

// Variables used for command line parameters
var (
	token      string
	projectID  string
	webhookURL string
)

func init() {
	err := godotenv.Load("configs.env")
	if err != nil {
		log.Fatal("Error loading configs.env file")
	}

	projectID = os.Getenv("PROJECT_ID")
	token = os.Getenv("DISCORD_TOKEN")
	webhookURL = os.Getenv("DISCORD_WEBHOOK")

	internal.SetProjectID(projectID)
	internal.SetWebhookURL(webhookURL)
}

func main() {
	ctx := context.Background()
	err := internal.CheckServersChangedWithWebhook(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
