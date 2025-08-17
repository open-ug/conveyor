package models

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ResourceDefinitionModel struct {
	Client *clientv3.Client
}

func NewResourceDefinitionModel(cli *clientv3.Client) *ResourceDefinitionModel {
	return &ResourceDefinitionModel{
		Client: cli,
	}
}

func (m *ResourceDefinitionModel) key(name string) string {
	return fmt.Sprintf("/resource_definitions/%s", name)
}

// Insert adds a new resource definition to the etcd store.
// It returns an error if a resource definition with the same name already exists.
func (m *ResourceDefinitionModel) Insert(name string, definition []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name)

	// Check existence
	getResp, err := m.Client.Get(ctx, key)
	if err != nil {
		return err
	}
	if len(getResp.Kvs) > 0 {
		return fmt.Errorf("resource definition with name %s already exists", name)
	}

	_, err = m.Client.Put(ctx, key, string(definition))
	return err
}

// FindOne retrieves a single resource definition by its name.
// It returns the resource definition data as a byte slice or an error if not found.
func (m *ResourceDefinitionModel) FindOne(name string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name)

	getResp, err := m.Client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(getResp.Kvs) == 0 {
		return nil, fmt.Errorf("resource definition with name %s not found", name)
	}

	return getResp.Kvs[0].Value, nil
}

// Delete removes a resource definition by its name.
// It returns an error if the resource definition does not exist.
func (m *ResourceDefinitionModel) Delete(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name)

	// Check existence
	getResp, err := m.Client.Get(ctx, key)
	if err != nil {
		return err
	}
	if len(getResp.Kvs) == 0 {
		return fmt.Errorf("resource definition with name %s not found", name)
	}

	_, err = m.Client.Delete(ctx, key)
	return err
}

// FindAll retrieves all resource definitions.
// It returns a slice of resource definition names or an error if no definitions are found.
func (m *ResourceDefinitionModel) FindAll() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	getResp, err := m.Client.Get(ctx, "/resource_definitions/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if len(getResp.Kvs) == 0 {
		return nil, fmt.Errorf("no resource definitions found")
	}

	var resourceDefinitions []string
	for _, kv := range getResp.Kvs {
		resourceDefinitions = append(resourceDefinitions, string(kv.Value))
	}
	return resourceDefinitions, nil
}

// Update modifies an existing resource definition's data.
// It returns an error if the resource definition does not exist.
func (m *ResourceDefinitionModel) Update(name string, definition []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name)

	// Check existence
	getResp, err := m.Client.Get(ctx, key)
	if err != nil {
		return err
	}
	if len(getResp.Kvs) == 0 {
		return fmt.Errorf("resource definition with name %s not found", name)
	}

	_, err = m.Client.Put(ctx, key, string(definition))
	return err
}
