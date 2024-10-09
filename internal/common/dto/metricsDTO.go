package dto

type MetricsDTO struct {
	IP          string    `json:"ip"`
	HostName    string    `json:"hostname"`
	MemoryUsage string    `json:"memory_usage"`
	CpuUsage    string    `json:"cpu_usage"`
	MemoryTotal string    `json:"memory_total"`
	CpuTotal    int       `json:"cpu_total"`
	ProcessList []Process `json:"Process_list"`
}

type Process struct {
	Pid     int32  `json:"pid"`
	Name    string `json:"name"`
	Command string `json:"command"`
	Memory  string `json:"memory"`
	Cpu     string `json:"cpu"`
}
