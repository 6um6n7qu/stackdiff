package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type slackPayload struct {
	Text string `json:"text"`
}

// sendSlack posts msg to the given Slack webhook URL.
func sendSlack(webhookURL, msg string) error {
	if webhookURL == "" {
		return fmt.Errorf("notify: slack webhook URL is empty")
	}
	payload := slackPayload{Text: msg}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notify: marshal slack payload: %w", err)
	}
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body)) //nolint:gosec
	if err != nil {
		return fmt.Errorf("notify: slack post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify: slack returned status %d", resp.StatusCode)
	}
	return nil
}
