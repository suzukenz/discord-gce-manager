package handlers

import (
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/context"

	"github.com/suzukenz/discord-gce-manager/internal/pkg/model"
)

// StopHandler implements stop gameserver command.
type StopHandler struct{}

func (h *StopHandler) validateOptions(cmd string, opts []string) (errMsg string) {
	if len(opts) != 1 {
		errMsg = fmt.Sprintf("引数が足りないか多すぎます。`例） %[1]s 7d2d : 特定のサーバーを停止(この場合7d2d)`", cmd)
		return errMsg
	}
	return ""
}
func (h *StopHandler) execute(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, opts []string) error {
	target := opts[0]

	datastoreClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	gs, err := model.GetGameServer(ctx, datastoreClient, target)
	if err != nil {
		return err
	}
	if gs == nil {
		s.ChannelMessageSend(m.ChannelID, "対象のサーバーが見つかりませんでした。")
		return nil
	}

	is, err := newGCEInstanceService(ctx)
	if err != nil {
		return err
	}
	running, err := gs.CheckServerIsRunning(is, projectID)
	if err != nil {
		return err
	}
	if !running {
		msg := fmt.Sprintf("%s サーバーは既に停止中です。", gs.ShowName)
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil
	}

	msg := fmt.Sprintf("30秒後に %s サーバーを停止します。", gs.ShowName)
	s.ChannelMessageSend(m.ChannelID, msg)
	time.Sleep(30 * time.Second)

	err = gs.StopServer(is, projectID)
	if err != nil {
		return err
	}
	msg = fmt.Sprintf("%s サーバーの停止を開始しました。", gs.ShowName)
	s.ChannelMessageSend(m.ChannelID, msg)

	return nil
}
