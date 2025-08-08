package utils

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
func NewLokiClient(url string) *LokiClient {
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

func formatLabels(labels map[string]string) string {
	var formattedLabels string
	for k, v := range labels {
		formattedLabels = fmt.Sprintf("%s%s=\"%s\",", formattedLabels, k, v)
	}
	return formattedLabels[:len(formattedLabels)-1] // Remove the trailing comma
}

// QueryLoki queries Loki for logs based on the provided labels and time range. The time range is optional.
func (lc *LokiClient) QueryLoki(labels map[string]string, start, end time.Time) ([]LokiStream, error) {
	query := fmt.Sprintf("{%s}", formatLabels(labels))
	if !start.IsZero() && !end.IsZero() {
		query += fmt.Sprintf(" @timestamp >= %d and @timestamp <= %d", start.UnixNano(), end.UnixNano())
	}

	fmt.Println("Querying Loki with query:", query)

	// direction Forwards
	resp, err := lc.Client.Get(fmt.Sprintf("%s/loki/api/v1/query_range?query=%s&direction=FORWARD", lc.URL, query))
	if err != nil {
		return nil, fmt.Errorf("failed to query loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("loki responded with status %s", resp.Status)
	}

	var result struct {
		Data struct {
			Result []LokiStream `json:"result"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode loki response: %w", err)
	}

	return result.Data.Result, nil
}
