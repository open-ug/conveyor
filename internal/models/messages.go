package models

import (
	"context"
	"encoding/json"
	"fmt"

	craneTypes "github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type DriverMessageModel struct {
	Client *clientv3.Client
	Prefix string // e.g., "driver-messages/"
}

func NewDriverMessageModel(cli *clientv3.Client) *DriverMessageModel {
	return &DriverMessageModel{
		Client: cli,
		Prefix: "driver-messages/",
	}
}

func (m *DriverMessageModel) Insert(message craneTypes.DriverMessage) error {
	key := m.Prefix + message.ID // assuming ID is a unique string field
	value, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %v", err)
	}

	_, err = m.Client.Put(context.Background(), key, string(value))
	if err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}
	return nil
}

func (m *DriverMessageModel) FindOne(id string) (*craneTypes.DriverMessage, error) {
	key := m.Prefix + id
	resp, err := m.Client.Get(context.Background(), key)
	if err != nil || len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("message not found: %v", err)
	}

	var msg craneTypes.DriverMessage
	if err := json.Unmarshal(resp.Kvs[0].Value, &msg); err != nil {
		return nil, fmt.Errorf("failed to decode message: %v", err)
	}
	return &msg, nil
}

func (m *DriverMessageModel) FindAll() ([]craneTypes.DriverMessage, error) {
	resp, err := m.Client.Get(context.Background(), m.Prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %v", err)
	}

	var messages []craneTypes.DriverMessage
	for _, kv := range resp.Kvs {
		var msg craneTypes.DriverMessage
		if err := json.Unmarshal(kv.Value, &msg); err == nil {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

func (m *DriverMessageModel) UpdateOne(id string, updated craneTypes.DriverMessage) error {
	return m.Insert(updated) // etcd has no partial update; replace the value
}

func (m *DriverMessageModel) DeleteOne(id string) error {
	key := m.Prefix + id
	_, err := m.Client.Delete(context.Background(), key)
	if err != nil {
		return fmt.Errorf("failed to delete message: %v", err)
	}
	return nil
}
