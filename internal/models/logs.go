package models

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"github.com/open-ug/conveyor/pkg/types"
)

type LogModel struct {
	DB *badger.DB
}

// Generate composite key
func makeKey(pipeline, driver, runid, timestamp string) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", pipeline, driver, runid, timestamp, uuid.New().String())
}

// Insert log entry
func (m *LogModel) Insert(log types.Log) error {
	key := []byte(makeKey(log.Pipeline, log.Driver, log.RunID, log.Timestamp))
	value, err := json.Marshal(log)
	if err != nil {
		return err
	}
	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Query logs by filters (any combination)
func (m *LogModel) Query(pipeline, driver, runid string) ([]types.Log, error) {
	logs := []types.Log{}

	err := m.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := ""
		if pipeline != "" {
			prefix += pipeline
		}
		prefix += "|"
		if driver != "" {
			prefix += driver
		}
		prefix += "|"
		if runid != "" {
			prefix += runid
		}

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := string(item.Key())

			if pipeline != "" && !matchPart(k, 0, pipeline) {
				continue
			}
			if driver != "" && !matchPart(k, 1, driver) {
				continue
			}
			if runid != "" && !matchPart(k, 2, runid) {
				continue
			}

			err := item.Value(func(v []byte) error {
				var l types.Log
				if e := json.Unmarshal(v, &l); e == nil {
					logs = append(logs, l)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return logs, err
}

// Helper: check Nth part of key
func matchPart(key string, idx int, expected string) bool {
	parts := make([]string, 0, 4)
	curr := ""
	for i := 0; i < len(key); i++ {
		if key[i] == '|' {
			parts = append(parts, curr)
			curr = ""
		} else {
			curr += string(key[i])
		}
	}
	parts = append(parts, curr)

	if idx < len(parts) {
		return parts[idx] == expected
	}
	return false
}
