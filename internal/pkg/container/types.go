package container

// ContainerDTO 容器数据结构
type ContainerDTO struct {
	Id       string   `json:"container_id"`
	Name     string   `json:"container_name"`
	Status   string   `json:"container_status"`
	MemUsage string   `json:"mem_usage"`
	CpuUsage string   `json:"cpu_usage"`
	Describe Describe `json:"describe"`
}

type Describe struct {
	Image  string            `json:"container_image"`
	Labels map[string]string `json:"labels"`
	Ports  []string          `json:"ports"`
}
