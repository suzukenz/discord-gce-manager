package config

import (
	"os"
)

// Config is application configurations.
type Config struct {
	GCPProjectID      string
	DiscordToken      string
	DiscordWebhookURL string
}

// NewConfig returns configurations.
func NewConfig() *Config {
	return &Config{
		GCPProjectID:      os.Getenv("PROJECT_ID"),
		DiscordToken:      os.Getenv("DISCORD_TOKEN"),
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK"),
	}
}
