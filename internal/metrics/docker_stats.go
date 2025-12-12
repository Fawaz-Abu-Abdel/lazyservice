package metrics

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/lazyservice/lazyservice/internal/app"
)

// DockerMetricsCollector collects metrics for Docker containers
type DockerMetricsCollector struct {
	client *client.Client
}

// NewDockerMetricsCollector creates a new Docker metrics collector
func NewDockerMetricsCollector() (*DockerMetricsCollector, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerMetricsCollector{client: cli}, nil
}

// CollectMetrics collects metrics for a Docker container
func (d *DockerMetricsCollector) CollectMetrics(ctx context.Context, service *app.Service) (*app.ServiceMetrics, error) {
	// Get container stats
	stats, err := d.client.ContainerStats(ctx, service.ID, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	// Parse stats
	var containerStats types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&containerStats); err != nil {
		if err != io.EOF {
			return nil, err
		}
	}

	// Calculate CPU percentage
	cpuPercent := calculateCPUPercent(&containerStats)

	// Calculate memory usage
	memoryUsage := containerStats.MemoryStats.Usage
	memoryLimit := containerStats.MemoryStats.Limit
	memoryPercent := 0.0
	if memoryLimit > 0 {
		memoryPercent = float64(memoryUsage) / float64(memoryLimit) * 100
	}

	// Network stats
	var networkIn, networkOut uint64
	for _, network := range containerStats.Networks {
		networkIn += network.RxBytes
		networkOut += network.TxBytes
	}

	// Disk I/O stats
	var diskRead, diskWrite uint64
	for _, bioEntry := range containerStats.BlkioStats.IoServiceBytesRecursive {
		if bioEntry.Op == "Read" {
			diskRead += bioEntry.Value
		} else if bioEntry.Op == "Write" {
			diskWrite += bioEntry.Value
		}
	}

	// Get container inspect info for uptime
	inspect, err := d.client.ContainerInspect(ctx, service.ID)
	uptime := time.Duration(0)
	if err == nil {
		startTime, err := time.Parse(time.RFC3339Nano, inspect.State.StartedAt)
		if err == nil {
			uptime = time.Since(startTime)
		}
	}

	metrics := &app.ServiceMetrics{
		CPUPercent:     cpuPercent,
		CPUCores:       cpuPercent / 100.0,
		MemoryUsage:    memoryUsage,
		MemoryLimit:    memoryLimit,
		MemoryPercent:  memoryPercent,
		NetworkIn:      networkIn,
		NetworkOut:     networkOut,
		DiskIORead:     diskRead,
		DiskIOWrite:    diskWrite,
		RequestsPerSec: 0, // Would need additional monitoring
		ErrorRate:      0, // Would need additional monitoring
		ResponseTime:   0, // Would need additional monitoring
		Uptime:         uptime,
	}

	return metrics, nil
}

// calculateCPUPercent calculates CPU percentage from stats
func calculateCPUPercent(stats *types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)

	if systemDelta > 0 && cpuDelta > 0 {
		cpuPercent := (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
		return cpuPercent
	}
	return 0.0
}

// Close closes the Docker client
func (d *DockerMetricsCollector) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}
