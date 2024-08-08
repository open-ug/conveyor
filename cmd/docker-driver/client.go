package dockerdriver

import (
	"fmt"

	"github.com/docker/docker/client"
)

func GetDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client: %v", err)
	}
	return cli, nil
}
