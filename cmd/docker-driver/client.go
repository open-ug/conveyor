/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package dockerdriver

import (
	"fmt"

	"github.com/docker/docker/client"
)

func GetDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error creating docker client: %v", err)
	}
	return cli, nil
}
