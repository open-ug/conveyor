package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	craneTypes "github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type DriverMessageModel struct {
	Client *clientv3.Client
	Prefix string // e.g., "driver-messages/"
	DB     *badger.DB
}

func NewDriverMessageModel(cli *clientv3.Client, db *badger.DB) *DriverMessageModel {
	return &DriverMessageModel{
		Client: cli,
		Prefix: "driver-messages/",
		DB:     db,
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

func (m *DriverMessageModel) BadgerInsert(message craneTypes.DriverMessage) error {
	err := m.DB.Update(func(txn *badger.Txn) error {
		key := m.Prefix + message.ID
		value, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to serialize message: %v", err)
		}

		err = txn.Set([]byte(key), value)
		if err != nil {
			return fmt.Errorf("failed to insert message into BadgerDB: %v", err)
		}
		return nil
	})
	return err
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

func (m *DriverMessageModel) BadgerFindOne(id string) (*craneTypes.DriverMessage, error) {
	var msg craneTypes.DriverMessage
	err := m.DB.View(func(txn *badger.Txn) error {
		key := m.Prefix + id
		item, err := txn.Get([]byte(key))
		if err != nil {
			return fmt.Errorf("message not found: %v", err)
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &msg)
		})
	})
	if err != nil {
		return nil, err
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

func (m *DriverMessageModel) BadgerFindAll() ([]craneTypes.DriverMessage, error) {
	var messages []craneTypes.DriverMessage
	err := m.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(m.Prefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var msg craneTypes.DriverMessage
				if err := json.Unmarshal(val, &msg); err != nil {
					return fmt.Errorf("failed to decode message: %v", err)
				}
				messages = append(messages, msg)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
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

func (m *DriverMessageModel) BadgerDeleteOne(id string) error {
	err := m.DB.Update(func(txn *badger.Txn) error {
		key := m.Prefix + id
		err := txn.Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("failed to delete message from BadgerDB: %v", err)
		}
		return nil
	})
	return err
}
