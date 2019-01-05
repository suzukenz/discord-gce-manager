package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
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
	// Register commands
	handlers := internal.NewHandlers()
	for _, cmds := range []struct {
		command string
		handler internal.Handler
	}{
		{command: "/check", handler: new(internal.CheckHandler)},
		{command: "/run", handler: new(internal.RunHandler)},
		{command: "/stop", handler: new(internal.StopHandler)},
		{command: "/channel", handler: new(internal.CheckChannelIDHandler)},
	} {
		err := handlers.Add(cmds.command, cmds.handler)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(createMessageCreatedHandler(handlers))
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
		return
	}

	// Start http server for execute manually by request.
	srv := startHTTPServer()

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Discord bot server is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

	// Http server shutdown
	if err := srv.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Fatalln("Failed to http server shutdown:", err)
	}
}

// Return function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func createMessageCreatedHandler(handlers *internal.Handlers) func(s *discordgo.Session, m *discordgo.MessageCreate) {
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

func httpCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	err := internal.CheckServerChangedWithWebhook(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(w, "execute check ok")
}

func startHTTPServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks/check", httpCheckHandler)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	return srv
}
