package deps

import (
	"fmt"
	"time"

	"go.etcd.io/etcd/server/v3/embed"
)

func StartEmbeddedEtcd() (*embed.Etcd, error) {
	cfg := embed.NewConfig()
	cfg.Dir = "default.etcd" // where etcd will store data
	cfg.Logger = "zap"
	// disable clustering if you only want a single node
	cfg.ClusterState = "new"

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		return nil, err
	}

	select {
	case <-e.Server.ReadyNotify():
		fmt.Println("Embedded etcd is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop()
		return nil, fmt.Errorf("etcd server took too long to start")
	}

	return e, nil
}
