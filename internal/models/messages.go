package models

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	craneTypes "github.com/open-ug/conveyor/pkg/types"
)

type DriverMessageModel struct {
	Prefix string // e.g., "driver-messages/"
	DB     *badger.DB
}

func NewDriverMessageModel(db *badger.DB) *DriverMessageModel {
	return &DriverMessageModel{
		Prefix: "driver-messages/",
		DB:     db,
	}
}

func (m *DriverMessageModel) Insert(message craneTypes.DriverMessage) error {
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

func (m *DriverMessageModel) DeleteOne(id string) error {
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
