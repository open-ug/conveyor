package log_test

import (
	"fmt"
	"testing"
	"time"

	log "github.com/open-ug/conveyor/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestLokiClient_PushAndQuery(t *testing.T) {
	// Loki endpoint (adjust if in CI or different port)
	lokiURL := "http://localhost:3100"
	if testing.Short() {
		t.Skip("Skipping Loki integration test in short mode.")
	}

	client := log.New(lokiURL)
	assert.NotNil(t, client)

	// Unique label to identify our test log
	labels := map[string]string{
		"job":    "test-job",
		"source": "unit-test",
		"uniq":   time.Now().Format("150405"), // time-based unique label
	}

	message := "Hello Loki from Go test"

	// Push log
	err := client.PushLog(labels, message)
	assert.NoError(t, err, "Expected to push log to Loki")

	// Wait for Loki to index the log
	time.Sleep(2 * time.Second)

	// Query logs (last 30s)
	//start := time.Now().Add(-30 * time.Second)
	//end := time.Now()
	results, err := client.QueryLoki(labels, time.Time{}, time.Time{})
	if err != nil {
		fmt.Printf("Error querying Loki: %v\n", err)
	}
	assert.NoError(t, err, "Expected to query Loki without error")
	assert.NotEmpty(t, results, "Expected to find at least one log entry")

	// Verify message appears in one of the streams
	found := false
	for _, stream := range results {
		for _, entry := range stream.Values {
			if entry[1] == message {
				found = true
				break
			}
		}
	}
	assert.True(t, found, "Expected to find pushed log message in Loki")
}
