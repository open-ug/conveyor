package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	craneTypes "github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ApplicationModel struct {
	Client *clientv3.Client
}

func NewApplicationModel(cli *clientv3.Client) *ApplicationModel {
	return &ApplicationModel{
		Client: cli,
	}
}

func (m *ApplicationModel) key(name string) string {
	return fmt.Sprintf("/applications/%s", name)
}

func (m *ApplicationModel) Insert(app craneTypes.Application) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(app.Name)

	// Check existence
	getResp, err := m.Client.Get(ctx, key)
	if err != nil {
		return err
	}
	if len(getResp.Kvs) > 0 {
		return fmt.Errorf("application with name %s already exists", app.Name)
	}

	// Encode to JSON
	val, err := json.Marshal(app)
	if err != nil {
		return err
	}

	_, err = m.Client.Put(ctx, key, string(val))
	return err
}

func (m *ApplicationModel) FindOne(name string) (*craneTypes.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name)

	getResp, err := m.Client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(getResp.Kvs) == 0 {
		return nil, fmt.Errorf("application %s not found", name)
	}

	var app craneTypes.Application
	if err := json.Unmarshal(getResp.Kvs[0].Value, &app); err != nil {
		return nil, err
	}
	return &app, nil
}

func (m *ApplicationModel) FindAll() ([]craneTypes.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	getResp, err := m.Client.Get(ctx, "/applications/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	apps := make([]craneTypes.Application, 0, len(getResp.Kvs))
	for _, kv := range getResp.Kvs {
		var app craneTypes.Application
		if err := json.Unmarshal(kv.Value, &app); err != nil {
			continue // skip corrupted data
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (m *ApplicationModel) Update(app craneTypes.Application) error {
	return m.Insert(app) // Overwrites the key
}

func (m *ApplicationModel) Delete(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := m.key(name)
	_, err := m.Client.Delete(ctx, key)
	return err
}
