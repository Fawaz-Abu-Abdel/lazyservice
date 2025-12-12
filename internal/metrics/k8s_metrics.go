package metrics

import (
	"context"

	"github.com/lazyservice/lazyservice/internal/app"
)

// K8sMetricsCollector collects metrics for Kubernetes pods (stub implementation)
type K8sMetricsCollector struct{}

// NewK8sMetricsCollector creates a new Kubernetes metrics collector
// Note: This is a stub implementation. Full K8s support requires additional dependencies.
func NewK8sMetricsCollector() (*K8sMetricsCollector, error) {
	return &K8sMetricsCollector{}, nil
}

// CollectMetrics collects metrics for a Kubernetes pod
func (k *K8sMetricsCollector) CollectMetrics(ctx context.Context, service *app.Service) (*app.ServiceMetrics, error) {
	// Stub implementation - returns empty metrics
	return &app.ServiceMetrics{}, nil
}
