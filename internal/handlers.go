package internal

import (
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

// Handler is command executer.
type Handler interface {
	validateOptions(cmd string, opts []string) (errMsg string)
	execute(ctx context.Context,
		s *discordgo.Session, m *discordgo.MessageCreate, opts []string) error
}

// Handlers manage and exec Handler.
type Handlers struct {
	handlers map[string]Handler
}

// NewHandlers create and retrun Handlers struct.
func NewHandlers() *Handlers {
	handlers := &Handlers{
		handlers: map[string]Handler{},
	}
	return handlers
}

// Add add Handler to Handlers.
func (handlers *Handlers) Add(command string, h Handler) error {
	_, ok := handlers.handlers[command]
	if ok {
		return fmt.Errorf("error command registration, command is dupulicated")
	}
	handlers.handlers[command] = h
	return nil
}

// Execute parse command from original message and execute command.
func (handlers *Handlers) Execute(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, message string) error {
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
	h, ok := handlers.handlers[cmd]
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

	var servers []*gameServer
	if specifiedTarget {
		svr, err := getGameServer(ctx, datastoreClient, target)
		if err != nil {
			return err
		}
		if svr != nil {
			servers = append(servers, svr)
		}
	} else {
		svrs, err := getGameServers(ctx, datastoreClient)
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
		exIP, _, err := gs.checkServerIsReady(is)

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

	gs, err := getGameServer(ctx, datastoreClient, target)
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
	alreadyRun, err := gs.checkServerIsRunning(is)
	if err != nil {
		return err
	}
	if alreadyRun {
		msg := fmt.Sprintf("%s サーバーは既に起動中です。接続できない場合はもう少し待つか、管理者に連絡してください。", gs.ShowName)
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil
	}

	err = gs.runServer(is)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%s サーバーを起動しました。接続できるようになるまでお待ちください。", gs.ShowName)
	s.ChannelMessageSend(m.ChannelID, msg)

	return nil
}

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

	gs, err := getGameServer(ctx, datastoreClient, target)
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
	running, err := gs.checkServerIsRunning(is)
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

	err = gs.stopServer(is)
	if err != nil {
		return err
	}
	msg = fmt.Sprintf("%s サーバーを停止しました。", gs.ShowName)
	s.ChannelMessageSend(m.ChannelID, msg)

	return nil
}

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

// CheckAllServerWithWebhook get all gameserver difinitions from DataStore and check if server is ready.
// Results will be sented to Discord Channnel by Webhook.
func CheckServersChangedWithWebhook(ctx context.Context) error {
	datastoreClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	servers, err := getGameServers(ctx, datastoreClient)
	if err != nil {
		return err
	}

	if len(servers) == 0 {
		return fmt.Errorf("対象のサーバーが見つかりませんでした。")
	}

	is, err := newGCEInstanceService(ctx)
	if err != nil {
		return err
	}

	for _, gs := range servers {
		var msg string

		exIP, status, _ := gs.checkServerIsReady(is)
		if status != gs.LastStatus {
			if status == StatusReady {
				msg = fmt.Sprintf("%s サーバーが `%s:%d` で起動完了しました。", gs.ShowName, exIP, gs.Port)
			}
			if status == StatusNotReady {
				msg = fmt.Sprintf("%s サーバーが停止しました。", gs.ShowName)
			}

			// update to current status.
			gs.LastStatus = status
			err := gs.saveToDatastore(ctx, datastoreClient)
			if err != nil {
				return err
			}
			log.Println(msg)
			err = sendMessageByWebhook(msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
