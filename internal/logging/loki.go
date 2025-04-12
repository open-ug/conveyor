package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LokiClient struct {
	URL    string
	Client *http.Client
}

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LokiPayload struct {
	Streams []LokiStream `json:"streams"`
}

// New creates a new LokiClient.
func New(url string) *LokiClient {
	return &LokiClient{
		URL:    url,
		Client: &http.Client{Timeout: 5 * time.Second},
	}
}

// PushLog sends a single log line to Loki.
func (lc *LokiClient) PushLog(labels map[string]string, message string) error {
	entry := LokiPayload{
		Streams: []LokiStream{
			{
				Stream: labels,
				Values: [][]string{
					{formatNano(time.Now()), message},
				},
			},
		},
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}

	resp, err := lc.Client.Post(lc.URL+"/loki/api/v1/push", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to push log: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("loki responded with status %s", resp.Status)
	}

	return nil
}

func formatNano(t time.Time) string {
	return fmt.Sprintf("%d", t.UnixNano())
}
