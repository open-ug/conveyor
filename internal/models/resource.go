package models

import (
	"context"
	"fmt"
	"time"

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
	return err
}

// FindOne retrieves a single resource by its name and type.
// It returns the resource data as a byte slice or an error if not found.
func (m *ResourceModel) FindOne(name string, resourceType string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name, resourceType)

	getResp, err := m.Client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(getResp.Kvs) == 0 {
		return nil, fmt.Errorf("resource with name %s and type %s not found", name, resourceType)
	}

	return getResp.Kvs[0].Value, nil
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
func (m *ResourceModel) Update(name string, resourceType string, resource []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name, resourceType)

	// Check existence
	getResp, err := m.Client.Get(ctx, key)
	if err != nil {
		return err
	}
	if len(getResp.Kvs) == 0 {
		return fmt.Errorf("resource with name %s and type %s not found", name, resourceType)
	}

	_, err = m.Client.Put(ctx, key, string(resource))
	return err
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
