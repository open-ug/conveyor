package dockerdriver

import (
	"context"
	"strconv"

	craneTypes "crane.cloud.cranom.tech/cmd/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gofiber/fiber/v2/log"
)

func CreateContainer(
	dockerClient *client.Client,
	app *craneTypes.Application,
) error {

	ctx := context.Background()

	// Pull the image (if not already pulled)
	reader, err := dockerClient.ImagePull(ctx, app.Spec.Image, image.PullOptions{})
	if err != nil {
		log.Fatalf("Error pulling image: %v", err)
		return err
	}
	defer reader.Close()

	containerCfg, hostCfg, networkCfg, err := GenerateContainerConfig(app)

	if err != nil {
		log.Fatalf("Error pulling image: %v", err)
		return err
	}

	// Create the container
	resp, err := dockerClient.ContainerCreate(ctx, containerCfg, hostCfg, networkCfg, nil, app.Name)

	if err != nil {
		log.Fatalf("Error creating container: %v", err)
		return err
	}

	log.Infof("Container %s created", resp.ID)

	return nil
}

func GenerateContainerConfig(
	app *craneTypes.Application,
) (*container.Config, *container.HostConfig, *network.NetworkingConfig, error) {
	envVars := []string{}

	for i := 0; i < len(app.Spec.Env); i++ {
		env := app.Spec.Env[i]
		envVars = append(envVars, env.Name+"="+env.Value)
	}
	containerCfg := container.Config{
		Image: app.Spec.Image,
		Env:   envVars,
	}

	networkName := app.Spec.Network

	networkCfg := network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {
				NetworkID: networkName,
			},
		},
	}
	containerPort := strconv.Itoa(app.Spec.Ports[0].Internal)
	hostPort := strconv.Itoa(app.Spec.Ports[0].External)

	hostCfg := container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(containerPort + "/tcp"): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: hostPort,
				},
			},
		},
	}

	return &containerCfg, &hostCfg, &networkCfg, nil

}

func StartContainer(
	dockerClient *client.Client,
	app *craneTypes.Application,
) error {
	ctx := context.Background()

	err := dockerClient.ContainerStart(ctx, app.Name, container.StartOptions{})
	if err != nil {
		log.Fatalf("Error starting container: %v", err)
		return err
	}

	log.Infof("Container %s started", app.Name)

	return nil
}

func StopContainer(
	dockerClient *client.Client,
	app *craneTypes.Application,
) error {
	ctx := context.Background()

	err := dockerClient.ContainerStop(ctx, app.Name, container.StopOptions{})
	if err != nil {
		log.Fatalf("Error stopping container: %v", err)
		return err
	}

	log.Infof("Container %s stopped", app.Name)

	return nil
}

func DeleteContainer(
	dockerClient *client.Client,
	app *craneTypes.Application,
) error {
	ctx := context.Background()

	err := dockerClient.ContainerRemove(ctx, app.Name, container.RemoveOptions{})
	if err != nil {
		log.Fatalf("Error removing container: %v", err)
		return err
	}

	log.Infof("Container %s removed", app.Name)

	return nil
}
