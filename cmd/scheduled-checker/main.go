package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"

	"github.com/suzukenz/discord-gce-manager/internal/pkg/config"
	"github.com/suzukenz/discord-gce-manager/internal/pkg/model"
)

func main() {
	ctx := context.Background()
	err := checkServersChangedWithWebhook(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

// checkServersChangedWithWebhook get all gameserver definitions from DataStore and check if server is ready.
// Results that status chenged will be sented to Discord Channel by Webhook.
func checkServersChangedWithWebhook(ctx context.Context) error {
	cfg := config.NewConfig()
	projectID := cfg.GCPProjectID

	datastoreClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	servers, err := model.GetGameServers(ctx, datastoreClient)
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

		exIP, status, _ := gs.CheckServerIsReady(is, projectID)
		if status != gs.LastStatus {
			if status == model.StatusReady {
				msg = fmt.Sprintf("%s サーバーが `%s:%d` で起動完了しました。", gs.ShowName, exIP, gs.Port)
			}
			if status == model.StatusNotReady {
				msg = fmt.Sprintf("%s サーバーが停止しました。", gs.ShowName)
			}

			// update to current status.
			gs.LastStatus = status
			err := gs.SaveToDatastore(ctx, datastoreClient)
			if err != nil {
				return err
			}
			log.Println(msg)
			err = sendMessageByWebhook(cfg.DiscordWebhookURL, msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func sendMessageByWebhook(webhookURL, message string) error {
	jsonStr := `{"content":"` + message + `"}`

	req, err := http.NewRequest(
		"POST",
		webhookURL,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("finish request webhookURL, status: %s, body: %s", resp.Status, string(bytes))
	log.Println(msg)

	return err
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
