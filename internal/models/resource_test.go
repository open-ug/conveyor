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

	resource := types.Resource{
		Name:     "pipeline-1",
		Resource: "pipe4",
		Spec: map[string]interface{}{
			"image": "ubuntu:latest",
			"steps": []map[string]string{
				{
					"name":    "list dir",
					"command": "ls",
				},
			},
		},
	}

	marshalResource, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal resource: %v", err)
	}

	resourcemodel := models.NewResourceModel(nil, db)

	t.Run("Insert Resource", func(t *testing.T) {
		err := resourcemodel.BadgerDBInsert(resource.Name, resource.Resource, marshalResource)
		if err != nil {
			t.Fatalf("failed to insert resource: %v", err)
		}

		// Verify the resource was inserted
		err = db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("/resources/pipe4/pipeline-1"))
			if err != nil {
				return err
			}
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			// we cannot directly compare the marshaled resource with the value from badger because the Insert function adds some metadata fields, we should only compare resource name and type for this test
			var insertedResource types.Resource
			err = json.Unmarshal(val, &insertedResource)
			if err != nil {
				return err
			}
			if insertedResource.Name != resource.Name {
				t.Errorf("expected resource name '%s', got '%s'", resource.Name, insertedResource.Name)
			}
			if insertedResource.Resource != resource.Resource {
				t.Errorf("expected resource type '%s', got '%s'", resource.Resource, insertedResource.Resource)
			}

			// check if spec is not empty
			if insertedResource.Spec == nil {
				t.Errorf("expected resource spec to be not nil")
			}

			// we can also check if the spec fields are correct
			specMap := insertedResource.Spec.(map[string]interface{})
			if specMap["image"] != "ubuntu:latest" {
				t.Errorf("expected image 'ubuntu:latest', got '%s'", specMap["image"])
			}
			steps := specMap["steps"].([]interface{})
			if len(steps) != 1 {
				t.Errorf("expected 1 step, got %d", len(steps))
			} else {
				step := steps[0].(map[string]interface{})
				if step["name"] != "list dir" {
					t.Errorf("expected step name 'list dir', got '%s'", step["name"])
				}
				if step["command"] != "ls" {
					t.Errorf("expected step command 'ls', got '%s'", step["command"])
				}
			}

			return nil
		})
		if err != nil {
			t.Fatalf("failed to verify resource insertion: %v", err)
		}
	})

	// verify versioned resource was inserted
	t.Run("Verify Versioned Resource", func(t *testing.T) {
		err = db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("/resources/pipe4/pipeline-1/1"))
			if err != nil {
				return err
			}
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			var versionedResource types.Resource
			err = json.Unmarshal(val, &versionedResource)
			if err != nil {
				return err
			}
			if versionedResource.Name != resource.Name {
				t.Errorf("expected versioned resource name '%s', got '%s'", resource.Name, versionedResource.Name)
			}
			if versionedResource.Resource != resource.Resource {
				t.Errorf("expected versioned resource type '%s', got '%s'", resource.Resource, versionedResource.Resource)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("failed to verify versioned resource insertion: %v", err)
		}
	})

}

type PipelineResource struct {
	Name     string               `json:"name"`
	Resource string               `json:"resource"`
	Spec     PipelineResourceSpec `json:"spec"`
}

type PipelineResourceSpec struct {
	Image string         `json:"image"`
	Steps []PipelineStep `json:"steps"`
}

type PipelineStep struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

var resource = PipelineResource{
	Name:     "pipeline-1",
	Resource: "pipe4",
	Spec: PipelineResourceSpec{
		Image: "ubuntu:latest",
		Steps: []PipelineStep{
			{
				Name:    "list dir",
				Command: "ls",
			},
		},
	},
}

func Test_Resource_FindOne(t *testing.T) {
	db := setupTestDB(t)

	resourcemodel := models.NewResourceModel(nil, db)

	rawResource, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal resource: %v", err)
	}

	// Insert a resource to find
	err = resourcemodel.BadgerDBInsert(resource.Name, resource.Resource, rawResource)
	if err != nil {
		t.Fatalf("failed to insert resource: %v", err)
	}

	t.Run("Find Existing Resource", func(t *testing.T) {
		foundResource, err := resourcemodel.BadgerDBFindOne(resource.Name, resource.Resource)
		if err != nil {
			t.Fatalf("failed to find resource: %v", err)
		}

		converted := PipelineResource{
			Name:     foundResource.Name,
			Resource: foundResource.Resource,
		}

		// Convert Spec map -> struct
		specBytes, err := json.Marshal(foundResource.Spec)
		if err != nil {
			t.Fatalf("failed to marshal spec: %v", err)
		}

		err = json.Unmarshal(specBytes, &converted.Spec)
		if err != nil {
			t.Fatalf("failed to unmarshal spec: %v", err)
		}

		if !reflect.DeepEqual(converted, resource) {
			t.Errorf("expected %+v, got %+v", resource, converted)
		}

	})

	t.Run("Find Non-Existing Resource", func(t *testing.T) {
		_, err := resourcemodel.BadgerDBFindOne("non-existent", "test-type")
		if err == nil {
			t.Fatalf("expected error when finding non-existent resource, got nil")
		}
	})
}

func Test_Resource_Delete(t *testing.T) {
	db := setupTestDB(t)

	resourcemodel := models.NewResourceModel(nil, db)

	resource := types.Resource{
		Name:     "test-resource",
		Resource: "test-type",
		Spec:     map[string]interface{}{"key": "value"},
	}

	marshalResource, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal resource: %v", err)
	}

	// Insert a resource to delete
	err = resourcemodel.BadgerDBInsert(resource.Name, resource.Resource, marshalResource)
	if err != nil {
		t.Fatalf("failed to insert resource: %v", err)
	}

	t.Run("Delete Existing Resource", func(t *testing.T) {
		err := resourcemodel.BadgerDBDelete("test-resource", "test-type")
		if err != nil {
			t.Fatalf("failed to delete resource: %v", err)
		}

		// Verify the resource was deleted
		err = db.View(func(txn *badger.Txn) error {
			_, err := txn.Get([]byte("/resources/test-type/test-resource"))
			if err == nil {
				return fmt.Errorf("expected resource to be deleted, but it still exists")
			}
			if err != badger.ErrKeyNotFound {
				return err
			}
			return nil
		})
		if err != nil {
			t.Fatalf("failed to verify resource deletion: %v", err)
		}
	})

}

// Test BadgerList
func Test_Resource_List(t *testing.T) {
	db := setupTestDB(t)

	resourcemodel := models.NewResourceModel(nil, db)

	// Insert multiple resources of PipelineResource type
	for i := 1; i <= 3; i++ {
		resource := PipelineResource{
			Name:     fmt.Sprintf("pipeline-%d", i),
			Resource: "pipe4",
			Spec: PipelineResourceSpec{
				Image: "ubuntu:latest",
				Steps: []PipelineStep{
					{
						Name:    "list dir",
						Command: "ls",
					},
				},
			},
		}

		rawResource, err := json.Marshal(resource)
		if err != nil {
			t.Fatalf("failed to marshal resource: %v", err)
		}

		err = resourcemodel.BadgerDBInsert(resource.Name, resource.Resource, rawResource)
		if err != nil {
			t.Fatalf("failed to insert resource: %v", err)
		}
	}

	t.Run("List Resources", func(t *testing.T) {
		resources, err := resourcemodel.BadgerDBList("pipe4")
		fmt.Printf("Listed Resources: %+v\n", resources)
		if err != nil {
			t.Fatalf("failed to list resources: %v", err)
		}

		if len(resources) != 3 {
			t.Errorf("expected 3 resources, got %d", len(resources))
		}

		for i, res := range resources {
			expectedName := fmt.Sprintf("pipeline-%d", i+1)
			if res.Name != expectedName {
				t.Errorf("expected resource name '%s', got '%s'", expectedName, res.Name)
			}
			if res.Resource != "pipe4" {
				t.Errorf("expected resource type 'pipe4', got '%s'", res.Resource)
			}
		}
	})

}

func Test_Resource_FindAll(t *testing.T) {
	db := setupTestDB(t)

	resourcemodel := models.NewResourceModel(nil, db)

	// Insert multiple resources of PipelineResource type
	for i := 1; i <= 3; i++ {
		resource := PipelineResource{
			Name:     fmt.Sprintf("pipeline-%d", i),
			Resource: "pipe4",
			Spec: PipelineResourceSpec{
				Image: "ubuntu:latest",
				Steps: []PipelineStep{
					{
						Name:    "list dir",
						Command: "ls",
					},
				},
			},
		}

		rawResource, err := json.Marshal(resource)
		if err != nil {
			t.Fatalf("failed to marshal resource: %v", err)
		}

		err = resourcemodel.BadgerDBInsert(resource.Name, resource.Resource, rawResource)
		if err != nil {
			t.Fatalf("failed to insert resource: %v", err)
		}
	}

	t.Run("Find All Resources", func(t *testing.T) {
		resources, err := resourcemodel.BadgerDBFindAll("pipe4")
		fmt.Printf("Found Resources: %+v\n", resources)
		if err != nil {
			t.Fatalf("failed to find all resources: %v", err)
		}

		if len(resources) != 3 {
			t.Errorf("expected 3 resources, got %d", len(resources))
		}

		for i, res := range resources {
			expectedName := fmt.Sprintf("pipeline-%d", i+1)
			if res.Name != expectedName {
				t.Errorf("expected resource name '%s', got '%s'", expectedName, res.Name)
			}
			if res.Resource != "pipe4" {
				t.Errorf("expected resource type 'pipe4', got '%s'", res.Resource)
			}
		}
	})

}

func Test_Resource_Update(t *testing.T) {
	db := setupTestDB(t)

	resourcemodel := models.NewResourceModel(nil, db)

	// Sample resource to insert and update
	resource := types.Resource{
		Name:     "pipeline-1",
		Resource: "pipe4",
		Spec: map[string]interface{}{
			"image": "ubuntu:latest",
			"steps": []map[string]string{
				{
					"name":    "list dir",
					"command": "ls",
				},
			},
		},
	}

	rawResource, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal resource: %v", err)
	}

	// Insert the resource to update
	err = resourcemodel.BadgerDBInsert(resource.Name, resource.Resource, rawResource)
	if err != nil {
		t.Fatalf("failed to insert resource: %v", err)
	}

	t.Run("Update Existing Resource", func(t *testing.T) {
		// Update the resource spec
		resource.Spec = map[string]interface{}{
			"image": "alpine:latest",
			"steps": []map[string]string{
				{
					"name":    "print working dir",
					"command": "pwd",
				},
			},
		}

		updatedResource, err := resourcemodel.BadgerDBUpdate(resource.Name, resource.Resource, resource)
		if err != nil {
			t.Fatalf("failed to update resource: %v", err)
		}

		// Verify the resource was updated
		if updatedResource.Spec.(map[string]interface{})["image"] != "alpine:latest" {
			t.Errorf("expected image 'alpine:latest', got '%s'", updatedResource.Spec.(map[string]interface{})["image"])
		}
	})
}

func Test_SaveDriverResult(t *testing.T) {
	db := setupTestDB(t)

	resourcemodel := models.NewResourceModel(nil, db)

	// Sample resource result to save
	resource := types.Resource{
		Name:     "pipeline-1",
		Resource: "pipe4",
		Spec: map[string]interface{}{
			"name": "success",
		},
	}

	// "{\"success\":false,\"message\":\"Unsupported event type: process\",\"data\":null,\"driver\":\"docker\"}",
	result := map[string]interface{}{
		"success": false,
		"message": "Unsupported event type: process",
		"data":    nil,
		"driver":  "docker",
	}

	jsonResource, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal resource: %v", err)
	}

	// insert resource
	err = resourcemodel.BadgerDBInsert(resource.Name, resource.Resource, jsonResource)
	if err != nil {
		t.Fatalf("failed to insert resource: %v", err)
	}

	t.Run("Save Driver Result", func(t *testing.T) {
		err := resourcemodel.SaveDriverResultBadgerDB("pipeline-1", "pipe4", "docker", result)
		if err != nil {
			t.Fatalf("failed to save driver result: %v", err)
		}

		// find the new updated resource
		updatedResource, err := resourcemodel.BadgerDBFindOne("pipeline-1", "pipe4")
		if err != nil {
			t.Fatalf("failed to find updated resource: %v", err)
		}

		// check if the driver result was saved in the resource metadata
		if updatedResource.Metadata["driverresults"] == nil {
			t.Fatalf("expected driver result in resource metadata, got nil")
		}

		// check the driver result values
		driverResults := updatedResource.Metadata["driverresults"].(map[string]interface{})
		if driverResults["docker"] == nil {
			t.Fatalf("expected docker driver result in resource metadata, got nil")
		}

		dockerResult := driverResults["docker"].(map[string]interface{})
		if dockerResult["success"] != false {
			t.Errorf("expected success 'false', got '%v'", dockerResult["success"])
		}
		if dockerResult["message"] != "Unsupported event type: process" {
			t.Errorf("expected message 'Unsupported event type: process', got '%v'", dockerResult["message"])
		}
	})
}
