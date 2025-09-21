package utils

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/spf13/viper"
	"go.etcd.io/etcd/server/v3/embed"
	"go.etcd.io/etcd/server/v3/etcdserver/api/v3client"
)

// EtcdClient wraps the embedded etcd server and client
type EtcdClient struct {
	Client     *clientv3.Client
	Ctx        context.Context
	Cancel     context.CancelFunc
	Endpoint   string
	ServerStop func() // clean shutdown of embedded etcd
}

// NewEtcdClient starts an embedded etcd and returns a connected client
func NewEtcdClient() (*EtcdClient, error) {
	// --- configure embedded etcd ---
	cfg := embed.NewConfig()
	conveyorDataDir := viper.GetString("api.data")
	cfg.Dir = conveyorDataDir + "/etcd"
	cfg.Logger = "zap"
	cfg.LogOutputs = []string{conveyorDataDir + "/etcd.log"}
	cfg.ClusterState = "new"

	if IsTestMode() {
		cfg.Dir = filepath.Join(os.TempDir(), cfg.Name)
		cfg.LogOutputs = []string{"stderr"}
		cfg.ListenClientUrls = []url.URL{{Scheme: "http", Host: "localhost:0"}}
		cfg.ListenPeerUrls = []url.URL{{Scheme: "http", Host: "localhost:0"}}
	}

	// Start etcd
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to start embedded etcd: %w", err)
	}

	// Wait until ready
	select {
	case <-e.Server.ReadyNotify():
		log.Println("Embedded etcd is ready on")
	case <-time.After(60 * time.Second):
		e.Server.Stop()
		return nil, fmt.Errorf("etcd server took too long to start")
	}

	// --- connect etcd client ---
	client := v3client.New(e.Server)

	ctx, cancel := context.WithCancel(context.Background())

	return &EtcdClient{
		Client:     client,
		Ctx:        ctx,
		Cancel:     cancel,
		ServerStop: e.Close,
	}, nil
}
