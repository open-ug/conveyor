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

func TestCreatePipeline(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)

	defer db.Close()

	// Create a new PipelineModel instance
	pm := models.NewPipelineModel(nil, db)

	// Define a test pipeline
	testPipeline := &types.Pipeline{
		Name:     "test-pipeline",
		Resource: "test-res",
	}

	t.Run("create pipeline", func(t *testing.T) {
		err := pm.BadgerCreatePipeline(testPipeline)
		if err != nil {
			t.Fatalf("Failed to create pipeline: %v", err)
		}

		// Verify that the pipeline was created in BadgerDB
		err = db.View(func(txn *badger.Txn) error {
			pipelineKey := fmt.Sprintf("/pipelines/%s", testPipeline.Name)
			item, err := txn.Get([]byte(pipelineKey))
			if err != nil {
				return err
			}

			var storedPipeline types.Pipeline
			err = item.Value(func(val []byte) error {
				return json.Unmarshal(val, &storedPipeline)
			})
			if err != nil {
				return err
			}

			if !reflect.DeepEqual(storedPipeline, *testPipeline) {
				t.Errorf("Stored pipeline does not match expected. Got %+v, want %+v", storedPipeline, *testPipeline)
			}

			return nil
		})
		if err != nil {
			t.Fatalf("Failed to verify pipeline in BadgerDB: %v", err)
		}
	})
}

func TestBadgerDBGetPipeline(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)
	defer db.Close()

	// Create a new PipelineModel instance
	pm := models.NewPipelineModel(nil, db)

	// Define a test pipeline
	testPipeline := &types.Pipeline{
		Name:     "test-pipeline",
		Resource: "test-res",
	}

	// Create the pipeline in BadgerDB
	err := pm.BadgerCreatePipeline(testPipeline)
	if err != nil {
		t.Fatalf("Failed to create pipeline: %v", err)
	}

	t.Run("get pipeline", func(t *testing.T) {
		var retrievedPipeline types.Pipeline

		err := db.View(func(txn *badger.Txn) error {
			pipelineKey := fmt.Sprintf("/pipelines/%s", testPipeline.Name)
			item, err := txn.Get([]byte(pipelineKey))
			if err != nil {
				return err
			}

			return item.Value(func(val []byte) error {
				return json.Unmarshal(val, &retrievedPipeline)
			})
		})
		if err != nil {
			t.Fatalf("Failed to retrieve pipeline from BadgerDB: %v", err)
		}

		if !reflect.DeepEqual(retrievedPipeline, *testPipeline) {
			t.Errorf("Retrieved pipeline does not match expected. Got %+v, want %+v", retrievedPipeline, *testPipeline)
		}
	})

	t.Run("check non-existing pipeline", func(t *testing.T) {
		_, err := pm.BadgerGetPipeline("non-existing-pipeline")
		if err == nil {
			t.Fatalf("Expected error when retrieving non-existing pipeline, but got none")
		}
	})
}

func TestBadgerUpdatePipeline(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)
	defer db.Close()

	// Create a new PipelineModel instance
	pm := models.NewPipelineModel(nil, db)

	// Define a test pipeline
	testPipeline := &types.Pipeline{
		Name:     "test-pipeline",
		Resource: "test-res",
	}

	// Create the pipeline in BadgerDB
	err := pm.BadgerCreatePipeline(testPipeline)
	if err != nil {
		t.Fatalf("Failed to create pipeline: %v", err)
	}

	t.Run("update pipeline", func(t *testing.T) {
		// Update the pipeline's resource
		testPipeline.Resource = "updated-res"
		err := pm.BadgerUpdatePipeline(testPipeline)
		if err != nil {
			t.Fatalf("Failed to update pipeline: %v", err)
		}

		var updatedPipeline types.Pipeline

		err = db.View(func(txn *badger.Txn) error {
			pipelineKey := fmt.Sprintf("/pipelines/%s", testPipeline.Name)
			item, err := txn.Get([]byte(pipelineKey))
			if err != nil {
				return err
			}

			return item.Value(func(val []byte) error {
				return json.Unmarshal(val, &updatedPipeline)
			})
		})
		if err != nil {
			t.Fatalf("Failed to retrieve updated pipeline from BadgerDB: %v", err)
		}

		if !reflect.DeepEqual(updatedPipeline, *testPipeline) {
			t.Errorf("Updated pipeline does not match expected. Got %+v, want %+v", updatedPipeline, *testPipeline)
		}
	})
}

func TestBadgerDeletePipeline(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)
	defer db.Close()

	// Create a new PipelineModel instance
	pm := models.NewPipelineModel(nil, db)

	// Define a test pipeline
	testPipeline := &types.Pipeline{
		Name:     "test-pipeline",
		Resource: "test-res",
	}

	// Create the pipeline in BadgerDB
	err := pm.BadgerCreatePipeline(testPipeline)
	if err != nil {
		t.Fatalf("Failed to create pipeline: %v", err)
	}

	t.Run("delete pipeline", func(t *testing.T) {
		err := pm.BadgerDeletePipeline(testPipeline.Name)
		if err != nil {
			t.Fatalf("Failed to delete pipeline: %v", err)
		}

		// Verify that the pipeline was deleted from BadgerDB
		err = db.View(func(txn *badger.Txn) error {
			pipelineKey := fmt.Sprintf("/pipelines/%s", testPipeline.Name)
			_, err := txn.Get([]byte(pipelineKey))
			if err == nil {
				t.Errorf("Expected error when retrieving deleted pipeline, but got none")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Failed to verify deletion of pipeline in BadgerDB: %v", err)
		}
	})
}

func TestBadgerListPipelines(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)
	defer db.Close()

	// Create a new PipelineModel instance
	pm := models.NewPipelineModel(nil, db)

	// Define test pipelines
	testPipelines := []*types.Pipeline{
		{Name: "pipeline-1", Resource: "res-1"},
		{Name: "pipeline-2", Resource: "res-2"},
	}

	// Create the pipelines in BadgerDB
	for _, pipeline := range testPipelines {
		err := pm.BadgerCreatePipeline(pipeline)
		if err != nil {
			t.Fatalf("Failed to create pipeline: %v", err)
		}
	}

	t.Run("list pipelines", func(t *testing.T) {
		var retrievedPipelines []*types.Pipeline

		err := db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = true
			it := txn.NewIterator(opts)
			defer it.Close()

			prefix := []byte("/pipelines/")
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				item := it.Item()
				var pipeline types.Pipeline
				err := item.Value(func(val []byte) error {
					return json.Unmarshal(val, &pipeline)
				})
				if err != nil {
					return err
				}
				retrievedPipelines = append(retrievedPipelines, &pipeline)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Failed to list pipelines from BadgerDB: %v", err)
		}

		if len(retrievedPipelines) != len(testPipelines) {
			t.Errorf("Expected %d pipelines, got %d", len(testPipelines), len(retrievedPipelines))
		}

		for _, expectedPipeline := range testPipelines {
			found := false
			for _, retrievedPipeline := range retrievedPipelines {
				if reflect.DeepEqual(expectedPipeline, retrievedPipeline) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected pipeline %+v not found in retrieved pipelines", expectedPipeline)
			}
		}
	})
}
