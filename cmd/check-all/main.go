package main

import (
	"context"
	"flag"
	"log"

	"github.com/suzukenz/discord-gce-manager/internal"
)

// Variables used for command line parameters
var (
	token      string
	projectID  string
	webhookURL string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.StringVar(&projectID, "p", "", "GCP ProjectID")
	flag.StringVar(&webhookURL, "d", "", "Discord Webhook URL")
	flag.Parse()

	internal.SetProjectID(projectID)
	internal.SetWebhookURL(webhookURL)
}

func main() {
	ctx := context.Background()
	err := internal.CheckAllServerWithWebhook(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
