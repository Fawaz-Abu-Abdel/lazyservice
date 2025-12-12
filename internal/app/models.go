package app

import "time"

// ServiceType represents the type of service
type ServiceType string

const (
	ServiceTypeDocker     ServiceType = "docker"
	ServiceTypeKubernetes ServiceType = "kubernetes"
	ServiceTypeSystemd    ServiceType = "systemd"
	ServiceTypeProcess    ServiceType = "process"
)

// ServiceStatus represents the status of a service
type ServiceStatus string

const (
	ServiceStatusRunning ServiceStatus = "running"
	ServiceStatusStopped ServiceStatus = "stopped"
	ServiceStatusFailed  ServiceStatus = "failed"
	ServiceStatusUnknown ServiceStatus = "unknown"
)

// Service represents a unified service across different platforms
type Service struct {
	ID          string
	Name        string
	Type        ServiceType
	Status      ServiceStatus
	CreatedAt   time.Time
	Ports       []string
	Image       string
	Namespace   string
	Version     string
	HealthCheck string
	Warnings    []string
	Environment map[string]string
	Metrics     *ServiceMetrics
	Statistics  *ServiceStatistics
}

// ServiceMetrics holds real-time metrics for a service
type ServiceMetrics struct {
	CPUPercent     float64
	CPUCores       float64
	MemoryUsage    uint64
	MemoryLimit    uint64
	MemoryPercent  float64
	NetworkIn      uint64
	NetworkOut     uint64
	DiskIORead     uint64
	DiskIOWrite    uint64
	RequestsPerSec float64
	ErrorRate      float64
	ResponseTime   float64
	Uptime         time.Duration
	ThreadCount    int32
	Connections    int
	History        *MetricsHistory
}

// MetricsHistory stores historical metrics data
type MetricsHistory struct {
	CPU      []float64
	Memory   []float64
	Network  []float64
	Requests []float64
	MaxSize  int
}

// NewMetricsHistory creates a new metrics history with max size
func NewMetricsHistory(maxSize int) *MetricsHistory {
	return &MetricsHistory{
		CPU:      make([]float64, 0, maxSize),
		Memory:   make([]float64, 0, maxSize),
		Network:  make([]float64, 0, maxSize),
		Requests: make([]float64, 0, maxSize),
		MaxSize:  maxSize,
	}
}

// AddMetrics adds a new metrics entry to history
func (h *MetricsHistory) AddMetrics(cpu, memory, network, requests float64) {
	h.CPU = append(h.CPU, cpu)
	h.Memory = append(h.Memory, memory)
	h.Network = append(h.Network, network)
	h.Requests = append(h.Requests, requests)

	// Keep only the last maxSize entries
	if len(h.CPU) > h.MaxSize {
		h.CPU = h.CPU[len(h.CPU)-h.MaxSize:]
	}
	if len(h.Memory) > h.MaxSize {
		h.Memory = h.Memory[len(h.Memory)-h.MaxSize:]
	}
	if len(h.Network) > h.MaxSize {
		h.Network = h.Network[len(h.Network)-h.MaxSize:]
	}
	if len(h.Requests) > h.MaxSize {
		h.Requests = h.Requests[len(h.Requests)-h.MaxSize:]
	}
}

// GetStatusEmoji returns an emoji for the service status
func (s *Service) GetStatusEmoji() string {
	switch s.Status {
	case ServiceStatusRunning:
		if len(s.Warnings) > 0 {
			return "🟡"
		}
		if s.Type == ServiceTypeKubernetes {
			return "🟢"
		}
		if s.Type == ServiceTypeDocker {
			return "🔵"
		}
		return "🟢"
	case ServiceStatusStopped:
		return "⚫"
	case ServiceStatusFailed:
		return "🔴"
	default:
		return "⚪"
	}
}
