package detectors

import (
	"context"
	"fmt"

	"github.com/lazyservice/lazyservice/internal/app"
)

// KubernetesDetector detects Kubernetes services (stub implementation)
type KubernetesDetector struct{}

// NewKubernetesDetector creates a new Kubernetes detector
// Note: This is a stub implementation. Full K8s support requires additional dependencies.
func NewKubernetesDetector() (*KubernetesDetector, error) {
	return nil, fmt.Errorf("kubernetes support not available in this build")
}

// Detect detects Kubernetes pods
func (k *KubernetesDetector) Detect(ctx context.Context) ([]*app.Service, error) {
	return []*app.Service{}, nil
}

// GetType returns the service type
func (k *KubernetesDetector) GetType() app.ServiceType {
	return app.ServiceTypeKubernetes
}
