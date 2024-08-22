package container

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"log"
)

type Docker struct {
	Client *client.Client
}

func NewDockerClient() *Docker {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &Docker{Client: cli}
}

func (d *Docker) ListContainer() ([]ContainerDTO, error) {
	var containerDTO []ContainerDTO
	list, err := d.Client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, v := range list {
		containerDTO = append(containerDTO, ContainerDTO{
			Id:     v.ID[0:10],
			Name:   v.Names[0],
			Status: v.Status,
			Describe: Describe{
				Labels: v.Labels,
				Image:  v.Image,
			},
		})
	}
	state, err := d.getContainerState(containerDTO[0].Name)
	if err != nil {
		log.Fatalf("state: %v", err)
	}
	log.Fatalf("containerJson: %s", state)
	return containerDTO, nil
}

func (d *Docker) getContainerState(containerName string) (string, error) {
	stats, err := d.Client.ContainerStats(context.Background(), containerName, false)
	if err != nil {
		return "", err
	}
	defer stats.Body.Close()
	body, err := io.ReadAll(stats.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (d *Docker) StartContainer(containerId string) error {
	if err := d.Client.ContainerStart(context.Background(), containerId, container.StartOptions{}); err != nil {
		return err
	}
	return nil
}

func (d *Docker) StopContainer(containerId string) error {
	if err := d.Client.ContainerStop(context.Background(), containerId, container.StopOptions{}); err != nil {
		return err
	}
	return nil
}

func (d *Docker) PublishContainer() error {
	//TODO implement me
	panic("implement me")
}
