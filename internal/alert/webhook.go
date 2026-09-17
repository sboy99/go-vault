package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sboy99/go-vault/pkg/logger"
)

// NotifyFailure posts a JSON payload to webhookURL if set.
func NotifyFailure(webhookURL, event, message string) {
	if webhookURL == "" {
		return
	}
	payload, _ := json.Marshal(map[string]string{
		"event":   event,
		"message": message,
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(payload))
	if err != nil {
		logger.Warn("alert webhook request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Warn("alert webhook send: %v", err)
		return
	}
	_ = resp.Body.Close()
}
