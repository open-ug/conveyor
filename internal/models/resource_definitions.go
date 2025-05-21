package models

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

/*## Resource Definitions

By default, Conveyor CI does not come with any Resource in its API. You have to install them into Conveyor CI. To install a Resource, you must define a Resource Definetion. A Resource defenition is method of defining the Schema of a Resource. It determines how the Resouce will be defined, what properties it will have and also the validation schema for the resource.

Resource definitions are created by following the [Open API Specification](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md). They can be defined either in Yaml or JSON syntax and saved into the Conveyor CI API. Resource Definitions have only two mandatory requirementsa and these are:

*/

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
