package watchdog

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"net"
	"opsPilot/internal/common/dto"
	"opsPilot/internal/pkg/log"
	"os"
)

// GetMetrics 获取当前主机监控指标
func GetMetrics() (metricsDTO dto.MetricsDTO) {

	hostname, _ := os.Hostname()
	ips, _ := net.LookupIP(hostname)
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		log.Logger.Errorf("get cpu usage failed err: %v", err)
		panic(err)
	}
	cpuTotal, _ := cpu.Counts(false)

	memoryStat, err := mem.VirtualMemory()
	if err != nil {
		log.Logger.Errorf("get memory usage failed err: %v", err)
		panic(err)
	}

	processes, _ := process.Processes()
	var processList []dto.Process
	for _, proc := range processes {
		procName, _ := proc.Name()
		command, _ := proc.Cmdline()
		if command == "" || len(command) == 0 {
			continue
		}
		procMemInfo, _ := proc.MemoryInfo()
		procCpuPercent, _ := proc.CPUPercent()
		processList = append(processList, dto.Process{
			Pid:     proc.Pid,
			Name:    procName,
			Command: command,
			Memory:  fmt.Sprintf("%vMB", procMemInfo.RSS/1024/1024),
			Cpu:     fmt.Sprintf("%2.f%%", procCpuPercent),
		})
	}

	metricsDTO = dto.MetricsDTO{
		HostName:    hostname,
		IP:          ips[0].String(),
		MemoryUsage: fmt.Sprintf("%vMB", memoryStat.Used/1024/1024),
		MemoryTotal: fmt.Sprintf("%vMB", memoryStat.Total/1024/1024),
		CpuUsage:    fmt.Sprintf("%.2f%%", cpuUsage[0]),
		CpuTotal:    cpuTotal,
		ProcessList: processList,
	}

	return metricsDTO
}
