package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ResourceModel struct {
	Client *clientv3.Client
}

func NewResourceModel(cli *clientv3.Client) *ResourceModel {
	return &ResourceModel{
		Client: cli,
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

// Delete removes a resource by its name and type.
// It returns an error if the resource does not exist.
func (m *ResourceModel) Delete(name string, resourceType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name, resourceType)

	_, err := m.Client.Delete(ctx, key)
	return err
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
	vesion, err := strconv.Atoi(currentResource.Metadata["version"])
	if err != nil {
		return types.Resource{}, fmt.Errorf("failed to parse version: %v", err)
	}

	if resource.Metadata == nil {
		resource.Metadata = make(map[string]string)
	}
	resource.Metadata["version"] = strconv.Itoa(vesion + 1) // Increment version

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

	// Marshal the driver result to JSON
	resultData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal driver result: %v", err)
	}

	// Update the resource's metadata with the driver result
	if resource.Metadata == nil {
		resource.Metadata = make(map[string]string)
	}
	resource.Metadata["driverresults."+driver] = string(resultData)

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
