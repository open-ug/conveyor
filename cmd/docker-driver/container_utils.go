/*
Copyright © 2024 Cranom Technologies Limited, Beingana Jim Junior and Contributors
*/
package dockerdriver

import (
	"context"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gofiber/fiber/v2/log"
	craneTypes "github.com/open-ug/conveyor/pkg/types"
)

func CreateAppVolumes(
	app *craneTypes.Application,
	dockerClient *client.Client,
) error {
	ctx := context.Background()

	for i := 0; i < len(app.Spec.Volumes); i++ {
		vol := app.Spec.Volumes[i]
		// check if volume exists
		_, verr := dockerClient.VolumeInspect(ctx, vol.VolumeName)
		if verr != nil {
			log.Infof("Volume %s already exists", vol.VolumeName)
			continue
		} else {
			_, err := dockerClient.VolumeCreate(ctx, volume.CreateOptions{
				Name: vol.VolumeName,
			})

			if err != nil {
				log.Fatalf("Error creating volume: %v", err)
				return err
			}
			log.Infof("Volume %s created", vol.VolumeName)
		}
	}

	return nil

}

func CreateNetwork(
	app *craneTypes.Application,
	dockerClient *client.Client,
) error {
	ctx := context.Background()

	// check if network exists
	_, eerr := dockerClient.NetworkInspect(ctx, app.Spec.Network, network.InspectOptions{})
	if eerr == nil {
		log.Infof("Network %s already exists", app.Spec.Network)
		return nil
	}

	_, err := dockerClient.NetworkCreate(ctx, app.Spec.Network, network.CreateOptions{})

	if err != nil {
		log.Fatalf("Error creating network: %v", err)
		return err
	}

	log.Infof("Network %s created", app.Spec.Network)

	return nil
}

func CreateContainer(
	dockerClient *client.Client,
	app *craneTypes.Application,
) error {

	ctx := context.Background()

	if app.Spec.Source.Type == "docker" {
		// Pull the image (if not already pulled)
		reader, err := dockerClient.ImagePull(ctx, app.Spec.Source.Image.ImageURI, image.PullOptions{})
		if err != nil {
			log.Fatalf("Error pulling image: %v", err)
			return err
		}
		defer reader.Close()
	}

	containerCfg, hostCfg, networkCfg, err := GenerateContainerConfig(app)
	if err != nil {
		log.Fatalf("Error pulling image: %v", err)
		return err
	}
	// create network

	if app.Spec.Network != "" {

		nerr := CreateNetwork(app, dockerClient)
		if nerr != nil {
			log.Fatalf("Error creating network: %v", nerr)
			return nerr
		}
	}

	verr := CreateAppVolumes(app, dockerClient)

	if verr != nil {
		log.Fatalf("Error creating volumes: %v", verr)
		return verr
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
	var imageURI string
	if app.Spec.Source.Type == "git" {
		// if the source is git, the imageURI is the buildpacks image
		imageURI = app.Name + "-bpimage"
	} else {
		imageURI = app.Spec.Source.Image.ImageURI
	}

	containerCfg := container.Config{
		Image: imageURI,
		Env:   envVars,
	}

	var networkCfg network.NetworkingConfig

	if app.Spec.Network != "" {
		networkName := app.Spec.Network

		networkCfg = network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				networkName: {
					NetworkID: networkName,
				},
			},
		}

	}

	// Create volume bindings
	volumeBindings := []string{}

	for i := 0; i < len(app.Spec.Volumes); i++ {
		vol := app.Spec.Volumes[i]
		volumeBindings = append(volumeBindings, vol.VolumeName+":"+vol.Path)
	}

	hostCfg := container.HostConfig{
		Binds: volumeBindings,
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
	}

	var containerPort string
	var hostPort string

	// If there are ports specified then create port bindings
	if len(app.Spec.Ports) != 0 {

		containerPort = strconv.Itoa(app.Spec.Ports[0].Internal)
		hostPort = strconv.Itoa(app.Spec.Ports[0].External)

		hostCfg.PortBindings = nat.PortMap{
			nat.Port(containerPort + "/tcp"): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: hostPort,
				},
			},
		}

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

func UpdateContainer(
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

	err = dockerClient.ContainerRemove(ctx, app.Name, container.RemoveOptions{})
	if err != nil {
		log.Fatalf("Error removing container: %v", err)
		return err
	}

	log.Infof("Container %s removed", app.Name)

	err = CreateContainer(dockerClient, app)
	if err != nil {
		log.Fatalf("Error creating container: %v", err)
		return err
	}

	err = StartContainer(dockerClient, app)
	if err != nil {
		log.Fatalf("Error starting container: %v", err)
		return err
	}

	return nil
}
