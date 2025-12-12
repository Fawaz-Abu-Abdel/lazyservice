package detectors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/lazyservice/lazyservice/internal/app"

)

// ProcessDetector detects running processes
type ProcessDetector struct {
	// Filter processes by name patterns
	namePatterns []string
}

// NewProcessDetector creates a new process detector
func NewProcessDetector(namePatterns []string) *ProcessDetector {
	if len(namePatterns) == 0 {
		// Default patterns for common service processes
		namePatterns = []string{
			"node", "python", "java", "ruby", "php", "nginx", "apache",
			"redis", "postgres", "mysql", "mongodb", "pm2",
		}
	}
	return &ProcessDetector{namePatterns: namePatterns}
}

// Detect detects running processes
func (p *ProcessDetector) Detect(ctx context.Context) ([]*app.Service, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	services := make([]*app.Service, 0)
	for _, proc := range processes {
		service, err := p.processToService(proc)
		if err != nil {
			continue
		}
		if service != nil {
			services = append(services, service)
		}
	}

	return services, nil
}

// GetType returns the service type
func (p *ProcessDetector) GetType() app.ServiceType {
	return app.ServiceTypeProcess
}

// processToService converts a process to a service
func (p *ProcessDetector) processToService(proc *process.Process) (*app.Service, error) {
	name, err := proc.Name()
	if err != nil {
		return nil, err
	}

	// Filter by name patterns
	if !p.matchesPattern(name) {
		return nil, nil
	}

	// Get command line
	cmdline, err := proc.Cmdline()
	if err != nil {
		cmdline = name
	}

	// Get create time
	createTime, err := proc.CreateTime()
	if err != nil {
		createTime = 0
	}

	// Get status
	status, err := proc.Status()
	if err != nil {
		status = []string{"unknown"}
	}

	// Convert status
	serviceStatus := p.convertStatus(status)

	// Get working directory
	cwd, _ := proc.Cwd()

	service := &app.Service{
		ID:        fmt.Sprintf("proc-%d", proc.Pid),
		Name:      name,
		Type:      app.ServiceTypeProcess,
		Status:    serviceStatus,
		CreatedAt: time.Unix(createTime/1000, 0),
		Environment: map[string]string{
			"PID":     fmt.Sprintf("%d", proc.Pid),
			"Command": cmdline,
			"CWD":     cwd,
		},
	}

	return service, nil
}

// matchesPattern checks if process name matches any pattern
func (p *ProcessDetector) matchesPattern(name string) bool {
	nameLower := strings.ToLower(name)
	for _, pattern := range p.namePatterns {
		if strings.Contains(nameLower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// convertStatus converts process status to service status
func (p *ProcessDetector) convertStatus(status []string) app.ServiceStatus {
	if len(status) == 0 {
		return app.ServiceStatusUnknown
	}

	statusStr := strings.ToLower(status[0])
	switch statusStr {
	case "running", "r":
		return app.ServiceStatusRunning
	case "sleeping", "s", "idle", "i":
		return app.ServiceStatusRunning
	case "stopped", "t":
		return app.ServiceStatusStopped
	case "zombie", "z", "dead", "x":
		return app.ServiceStatusFailed
	default:
		return app.ServiceStatusUnknown
	}
}
