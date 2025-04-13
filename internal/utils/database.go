/*
Copyright © 2024 Cranom Technologies Limited, Beingana Jim Junior and Contributors
*/
package utils

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdClient is a struct that holds the etcd client and the context
type EtcdClient struct {
	Client   *clientv3.Client
	Ctx      context.Context
	Cancel   context.CancelFunc
	Endpoint string
}

// NewEtcdClient creates a new etcd client
func NewEtcdClient(endpoint string) (*EtcdClient, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 5 * 1000,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EtcdClient{
		Client:   client,
		Ctx:      ctx,
		Cancel:   cancel,
		Endpoint: endpoint,
	}, nil
}
