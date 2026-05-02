package models_test

import (
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/open-ug/conveyor/internal/models"
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

func Test_Resource_Insert(t *testing.T) {
	db := setupTestDB(t)

	resourcemodel := models.NewResourceModel(nil, db)

	t.Run("Insert Resource", func(t *testing.T) {
		err := resourcemodel.BadgerDBInsert("test-resource", "test-type", []byte("test-data"))
		if err != nil {
			t.Fatalf("failed to insert resource: %v", err)
		}

		// Verify the resource was inserted
		err = db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("/resources/test-type/test-resource"))
			if err != nil {
				return err
			}
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			if string(val) != "test-data" {
				t.Errorf("expected 'test-data', got '%s'", string(val))
			}
			return nil
		})
		if err != nil {
			t.Fatalf("failed to verify resource: %v", err)
		}
	})

	// verify versioned resource was inserted
	t.Run("Verify Versioned Resource", func(t *testing.T) {
		err := db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("/resources/test-type/test-resource/1"))
			if err != nil {
				return err
			}
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			if string(val) != "test-data" {
				t.Errorf("expected 'test-data', got '%s'", string(val))
			}
			return nil
		})
		if err != nil {
			t.Fatalf("failed to verify versioned resource: %v", err)
		}
	})

}
