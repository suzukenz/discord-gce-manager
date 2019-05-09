package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/context"
)

// CheckChannelIDHandler implements check channel id.
type CheckChannelIDHandler struct{}

func (h *CheckChannelIDHandler) validateOptions(cmd string, opts []string) (errMsg string) {
	return ""
}
func (h *CheckChannelIDHandler) execute(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, opts []string) error {
	msg := fmt.Sprintf("このチャンネルのIDは `%s` です。", m.ChannelID)
	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}
