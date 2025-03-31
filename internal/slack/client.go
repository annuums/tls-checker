package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anuums/tls-checker/internal/util"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anuums/tls-checker/internal/config"
	"github.com/anuums/tls-checker/internal/tlscheck"
)

type SlackClient struct {
	token string
}

func NewSlackClient(token string) *SlackClient {

	return &SlackClient{token: token}
}

func (c *SlackClient) SendAlert(vc *tlscheck.ValidateCert) error {

	channel := config.RuntimeConf.ChannelID
	tz := config.RuntimeConf.Timezone
	colorBar := config.RuntimeConf.ColorBar

	if tz == "" {

		tz = "Etc/UTC"
	}

	if colorBar != "" && util.IsValidRgb(colorBar) == false {

		colorBar = ""
	}

	log.Printf("Sending alert for [%s].\n", vc.Hostname)

	now := time.Now()
	loc, err := time.LoadLocation(tz)
	if err != nil {

		return fmt.Errorf("failed to load time location :: %w", err)
	}

	header := fmt.Sprintf("[%s] SSL/TLS Expiration Alert", vc.Hostname)
	payload := Payload{
		Channel:     channel,
		Color:       colorBar,
		Description: header,
		Header:      header,
		Text: []Text{
			{
				Text: "Hostname",
				Bold: true,
			},
			{
				Text: vc.Hostname,
				Bold: false,
			},
			{
				Text: "Expired At",
				Bold: true,
			},
			{
				Text: vc.Cert.NotAfter.In(loc).Format(time.DateTime),
				Bold: false,
			},
			{
				Text: "Remaining Until",
				Bold: true,
			},
			{
				Text: fmt.Sprintf("%v Days", int(vc.Cert.NotAfter.In(loc).Sub(now).Hours()/24)),
				Bold: false,
			},
			{
				Text: "Certificate Scope",
				Bold: true,
			},
			{
				Text: strings.Join(vc.Cert.DNSNames, ", "),
				Bold: false,
			},
		},
	}

	slackPayload := payload.ToSlackMessage()

	_, err = c.postMessageToSlack("chat.postMessage", slackPayload)
	if err != nil {

		return fmt.Errorf("failed to post slack message :: %w", err)
	}

	return nil
}

func (c *SlackClient) postMessageToSlack(apiMethod string, message *SlackMessage) (*SlackApiResponse, error) {

	slackToken := config.RuntimeConf.SlackConfig.Token

	if slackToken == "" {

		return nil, errors.New("SLACK_TOKEN environment variable not set")
	}

	url := fmt.Sprintf("https://slack.com/api/%s", apiMethod)

	jsonValue, err := json.Marshal(message)
	if err != nil {

		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {

		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+slackToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {

		if resp.Body != nil {

			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var response *SlackApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {

		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.OK == false {

		debugBody, _ := json.MarshalIndent(response, "", "  ")
		return response, fmt.Errorf(
			"failed to post message :: %s\nFull response:\n%s",
			response.Error,
			string(debugBody),
		)
	}

	return response, nil
}
