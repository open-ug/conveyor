package models_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/pkg/types"
)

func TestInsertDriverMessage(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)
	defer db.Close()

	// Create a new DriverMessageModel instance
	dm := models.NewDriverMessageModel(nil, db)

	// Define a test driver message
	testMessage := types.DriverMessage{
		ID:      "test-message-id",
		Payload: "This is a test message",
		Event:   "test-event",
	}

	t.Run("insert driver message", func(t *testing.T) {
		err := dm.BadgerInsert(testMessage)
		if err != nil {
			t.Fatalf("Failed to insert driver message: %v", err)
		}

		// Verify that the driver message was inserted in BadgerDB
		err = db.View(func(txn *badger.Txn) error {
			key := fmt.Sprintf("driver-messages/%s", testMessage.ID)
			item, err := txn.Get([]byte(key))
			if err != nil {
				return err
			}

			var storedMessage types.DriverMessage
			err = item.Value(func(val []byte) error {
				return json.Unmarshal(val, &storedMessage)
			})
			if err != nil {
				return err
			}

			if !reflect.DeepEqual(storedMessage, testMessage) {
				t.Errorf("Stored driver message does not match expected. Got %+v, want %+v", storedMessage, testMessage)
			}

			return nil
		})
		if err != nil {
			t.Fatalf("Failed to verify driver message in BadgerDB: %v", err)
		}
	})
}

func TestFindOneDriverMessage(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)
	defer db.Close()

	// Create a new DriverMessageModel instance
	dm := models.NewDriverMessageModel(nil, db)

	// Define a test driver message
	testMessage := types.DriverMessage{
		ID:      "test-message-id",
		Payload: "This is a test message",
		Event:   "test-event",
	}

	// Insert the test message into BadgerDB
	err := dm.BadgerInsert(testMessage)
	if err != nil {
		t.Fatalf("Failed to insert driver message: %v", err)
	}

	t.Run("find one driver message", func(t *testing.T) {
		foundMessage, err := dm.BadgerFindOne(testMessage.ID)
		if err != nil {
			t.Fatalf("Failed to find driver message: %v", err)
		}

		if !reflect.DeepEqual(*foundMessage, testMessage) {
			t.Errorf("Found driver message does not match expected. Got %+v, want %+v", *foundMessage, testMessage)
		}
	})
}

func TestBadgerFindAllDriverMessages(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)
	defer db.Close()

	// Create a new DriverMessageModel instance
	dm := models.NewDriverMessageModel(nil, db)

	// Define test driver messages
	testMessages := []types.DriverMessage{
		{ID: "msg1", Payload: "Message 1", Event: "event1"},
		{ID: "msg2", Payload: "Message 2", Event: "event2"},
	}

	// Insert the test messages into BadgerDB
	for _, msg := range testMessages {
		err := dm.BadgerInsert(msg)
		if err != nil {
			t.Fatalf("Failed to insert driver message %s: %v", msg.ID, err)
		}
	}

	t.Run("find all driver messages", func(t *testing.T) {
		foundMessages, err := dm.BadgerFindAll()
		if err != nil {
			t.Fatalf("Failed to find all driver messages: %v", err)
		}

		if len(foundMessages) != len(testMessages) {
			t.Errorf("Expected %d driver messages, got %d", len(testMessages), len(foundMessages))
		}

		for _, msg := range testMessages {
			found := false
			for _, foundMsg := range foundMessages {
				if reflect.DeepEqual(msg, foundMsg) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Driver message %+v not found in results", msg)
			}
		}
	})
}

func TestBadgerDeleteDriverMessage(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)
	defer db.Close()

	// Create a new DriverMessageModel instance
	dm := models.NewDriverMessageModel(nil, db)

	// Define a test driver message
	testMessage := types.DriverMessage{
		ID:      "test-message-id",
		Payload: "This is a test message",
		Event:   "test-event",
	}

	// Insert the test message into BadgerDB
	err := dm.BadgerInsert(testMessage)
	if err != nil {
		t.Fatalf("Failed to insert driver message: %v", err)
	}

	t.Run("delete driver message", func(t *testing.T) {
		err := dm.BadgerDeleteOne(testMessage.ID)
		if err != nil {
			t.Fatalf("Failed to delete driver message: %v", err)
		}

		// Verify that the driver message was deleted from BadgerDB
		err = db.View(func(txn *badger.Txn) error {
			key := fmt.Sprintf("driver-messages/%s", testMessage.ID)
			_, err := txn.Get([]byte(key))
			if err == nil {
				t.Errorf("Expected error when retrieving deleted driver message, but got none")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Failed to verify deletion of driver message in BadgerDB: %v", err)
		}
	})
}
