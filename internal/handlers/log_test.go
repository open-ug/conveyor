package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/pkg/types"
)

func setupTestHandler(t *testing.T) *LogHandler {
	opts := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opts)
	if err != nil {
		t.Fatalf("failed to open badger db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	model := &models.LogModel{DB: db}
	return &LogHandler{Model: model}
}

func TestCreateLogAndGetLogs(t *testing.T) {
	app := fiber.New()
	handler := setupTestHandler(t)
	app.Post("/logs", handler.CreateLog)
	app.Get("/logs", handler.GetLogs)

	log := types.Log{
		RunID:     "abc123",
		Driver:    "jenkins",
		Pipeline:  "build",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   "Pipeline started",
	}
	body, _ := json.Marshal(log)

	// Test POST /logs
	req := httptest.NewRequest(http.MethodPost, "/logs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /logs failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}

	// Test GET /logs with all filters
	url := "/logs?pipeline=build&driver=jenkins&runid=abc123"
	req = httptest.NewRequest(http.MethodGet, url, nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("GET /logs failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	var logs []types.Log
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("expected 1 log, got %d", len(logs))
	}
	if logs[0].Message != log.Message {
		t.Errorf("expected message %q, got %q", log.Message, logs[0].Message)
	}

	// Test GET /logs with partial filter
	url = "/logs?pipeline=build"
	req = httptest.NewRequest(http.MethodGet, url, nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("GET /logs partial filter failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	logs = nil
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("expected 1 log for pipeline, got %d", len(logs))
	}

	// Test GET /logs with no match
	url = "/logs?pipeline=other"
	req = httptest.NewRequest(http.MethodGet, url, nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("GET /logs no match failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	logs = nil
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("expected 0 logs, got %d", len(logs))
	}
}
