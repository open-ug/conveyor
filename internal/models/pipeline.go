package models

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/open-ug/conveyor/pkg/types"
)

type PipelineModel struct {
	DB *badger.DB
}

func NewPipelineModel(db *badger.DB) *PipelineModel {
	return &PipelineModel{
		DB: db,
	}
}

func (pm *PipelineModel) CreatePipeline(pipeline *types.Pipeline) error {
	err := pm.DB.Update(func(txn *badger.Txn) error {
		pipelineKey := fmt.Sprintf("/pipelines/%s", pipeline.Name)
		pipelineValue, err := json.Marshal(pipeline)
		if err != nil {
			return err
		}

		err = txn.Set([]byte(pipelineKey), pipelineValue)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (pm *PipelineModel) GetPipeline(name string) (*types.Pipeline, error) {
	var pipeline *types.Pipeline

	err := pm.DB.View(func(txn *badger.Txn) error {
		pipelineKey := fmt.Sprintf("/pipelines/%s", name)
		item, err := txn.Get([]byte(pipelineKey))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &pipeline)
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func (pm *PipelineModel) UpdatePipeline(pipeline *types.Pipeline) error {
	err := pm.DB.Update(func(txn *badger.Txn) error {
		pipelineKey := fmt.Sprintf("/pipelines/%s", pipeline.Name)
		pipelineValue, err := json.Marshal(pipeline)
		if err != nil {
			return err
		}

		err = txn.Set([]byte(pipelineKey), pipelineValue)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (pm *PipelineModel) DeletePipeline(name string) error {
	err := pm.DB.Update(func(txn *badger.Txn) error {
		pipelineKey := fmt.Sprintf("/pipelines/%s", name)
		err := txn.Delete([]byte(pipelineKey))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (pm *PipelineModel) ListPipelines() ([]*types.Pipeline, error) {
	var pipelines []*types.Pipeline

	err := pm.DB.View(func(txn *badger.Txn) error {
		pipelinePrefix := "/pipelines/"
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte(pipelinePrefix)); it.ValidForPrefix([]byte(pipelinePrefix)); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var pipeline types.Pipeline
				err := json.Unmarshal(val, &pipeline)
				if err != nil {
					return err
				}
				pipelines = append(pipelines, &pipeline)
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

	return pipelines, nil
}
