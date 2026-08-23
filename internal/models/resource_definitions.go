package models

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/open-ug/conveyor/pkg/types"
)

type ResourceDefinitionModel struct {
	DB *badger.DB
}

func NewResourceDefinitionModel(db *badger.DB) *ResourceDefinitionModel {
	return &ResourceDefinitionModel{
		DB: db,
	}
}

func (m *ResourceDefinitionModel) key(name string) string {
	return fmt.Sprintf("/resource_definitions/%s", name)
}

// Insert adds a new resource definition
// It returns an error if a resource definition with the same name already exists.
func (m *ResourceDefinitionModel) Insert(name string, definition *types.ResourceDefinition) error {
	err := m.DB.Update(func(txn *badger.Txn) error {
		key := m.key(name)

		// Check existence
		_, err := txn.Get([]byte(key))
		if err == nil {
			return fmt.Errorf("resource definition with name %s already exists", name)
		} else if err != badger.ErrKeyNotFound {
			return err
		}

		definitionValue, err := json.Marshal(definition)
		if err != nil {
			return err
		}

		err = txn.Set([]byte(key), definitionValue)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// FindOne retrieves a single resource definition by its name.
// It returns the resource definition data as a byte slice or an error if not found.
func (m *ResourceDefinitionModel) FindOne(name string) (*types.ResourceDefinition, error) {
	var definition types.ResourceDefinition

	err := m.DB.View(func(txn *badger.Txn) error {
		key := m.key(name)

		item, err := txn.Get([]byte(key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return fmt.Errorf("resource definition with name %s not found", name)
			}
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &definition)
		})
	})

	if err != nil {
		return nil, err
	}

	return &definition, nil
}

// Delete removes a resource definition by its name.
// It returns an error if the resource definition does not exist.
func (m *ResourceDefinitionModel) Delete(name string) error {
	err := m.DB.Update(func(txn *badger.Txn) error {
		key := m.key(name)

		// Check existence
		_, err := txn.Get([]byte(key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return fmt.Errorf("resource definition with name %s not found", name)
			}
			return err
		}

		err = txn.Delete([]byte(key))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// FindAll retrieves all resource definitions.
// It returns a slice of resource definition names or an error if no definitions are found.
func (m *ResourceDefinitionModel) FindAll() ([]*types.ResourceDefinition, error) {
	var resourceDefinitions []*types.ResourceDefinition

	err := m.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("/resource_definitions/")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var definition types.ResourceDefinition
				if err := json.Unmarshal(val, &definition); err != nil {
					return err
				}
				resourceDefinitions = append(resourceDefinitions, &definition)
				return nil
			})
			if err != nil {
				return err
			}
		}

		if len(resourceDefinitions) == 0 {
			return fmt.Errorf("no resource definitions found")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resourceDefinitions, nil
}

// Update modifies an existing resource definition's data.
// It returns an error if the resource definition does not exist.
func (m *ResourceDefinitionModel) Update(name string, definition *types.ResourceDefinition) error {
	err := m.DB.Update(func(txn *badger.Txn) error {
		key := m.key(name)

		// Check existence
		_, err := txn.Get([]byte(key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return fmt.Errorf("resource definition with name %s not found", name)
			}
			return err
		}

		definitionValue, err := json.Marshal(definition)
		if err != nil {
			return err
		}

		err = txn.Set([]byte(key), definitionValue)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
