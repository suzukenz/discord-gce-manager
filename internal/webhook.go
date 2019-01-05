package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func sendMessageByWebhook(message string) error {
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
