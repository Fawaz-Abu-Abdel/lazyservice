package detectors

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/lazyservice/lazyservice/internal/app"
)

// SystemdDetector detects systemd services
type SystemdDetector struct{}

// NewSystemdDetector creates a new systemd detector
func NewSystemdDetector() *SystemdDetector {
	return &SystemdDetector{}
}

// Detect detects systemd services
func (s *SystemdDetector) Detect(ctx context.Context) ([]*app.Service, error) {
	// Check if systemctl is available
	_, err := exec.LookPath("systemctl")
	if err != nil {
		// systemctl not available (Windows, Mac without systemd, etc.)
		return []*app.Service{}, nil
	}

	cmd := exec.CommandContext(ctx, "systemctl", "list-units", "--type=service", "--all", "--no-pager", "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	services := make([]*app.Service, 0)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "UNIT") || strings.HasPrefix(line, "●") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		serviceName := fields[0]
		loadState := fields[1]
		activeState := fields[2]
		subState := fields[3]

		// Skip if not loaded
		if loadState != "loaded" {
			continue
		}

		status := s.convertStatus(activeState, subState)

		// Get service details
		service := &app.Service{
			ID:          serviceName,
			Name:        strings.TrimSuffix(serviceName, ".service"),
			Type:        app.ServiceTypeSystemd,
			Status:      status,
			CreatedAt:   time.Now(), // systemd doesn't easily expose creation time
			Environment: make(map[string]string),
		}

		// Try to get more details
		s.enrichServiceDetails(ctx, service)

		services = append(services, service)
	}

	return services, nil
}

// GetType returns the service type
func (s *SystemdDetector) GetType() app.ServiceType {
	return app.ServiceTypeSystemd
}

// convertStatus converts systemd status to service status
func (s *SystemdDetector) convertStatus(activeState, subState string) app.ServiceStatus {
	if activeState == "active" && subState == "running" {
		return app.ServiceStatusRunning
	}
	if activeState == "inactive" || activeState == "dead" {
		return app.ServiceStatusStopped
	}
	if activeState == "failed" {
		return app.ServiceStatusFailed
	}
	return app.ServiceStatusUnknown
}

// enrichServiceDetails gets additional service details
func (s *SystemdDetector) enrichServiceDetails(ctx context.Context, service *app.Service) {
	// Get service file path
	cmd := exec.CommandContext(ctx, "systemctl", "show", service.ID, "-p", "FragmentPath", "--value")
	if output, err := cmd.Output(); err == nil {
		path := strings.TrimSpace(string(output))
		if path != "" {
			service.Environment["ServiceFile"] = path
		}
	}

	// Get main PID
	cmd = exec.CommandContext(ctx, "systemctl", "show", service.ID, "-p", "MainPID", "--value")
	if output, err := cmd.Output(); err == nil {
		pid := strings.TrimSpace(string(output))
		if pid != "0" && pid != "" {
			service.Environment["PID"] = pid
		}
	}

	// Check for warnings (if service is in degraded state)
	if service.Status == app.ServiceStatusRunning {
		cmd = exec.CommandContext(ctx, "systemctl", "is-active", service.ID)
		if err := cmd.Run(); err != nil {
			service.Warnings = append(service.Warnings, "Service may be degraded")
		}
	}
}

// GetServiceLogs gets logs for a systemd service
func (s *SystemdDetector) GetServiceLogs(ctx context.Context, serviceName string, lines int) (string, error) {
	cmd := exec.CommandContext(ctx, "journalctl", "-u", serviceName, "-n", fmt.Sprintf("%d", lines), "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
