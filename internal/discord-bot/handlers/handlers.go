package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"

	"github.com/suzukenz/discord-gce-manager/internal/pkg/config"
)

var (
	projectID string
)

func init() {
	cfg := config.NewConfig()
	projectID = cfg.GCPProjectID
}

// Handler is command executer.
type Handler interface {
	validateOptions(cmd string, opts []string) (errMsg string)
	execute(ctx context.Context,
		s *discordgo.Session, m *discordgo.MessageCreate, opts []string) error
}

// Handlers manage and exec Handler.
type Handlers struct {
	Handlers map[string]Handler
}

// NewHandlers create and retrun Handlers struct.
func NewHandlers() *Handlers {
	Handlers := &Handlers{
		Handlers: map[string]Handler{},
	}
	return Handlers
}

// Add add Handler to Handlers.
func (Handlers *Handlers) Add(command string, h Handler) error {
	_, ok := Handlers.Handlers[command]
	if ok {
		return fmt.Errorf("error command registration, command is dupulicated")
	}
	Handlers.Handlers[command] = h
	return nil
}

// Execute parse command from original message and execute command.
func (Handlers *Handlers) Execute(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, message string) error {
	if !strings.HasPrefix(message, "/") {
		// ignore message that not command
		return nil
	}

	// extract command
	cmds := strings.Split(message, " ")
	var options []string
	if len(cmds) > 1 {
		options = cmds[1:]
	}

	cmd := cmds[0]
	h, ok := Handlers.Handlers[cmd]
	if !ok {
		// ignore message that not command
		log.Printf("command not found, received: %s", cmd)
		return nil
	}

	ret := h.validateOptions(cmd, options)
	if len(ret) > 0 {
		s.ChannelMessageSend(m.ChannelID, ret)
		return nil
	}

	return h.execute(ctx, s, m, options)
}

func newGCEInstanceService(ctx context.Context) (*compute.InstancesService, error) {
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		return nil, err
	}
	computeService, err := compute.New(client)
	if err != nil {
		return nil, fmt.Errorf("error opening connection, err: %s", err)
	}

	return compute.NewInstancesService(computeService), nil
}
