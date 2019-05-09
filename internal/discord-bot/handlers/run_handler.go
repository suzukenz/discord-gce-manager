package handlers

import (
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/context"

	"github.com/suzukenz/discord-gce-manager/internal/pkg/model"
)

// RunHandler implements run gameserver command.
type RunHandler struct{}

func (h *RunHandler) validateOptions(cmd string, opts []string) (errMsg string) {
	if len(opts) != 1 {
		errMsg = fmt.Sprintf("引数が足りないか多すぎます。`例） %[1]s 7d2d : 特定のサーバーを起動(この場合7d2d)`", cmd)
		return errMsg
	}
	return ""
}
func (h *RunHandler) execute(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, opts []string) error {
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
	alreadyRun, err := gs.CheckServerIsRunning(is, projectID)
	if err != nil {
		return err
	}
	if alreadyRun {
		msg := fmt.Sprintf("%s サーバーは既に起動中です。接続できない場合はもう少し待つか、管理者に連絡してください。", gs.ShowName)
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil
	}

	err = gs.RunServer(is, projectID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%s サーバーを起動しました。接続できるようになるまでお待ちください。", gs.ShowName)
	s.ChannelMessageSend(m.ChannelID, msg)

	return nil
}
