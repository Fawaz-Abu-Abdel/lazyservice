package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/lazyservice/lazyservice/internal/app"
)

// ProcessMetricsCollector collects metrics for processes
type ProcessMetricsCollector struct{}

// NewProcessMetricsCollector creates a new process metrics collector
func NewProcessMetricsCollector() *ProcessMetricsCollector {
	return &ProcessMetricsCollector{}
}

// CollectMetrics collects metrics for a process
func (p *ProcessMetricsCollector) CollectMetrics(ctx context.Context, service *app.Service) (*app.ServiceMetrics, error) {
	// Get PID from service
	pidStr, ok := service.Environment["PID"]
	if !ok {
		return nil, nil
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return nil, err
	}

	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}

	// Get CPU percent
	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		cpuPercent = 0
	}

	// Get memory info
	memInfo, err := proc.MemoryInfo()
	memoryUsage := uint64(0)
	if err == nil {
		memoryUsage = memInfo.RSS
	}

	// Get memory percent
	memPercent, err := proc.MemoryPercent()
	if err != nil {
		memPercent = 0
	}

	// Get IO counters
	ioCounters, err := proc.IOCounters()
	diskRead := uint64(0)
	diskWrite := uint64(0)
	if err == nil {
		diskRead = ioCounters.ReadBytes
		diskWrite = ioCounters.WriteBytes
	}

	// Get create time for uptime
	createTime, err := proc.CreateTime()
	uptime := time.Duration(0)
	if err == nil {
		uptime = time.Since(time.Unix(createTime/1000, 0))
	}

	// Get thread count
	threadCount, err := proc.NumThreads()
	if err != nil {
		threadCount = 0
	}

	// Get connection count
	connections, err := proc.Connections()
	connCount := 0
	if err == nil {
		connCount = len(connections)
	}

	// Network stats - not easily available per process
	// Would need additional tools like nethogs or tcpdump integration

	metrics := &app.ServiceMetrics{
		CPUPercent:     cpuPercent,
		CPUCores:       cpuPercent / 100.0,
		MemoryUsage:    memoryUsage,
		MemoryLimit:    0, // Not applicable for processes
		MemoryPercent:  float64(memPercent),
		NetworkIn:      0,
		NetworkOut:     0,
		DiskIORead:     diskRead,
		DiskIOWrite:    diskWrite,
		RequestsPerSec: 0,
		ErrorRate:      0,
		ResponseTime:   0,
		Uptime:         uptime,
		ThreadCount:    threadCount,
		Connections:    connCount,
	}

	return metrics, nil
}
