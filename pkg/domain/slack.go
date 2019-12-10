package domain

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)


var c *Config
//this url can be changed into specific webhook address
func init() {
	path := "config"
	c = &Config{}
	if err := c.ReadConfig(path); err != nil {
		log.Println(err)
	}
}

type SlackRequestBody struct {
	Text string `json:"text"`
}

func SendSlackMessage(msg string) {
	c.ResetConfig()
	err := SendSlackNotification(c.Slackwebhook, msg)
	if err != nil {
		log.Fatal(err)
	}
}

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func SendSlackNotification(webhookUrl string, msg string) error {
	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}
