package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ResourceModel struct {
	Client *clientv3.Client
	DB     *badger.DB
}

func NewResourceModel(cli *clientv3.Client, db *badger.DB) *ResourceModel {
	return &ResourceModel{
		Client: cli,
		DB:     db,
	}
}

// key generates a unique key for a resource based on its name and type.
func (m *ResourceModel) key(name string, resourceType string) string {
	return fmt.Sprintf("/resources/%s/%s", resourceType, name)
}

// Insert adds a new resource to the etcd store.
// It returns an error if a resource with the same name and type already exists.
func (m *ResourceModel) Insert(name string, resourceType string, resource []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name, resourceType)

	// Check existence
	getResp, err := m.Client.Get(ctx, key)
	if err != nil {
		return err
	}
	if len(getResp.Kvs) > 0 {
		return fmt.Errorf("resource with name %s and type %s already exists", name, resourceType)
	}

	_, err = m.Client.Put(ctx, key, string(resource))
	if err != nil {
		return fmt.Errorf("failed to insert resource: %v", err)
	}

	_, err = m.Client.Put(ctx, key+"/1", string(resource))
	return err
}

func (m *ResourceModel) BadgerDBInsert(name string, resourceType string, resource []byte) error {
	key := []byte(m.key(name, resourceType))

	// first unmarshal to resource inorder to set the version 1
	var real_resource types.Resource

	err := json.Unmarshal(resource, &real_resource)
	if err != nil {
		return fmt.Errorf("failed to unmarshal resource: %v", err)
	}

	setResourceVersion(&real_resource, 1)

	// marshal back to json
	resourceData, err := json.Marshal(real_resource)
	if err != nil {
		return fmt.Errorf("failed to marshal resource: %v", err)
	}

	// check existence
	err = m.DB.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err == nil {
			return fmt.Errorf("resource with name %s and type %s already exists", name, resourceType)
		}
		if err != badger.ErrKeyNotFound {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// insert resource
	err = m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, resourceData)
	})
	if err != nil {
		return fmt.Errorf("failed to insert resource: %v", err)
	}

	// insert versioned resource
	versionedKey := []byte(m.key(name, resourceType) + "/1")
	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(versionedKey, resourceData)
	})
}

// FindOne retrieves a single resource by its name and type.
// It returns the resource data as a byte slice or an error if not found.
func (m *ResourceModel) FindOne(name string, resourceType string) (types.Resource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name, resourceType)

	getResp, err := m.Client.Get(ctx, key)
	if err != nil {
		return types.Resource{}, err
	}
	if len(getResp.Kvs) == 0 {
		return types.Resource{}, fmt.Errorf("resource with name %s and type %s not found", name, resourceType)
	}

	resource := types.Resource{}
	err = json.Unmarshal(getResp.Kvs[0].Value, &resource)
	if err != nil {
		return types.Resource{}, fmt.Errorf("failed to unmarshal resource: %v", err)
	}

	return resource, nil
}

func (m *ResourceModel) BadgerDBFindOne(name string, resourceType string) (types.Resource, error) {
	key := []byte(m.key(name, resourceType))

	var resource types.Resource
	err := m.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return fmt.Errorf("resource with name %s and type %s not found", name, resourceType)
			}
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &resource)
		})
	})
	if err != nil {
		return types.Resource{}, err
	}

	return resource, nil
}

// Delete removes a resource by its name and type.
// It returns an error if the resource does not exist.
func (m *ResourceModel) Delete(name string, resourceType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name, resourceType)

	_, err := m.Client.Delete(ctx, key)
	return err
}

func (m *ResourceModel) BadgerDBDelete(name string, resourceType string) error {
	key := []byte(m.key(name, resourceType))

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// List retrieves all resources of a specific type.
// It returns a slice of resource names or an error if the operation fails.
func (m *ResourceModel) List(resourceType string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := fmt.Sprintf("/resources/%s/", resourceType)

	getResp, err := m.Client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var resources []string
	for _, kv := range getResp.Kvs {
		resources = append(resources, string(kv.Key))
	}
	return resources, nil
}

func (m *ResourceModel) BadgerDBList(resourceType string) ([]types.Resource, error) {
	prefix := []byte(fmt.Sprintf("/resources/%s/", resourceType))

	var resources []types.Resource

	err := m.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			key := string(item.Key())

			// Skip versioned resources like:
			// /resources/pipe4/pipeline-1/1
			remaining := strings.TrimPrefix(key, string(prefix))
			if strings.Contains(remaining, "/") {
				continue
			}

			err := item.Value(func(val []byte) error {
				var res types.Resource

				if err := json.Unmarshal(val, &res); err != nil {
					return fmt.Errorf("failed to unmarshal resource: %v", err)
				}

				resources = append(resources, res)
				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %v", err)
	}

	return resources, nil
}

// Update modifies an existing resource's data.
// It returns an error if the resource does not exist.
func (m *ResourceModel) Update(name string, resourceType string, resource types.Resource) (types.Resource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name, resourceType)

	// Check existence
	currentResource, err := m.FindOne(name, resourceType)
	if err != nil {
		return types.Resource{}, fmt.Errorf("resource with name %s and type %s not found: %v ", name, resourceType, err)
	}

	resource.ID = currentResource.ID // Ensure the ID remains unchanged
	version, err := getResourceVersion(currentResource)
	if err != nil {
		return types.Resource{}, fmt.Errorf("failed to get current resource version: %v", err)
	}

	setResourceVersion(&currentResource, version+1)

	resource.Metadata = currentResource.Metadata // Preserve existing metadata and versioning

	// Marshal the updated resource to JSON
	resourceData, err := json.Marshal(resource)
	if err != nil {
		return types.Resource{}, fmt.Errorf("failed to marshal resource: %v", err)
	}

	_, err = m.Client.Put(ctx, key, string(resourceData))
	if err != nil {
		return types.Resource{}, fmt.Errorf("failed to update resource: %v", err)
	}

	// save versioned resource
	versionedKey := fmt.Sprintf("%s/%s", key, resource.Metadata["version"])
	_, err = m.Client.Put(ctx, versionedKey, string(resourceData))
	if err != nil {
		return types.Resource{}, fmt.Errorf("failed to save versioned resource: %v", err)
	}
	return resource, err
}

func getResourceVersion(resource types.Resource) (int, error) {
	if resource.Metadata == nil {
		return 0, fmt.Errorf("resource metadata is nil")
	}

	versionStr, ok := resource.Metadata["version"].(string)
	if !ok {
		return 0, fmt.Errorf("version not found in resource metadata or is not a string")
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse version: %v", err)
	}

	return version, nil
}

func setResourceVersion(resource *types.Resource, version int) {
	if resource.Metadata == nil {
		resource.Metadata = make(map[string]interface{})
	}
	resource.Metadata["version"] = strconv.Itoa(version)
}

func (m *ResourceModel) BadgerDBUpdate(name string, resourceType string, resource types.Resource) (types.Resource, error) {
	key := []byte(m.key(name, resourceType))

	// Check existence
	currentResource, err := m.BadgerDBFindOne(name, resourceType)
	if err != nil {
		return types.Resource{}, fmt.Errorf("resource with name %s and type %s not found: %v ", name, resourceType, err)
	}

	resource.ID = currentResource.ID // Ensure the ID remains unchanged
	version, err := getResourceVersion(currentResource)
	if err != nil {
		return types.Resource{}, fmt.Errorf("failed to get current resource version: %v", err)
	}

	setResourceVersion(&currentResource, version+1)

	resource.Metadata = currentResource.Metadata // Preserve existing metadata and versioning

	// Marshal the updated resource to JSON
	resourceData, err := json.Marshal(resource)
	if err != nil {
		return types.Resource{}, fmt.Errorf("failed to marshal resource: %v", err)
	}

	err = m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, resourceData)
	})
	if err != nil {
		return types.Resource{}, fmt.Errorf("failed to update resource: %v", err)
	}

	// save versioned resource
	versionedKey := []byte(fmt.Sprintf("%s/%s", key, resource.Metadata["version"]))
	err = m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(versionedKey, resourceData)
	})
	if err != nil {
		return types.Resource{}, fmt.Errorf("failed to save versioned resource: %v", err)
	}

	return resource, nil
}

// FindAll retrieves all resources of a specific type.
// It returns a slice of resource names or an error if no resources are found.
func (m *ResourceModel) FindAll(resourceType string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	getResp, err := m.Client.Get(ctx, fmt.Sprintf("/resources/%s/", resourceType), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if len(getResp.Kvs) == 0 {
		return nil, fmt.Errorf("no resources of type %s found", resourceType)
	}

	var resources []string
	for _, kv := range getResp.Kvs {
		resources = append(resources, string(kv.Value))
	}
	return resources, nil
}

func (m *ResourceModel) BadgerDBFindAll(resourceType string) ([]types.Resource, error) {
	prefix := []byte(fmt.Sprintf("/resources/%s/", resourceType))

	var resources []types.Resource
	err := m.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			key := string(item.Key())

			// Skip versioned resources like:
			// /resources/pipe4/pipeline-1/1
			remaining := strings.TrimPrefix(key, string(prefix))
			if strings.Contains(remaining, "/") {
				continue
			}

			err := item.Value(func(val []byte) error {
				var res types.Resource

				if err := json.Unmarshal(val, &res); err != nil {
					return fmt.Errorf("failed to unmarshal resource: %v", err)
				}

				resources = append(resources, res)
				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find resources: %v", err)
	}

	if len(resources) == 0 {
		return nil, fmt.Errorf("no resources of type %s found", resourceType)
	}

	return resources, nil
}

// FindByVersion retrieves a specific version of a resource by its name and type and version.
// It returns the resource data or an error if not found.
func (m *ResourceModel) FindByVersion(name string, resourceType string, version string) (types.Resource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := fmt.Sprintf("/resources/%s/%s/%s", resourceType, name, version)

	getResp, err := m.Client.Get(ctx, key)
	if err != nil {
		return types.Resource{}, err
	}
	if len(getResp.Kvs) == 0 {
		return types.Resource{}, fmt.Errorf("resource with name %s, type %s and version %s not found", name, resourceType, version)
	}

	resource := types.Resource{}
	err = json.Unmarshal(getResp.Kvs[0].Value, &resource)
	if err != nil {
		return types.Resource{}, fmt.Errorf("failed to unmarshal resource: %v", err)
	}

	return resource, nil
}

/// A function that saves the driver result. This data is then stored in the metadata.driverresults.[driver] field of the resource and is arbitrary data types

func (m *ResourceModel) SaveDriverResult(name string, resourceType string, driver string, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Retrieve the current resource
	resource, err := m.FindOne(name, resourceType)
	if err != nil {
		return fmt.Errorf("failed to find resource: %v", err)
	}

	// Update the resource's metadata with the driver result
	err = updateDriverResult(&resource, driver, result)
	if err != nil {
		return fmt.Errorf("failed to update driver result: %v", err)
	}

	// Marshal the updated resource to JSON
	resourceData, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("failed to marshal updated resource: %v", err)
	}

	// Save the updated resource back to etcd
	key := m.key(name, resourceType)
	_, err = m.Client.Put(ctx, key, string(resourceData))
	if err != nil {
		return fmt.Errorf("failed to save updated resource: %v", err)
	}

	return nil
}

func (m *ResourceModel) SaveDriverResultBadgerDB(name string, resourceType string, driver string, result interface{}) error {
	// Retrieve the current resource
	resource, err := m.BadgerDBFindOne(name, resourceType)
	if err != nil {
		return fmt.Errorf("failed to find resource: %v", err)
	}

	// Update the resource's metadata with the driver result
	err = updateDriverResult(&resource, driver, result)
	if err != nil {
		return fmt.Errorf("failed to update driver result: %v", err)
	}

	// Marshal the updated resource to JSON
	resourceData, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("failed to marshal updated resource: %v", err)
	}

	// Save the updated resource back to BadgerDB
	key := []byte(m.key(name, resourceType))
	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, resourceData)
	})
}

func updateDriverResult(resource *types.Resource, driver string, result interface{}) error {
	// generate metadata object for driver results
	driverResults := make(map[string]interface{})
	if resource.Metadata != nil {
		if existingResults, ok := resource.Metadata["driverresults"]; ok {
			if existingResultsMap, ok := existingResults.(map[string]interface{}); ok {
				driverResults = existingResultsMap
			}
		}
	}

	// Update the driver result
	driverResults[driver] = result

	// set the updated driver results back to resource metadata
	if resource.Metadata == nil {
		resource.Metadata = make(map[string]interface{})
	}

	resource.Metadata["driverresults"] = driverResults

	return nil
}
