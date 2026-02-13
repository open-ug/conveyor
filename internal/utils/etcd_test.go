package utils_test

import (
	"context"
	"testing"
	"time"

	"github.com/open-ug/conveyor/internal/config"
	"github.com/open-ug/conveyor/internal/config/initialize"
	"github.com/open-ug/conveyor/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewEtcdClient_Integration(t *testing.T) {

	configFile, err := initialize.Run(&initialize.Options{
		Force:   true,
		TempDir: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}
	config.LoadTestEnvConfig(configFile)

	cfg, err := config.GetTestConfig()
	if err != nil {
		t.Fatalf("failed to get test config: %v", err)
	}

	client, err := utils.NewEtcdClient(&cfg)
	assert.NoError(t, err, "Expected to connect to etcd without error")
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)

	// Ensure context and cancel function are set
	assert.NotNil(t, client.Ctx)
	assert.NotNil(t, client.Cancel)

	// Perform a simple Put/Get to verify connectivity
	key := "test-key"
	value := "hello-etcd"

	// Put key
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = client.Client.Put(ctx, key, value)
	assert.NoError(t, err, "Expected to put key in etcd without error")

	// Get key
	resp, err := client.Client.Get(ctx, key)
	assert.NoError(t, err, "Expected to get key from etcd without error")
	assert.Equal(t, int64(1), resp.Count, "Expected 1 key to be returned")
	assert.Equal(t, value, string(resp.Kvs[0].Value), "Expected value to match")

	// Cleanup: delete key
	_, err = client.Client.Delete(ctx, key)
	assert.NoError(t, err, "Expected to delete key from etcd without error")

	// Close client
	client.Cancel()
	client.Client.Close()
	client.ServerStop()
}
