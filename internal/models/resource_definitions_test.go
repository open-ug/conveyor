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

func TestInsertResourceDefinition(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)
	defer db.Close()

	// Create a new ResourceDefinitionModel instance
	rm := models.NewResourceDefinitionModel(nil, db)

	// Define a test resource definition
	testResourceDef := &types.ResourceDefinition{
		Name: "test-resource",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"field1": map[string]interface{}{
					"type": "string",
				},
				"field2": map[string]interface{}{
					"type": "integer",
				},
			},
			"required": []string{"field1", "field2"},
		},
	}

	t.Run("insert resource definition", func(t *testing.T) {
		err := rm.BadgerInsert(testResourceDef.Name, testResourceDef)
		if err != nil {
			t.Fatalf("Failed to insert resource definition: %v", err)
		}

		// Verify that the resource definition was inserted in BadgerDB
		err = db.View(func(txn *badger.Txn) error {
			key := fmt.Sprintf("/resource_definitions/%s", testResourceDef.Name)
			item, err := txn.Get([]byte(key))
			if err != nil {
				return err
			}

			var storedResourceDef types.ResourceDefinition
			err = item.Value(func(val []byte) error {
				return json.Unmarshal(val, &storedResourceDef)
			})
			if err != nil {
				return err
			}

			// manually compare the schema maps since reflect.DeepEqual may not work as expected for maps with interface{} values
			if storedResourceDef.Name != testResourceDef.Name {
				t.Errorf("Stored resource definition name does not match expected. Got %s, want %s", storedResourceDef.Name, testResourceDef.Name)
			}

			storedShema := storedResourceDef.Schema.(map[string]interface{})
			expectedSchema := testResourceDef.Schema.(map[string]interface{})

			if storedShema["type"] != expectedSchema["type"] {
				t.Errorf("Stored resource definition schema type does not match expected. Got %s, want %s", storedShema["type"], expectedSchema["type"])
			}

			if storedProps, ok := storedShema["properties"].(map[string]interface{}); ok {
				if expectedProps, ok := expectedSchema["properties"].(map[string]interface{}); ok {
					for key, expectedProp := range expectedProps {
						if storedProp, ok := storedProps[key]; ok {
							if !reflect.DeepEqual(storedProp, expectedProp) {
								t.Errorf("Stored resource definition property %s does not match expected. Got %+v, want %+v", key, storedProp, expectedProp)
							}
						} else {
							t.Errorf("Stored resource definition is missing property %s", key)
						}
					}
				} else {
					t.Errorf("Expected schema properties is not a map")
				}
			} else {
				t.Errorf("Stored schema properties is not a map")
			}

			return nil
		})
		if err != nil {
			t.Fatalf("Failed to verify resource definition in BadgerDB: %v", err)
		}
	})
}

func TestBadgerFindOneResourceDefinition(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)
	defer db.Close()

	// Create a new ResourceDefinitionModel instance
	rm := models.NewResourceDefinitionModel(nil, db)

	// Define a test resource definition
	testResourceDef := &types.ResourceDefinition{
		Name: "test-resource",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"field1": map[string]interface{}{
					"type": "string",
				},
				"field2": map[string]interface{}{
					"type": "integer",
				},
			},
			"required": []string{"field1", "field2"},
		},
	}

	t.Run("find one resource definition", func(t *testing.T) {
		// Insert the test resource definition into BadgerDB
		err := rm.BadgerInsert(testResourceDef.Name, testResourceDef)
		if err != nil {
			t.Fatalf("Failed to insert resource definition: %v", err)
		}

		// Retrieve the resource definition from BadgerDB
		retrievedResourceDef, err := rm.BadgerFindOne(testResourceDef.Name)
		if err != nil {
			t.Fatalf("Failed to find resource definition: %v", err)
		}

		// manually compare the schema maps since reflect.DeepEqual may not work as expected for maps with interface{} values
		if retrievedResourceDef.Name != testResourceDef.Name {
			t.Errorf("Retrieved resource definition name does not match expected. Got %s, want %s", retrievedResourceDef.Name, testResourceDef.Name)
		}

		retrievedSchema := retrievedResourceDef.Schema.(map[string]interface{})
		expectedSchema := testResourceDef.Schema.(map[string]interface{})

		if retrievedSchema["type"] != expectedSchema["type"] {
			t.Errorf("Retrieved resource definition schema type does not match expected. Got %s, want %s", retrievedSchema["type"], expectedSchema["type"])
		}

		if retrievedProps, ok := retrievedSchema["properties"].(map[string]interface{}); ok {
			if expectedProps, ok := expectedSchema["properties"].(map[string]interface{}); ok {
				for key, expectedProp := range expectedProps {
					if retrievedProp, ok := retrievedProps[key]; ok {
						if !reflect.DeepEqual(retrievedProp, expectedProp) {
							t.Errorf("Retrieved resource definition property %s does not match expected. Got %+v, want %+v", key, retrievedProp, expectedProp)
						}
					} else {
						t.Errorf("Retrieved resource definition is missing property %s", key)
					}
				}
			} else {
				t.Errorf("Expected schema properties is not a map")
			}
		} else {
			t.Errorf("Retrieved schema properties is not a map")
		}
	})
}

func TestBadgerUpdateResourceDefinition(t *testing.T) {
	// Create a mock etcd client and badger DB for testing
	db := setupTestDB(t)
	defer db.Close()

	// Create a new ResourceDefinitionModel instance
	rm := models.NewResourceDefinitionModel(nil, db)

	// Define a test resource definition
	testResourceDef := &types.ResourceDefinition{
		Name: "test-resource",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"field1": map[string]interface{}{
					"type": "string",
				},
				"field2": map[string]interface{}{
					"type": "integer",
				},
			},
			"required": []string{"field1", "field2"},
		},
	}

	t.Run("update resource definition", func(t *testing.T) {
		// Insert the test resource definition into BadgerDB
		err := rm.BadgerInsert(testResourceDef.Name, testResourceDef)
		if err != nil {
			t.Fatalf("Failed to insert resource definition: %v", err)
		}

		// Update the resource definition's schema, first make a copy of the original schema to avoid modifying the original testResourceDef
		updatedSchema := map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"field1": map[string]interface{}{
					"type": "string",
				},
				"field2": map[string]interface{}{
					"type": "integer",
				},
				"field3": map[string]interface{}{
					"type": "boolean",
				},
			},
			"required": []string{"field1", "field2", "field3"},
		}

		testResourceDef.Schema = updatedSchema

		err = rm.BadgerUpdate(testResourceDef.Name, testResourceDef)
		if err != nil {
			t.Fatalf("Failed to update resource definition: %v", err)
		}

		var updatedResourceDef types.ResourceDefinition

		err = db.View(func(txn *badger.Txn) error {
			key := fmt.Sprintf("/resource_definitions/%s", testResourceDef.Name)
			item, err := txn.Get([]byte(key))
			if err != nil {
				return err
			}

			return item.Value(func(val []byte) error {
				return json.Unmarshal(val, &updatedResourceDef)
			})
		})
		if err != nil {
			t.Fatalf("Failed to retrieve updated resource definition from BadgerDB: %v", err)
		}

		// manually compare the schema maps since reflect.DeepEqual may not work as expected for maps with interface{} values
		if updatedResourceDef.Name != testResourceDef.Name {
			t.Errorf("Updated resource definition name does not match expected. Got %s, want %s", updatedResourceDef.Name, testResourceDef.Name)
		}

		updatedSchemaRetrieved := updatedResourceDef.Schema.(map[string]interface{})
		expectedSchema := testResourceDef.Schema.(map[string]interface{})

		if updatedSchemaRetrieved["type"] != expectedSchema["type"] {
			t.Errorf("Updated resource definition schema type does not match expected. Got %s, want %s", updatedSchemaRetrieved["type"], expectedSchema["type"])
		}

		if updatedProps, ok := updatedSchemaRetrieved["properties"].(map[string]interface{}); ok {
			if expectedProps, ok := expectedSchema["properties"].(map[string]interface{}); ok {
				for key, expectedProp := range expectedProps {
					if updatedProp, ok := updatedProps[key]; ok {
						if !reflect.DeepEqual(updatedProp, expectedProp) {
							t.Errorf("Updated resource definition property %s does not match expected. Got %+v, want %+v", key, updatedProp, expectedProp)
						}
					} else {
						t.Errorf("Updated resource definition is missing property %s", key)
					}
				}
			} else {
				t.Errorf("Expected schema properties is not a map")
			}
		} else {
			t.Errorf("Updated schema properties is not a map")
		}
	})
}
