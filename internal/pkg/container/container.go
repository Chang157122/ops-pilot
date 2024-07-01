package container

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

// GetContainerList 获取容器列表
func (d *Docker) GetContainerList() ContainerDTO {
	var containerDTO ContainerDTO
	list, err := d.Client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, v := range list {
		containerDTO.Id = v.ID[0:10]
		containerDTO.Name = v.Names[0]
		containerDTO.Image = v.Image
		containerDTO.Status = v.Status
	}
	return containerDTO
}
