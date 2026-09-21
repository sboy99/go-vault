package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sboy99/go-vault/internal/domain"
	"github.com/sboy99/go-vault/pkg/logger"
)

// WebhookNotifier posts failure events to a webhook URL.
type WebhookNotifier struct {
	URL string
}

var _ domain.Notifier = (*WebhookNotifier)(nil)

func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{URL: url}
}

// NotifyFailure posts a JSON payload to the webhook URL if set.
func (n *WebhookNotifier) NotifyFailure(ctx context.Context, event, message string) {
	if n == nil || n.URL == "" {
		return
	}
	payload, _ := json.Marshal(map[string]string{
		"event":   event,
		"message": message,
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.URL, bytes.NewReader(payload))
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
