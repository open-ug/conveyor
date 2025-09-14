package models

import (
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/open-ug/conveyor/pkg/types"
)

func setupTestDB(t *testing.T) *badger.DB {
	opts := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opts)
	if err != nil {
		t.Fatalf("failed to open badger db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestLogModel_InsertAndQuery(t *testing.T) {
	db := setupTestDB(t)
	model := &LogModel{DB: db}

	log := types.Log{
		RunID:     "abc123",
		Driver:    "jenkins",
		Pipeline:  "build",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   "Pipeline started",
	}

	// Insert log
	if err := model.Insert(log); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Query by all fields
	logs, err := model.Query("build", "jenkins", "abc123")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}
	if logs[0].Message != log.Message {
		t.Errorf("expected message %q, got %q", log.Message, logs[0].Message)
	}

	// Query by pipeline only
	logs, err = model.Query("build", "", "")
	if err != nil || len(logs) != 1 {
		t.Errorf("expected 1 log for pipeline, got %d, err: %v", len(logs), err)
	}

	// Query by driver only
	logs, err = model.Query("", "jenkins", "")
	if err != nil || len(logs) != 1 {
		t.Errorf("expected 1 log for driver, got %d, err: %v", len(logs), err)
	}

	// Query by runid only
	logs, err = model.Query("", "", "abc123")
	if err != nil || len(logs) != 1 {
		t.Errorf("expected 1 log for runid, got %d, err: %v", len(logs), err)
	}

	// Query with no match
	logs, err = model.Query("other", "jenkins", "abc123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("expected 0 logs, got %d", len(logs))
	}
}
