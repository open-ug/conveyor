package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type PipelineModel struct {
	Client *clientv3.Client
	DB     *badger.DB
}

func NewPipelineModel(cli *clientv3.Client, db *badger.DB) *PipelineModel {
	return &PipelineModel{
		Client: cli,
		DB:     db,
	}
}

func (pm *PipelineModel) CreatePipeline(pipeline *types.Pipeline) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipelineKey := fmt.Sprintf("/pipelines/%s", pipeline.Name)
	pipelineValue, err := json.Marshal(pipeline)
	if err != nil {
		return err
	}

	_, err = pm.Client.Put(ctx, pipelineKey, string(pipelineValue))
	if err != nil {
		return err
	}

	return nil
}

func (pm *PipelineModel) BadgerCreatePipeline(pipeline *types.Pipeline) error {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipelineKey := fmt.Sprintf("/pipelines/%s", name)
	resp, err := pm.Client.Get(ctx, pipelineKey)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("pipeline not found")
	}

	var pipeline types.Pipeline
	err = json.Unmarshal(resp.Kvs[0].Value, &pipeline)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (pm *PipelineModel) BadgerGetPipeline(name string) (*types.Pipeline, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipelineKey := fmt.Sprintf("/pipelines/%s", pipeline.Name)
	pipelineValue, err := json.Marshal(pipeline)
	if err != nil {
		return err
	}

	_, err = pm.Client.Put(ctx, pipelineKey, string(pipelineValue))
	if err != nil {
		return err
	}

	return nil
}

func (pm *PipelineModel) BadgerUpdatePipeline(pipeline *types.Pipeline) error {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipelineKey := fmt.Sprintf("/pipelines/%s", name)
	_, err := pm.Client.Delete(ctx, pipelineKey)
	if err != nil {
		return err
	}

	return nil
}

func (pm *PipelineModel) BadgerDeletePipeline(name string) error {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipelinePrefix := "/pipelines/"
	resp, err := pm.Client.Get(ctx, pipelinePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var pipelines []*types.Pipeline
	for _, kv := range resp.Kvs {
		var pipeline types.Pipeline
		err := json.Unmarshal(kv.Value, &pipeline)
		if err != nil {
			return nil, err
		}
		pipelines = append(pipelines, &pipeline)
	}

	return pipelines, nil
}

func (pm *PipelineModel) BadgerListPipelines() ([]*types.Pipeline, error) {
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
