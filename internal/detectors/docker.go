package detectors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/lazyservice/lazyservice/internal/app"
)

// DockerDetector detects Docker containers
type DockerDetector struct {
	client *client.Client
}

// NewDockerDetector creates a new Docker detector
func NewDockerDetector() (*DockerDetector, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerDetector{client: cli}, nil
}

// Detect detects Docker containers
func (d *DockerDetector) Detect(ctx context.Context) ([]*app.Service, error) {
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	services := make([]*app.Service, 0)
	for _, c := range containers {
		service := d.containerToService(c)
		services = append(services, service)
	}

	return services, nil
}

// GetType returns the service type
func (d *DockerDetector) GetType() app.ServiceType {
	return app.ServiceTypeDocker
}

// containerToService converts a Docker container to a service
func (d *DockerDetector) containerToService(c types.Container) *app.Service {
	status := app.ServiceStatusUnknown
	if c.State == "running" {
		status = app.ServiceStatusRunning
	} else if c.State == "exited" {
		status = app.ServiceStatusStopped
	} else if c.State == "dead" || c.State == "restarting" {
		status = app.ServiceStatusFailed
	}

	// Extract ports
	ports := make([]string, 0)
	for _, port := range c.Ports {
		if port.PublicPort > 0 {
			ports = append(ports, fmt.Sprintf("%d", port.PublicPort))
		}
	}

	// Extract image name without tag
	imageParts := strings.Split(c.Image, ":")
	image := c.Image
	version := "latest"
	if len(imageParts) > 1 {
		image = imageParts[0]
		version = imageParts[1]
	}

	// Get container name (remove leading slash)
	name := c.Names[0]
	if strings.HasPrefix(name, "/") {
		name = name[1:]
	}

	// Determine health check
	healthCheck := ""
	if c.State == "running" {
		if len(ports) > 0 {
			healthCheck = fmt.Sprintf("http://localhost:%s/health", ports[0])
		}
	}

	// Check for warnings
	warnings := make([]string, 0)
	if c.State == "restarting" {
		warnings = append(warnings, "Container is restarting")
	}

	service := &app.Service{
		ID:          c.ID[:12],
		Name:        name,
		Type:        app.ServiceTypeDocker,
		Status:      status,
		CreatedAt:   time.Unix(c.Created, 0),
		Ports:       ports,
		Image:       image,
		Version:     version,
		HealthCheck: healthCheck,
		Warnings:    warnings,
		Environment: make(map[string]string),
	}

	return service
}

// Close closes the Docker client
func (d *DockerDetector) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}
