package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/suzukenz/discord-gce-manager/internal/discord-bot/handlers"
	"github.com/suzukenz/discord-gce-manager/internal/pkg/config"
)

func main() {
	// Register commands
	hdls := handlers.NewHandlers()
	for _, cmds := range []struct {
		command string
		handler handlers.Handler
	}{
		{command: "/check", handler: new(handlers.CheckHandler)},
		{command: "/run", handler: new(handlers.RunHandler)},
		{command: "/stop", handler: new(handlers.StopHandler)},
		{command: "/channel", handler: new(handlers.CheckChannelIDHandler)},
	} {
		err := hdls.Add(cmds.command, cmds.handler)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}

	// Create a new Discord session using the provided bot token.
	cfg := config.NewConfig()
	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalln("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(createMessageCreatedHandler(hdls))

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Discord bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// Return function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func createMessageCreatedHandler(handlers *handlers.Handlers) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		msg := m.Content
		ctx := context.Background()
		err := handlers.Execute(ctx, s, m, msg)
		if err != nil {
			log.Printf("fail command, received: %s, err: %s", msg, err)
			s.ChannelMessageSend(m.ChannelID, "コマンド実行に失敗しました。")
		}
	}
}
