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

func (m *ResourceModel) key(name string, resourceType string) string {
	return fmt.Sprintf("/resources/%s/%s", resourceType, name)
}

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

func (m *ResourceModel) Delete(name string, resourceType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name, resourceType)

	_, err := m.Client.Delete(ctx, key)
	return err
}

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
