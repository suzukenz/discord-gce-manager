package handlers

import (
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/context"

	"github.com/suzukenz/discord-gce-manager/internal/pkg/model"
)

// CheckHandler implements check gameserver command.
type CheckHandler struct{}

func (h *CheckHandler) validateOptions(cmd string, opts []string) (errMsg string) {
	if len(opts) > 1 {
		errMsg = fmt.Sprintf("引数が多すぎます。`例） %[1]s : 全サーバーをチェック　%[1]s 7d2d : 特定のサーバーをチェック(この場合7d2d)`", cmd)
		return errMsg
	}
	return ""
}
func (h *CheckHandler) execute(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, opts []string) error {
	var target string
	var specifiedTarget bool
	if len(opts) > 0 {
		target = opts[0]
		specifiedTarget = true
	}

	// get gameserver difinitions from datastore
	datastoreClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	var servers []*model.GameServer
	if specifiedTarget {
		svr, err := model.GetGameServer(ctx, datastoreClient, target)
		if err != nil {
			return err
		}
		if svr != nil {
			servers = append(servers, svr)
		}
	} else {
		svrs, err := model.GetGameServers(ctx, datastoreClient)
		if err != nil {
			return err
		}
		servers = svrs
	}

	if len(servers) == 0 {
		s.ChannelMessageSend(m.ChannelID, "対象のサーバーが見つかりませんでした。")
		return nil
	}

	is, err := newGCEInstanceService(ctx)
	if err != nil {
		return err
	}

	for _, gs := range servers {
		exIP, _, err := gs.CheckServerIsReady(is, projectID)

		var msg string
		if err != nil {
			log.Println(err)
			msg = fmt.Sprintf("%s サーバーは停止しています。", gs.ShowName)
		} else {
			msg = fmt.Sprintf("%s サーバーは `%s:%d` で起動中です。", gs.ShowName, exIP, gs.Port)
		}
		s.ChannelMessageSend(m.ChannelID, msg)
	}

	return nil
}
