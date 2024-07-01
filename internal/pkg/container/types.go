package container

// ContainerDTO 容器数据结构
type ContainerDTO struct {
	Id     string `json:"container_id"`
	Name   string `json:"container_name"`
	Image  string `json:"container_image"`
	Status string `json:"container_status"`
	Ports  string `json:"ports"`
}
