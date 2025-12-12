package app

import (
	"context"
	"testing"
	"time"
)

// MockDetector implements ServiceDetector for testing
type MockDetector struct {
	services []*Service
	err      error
}

func (m *MockDetector) Detect(ctx context.Context) ([]*Service, error) {
	return m.services, m.err
}

func (m *MockDetector) GetType() ServiceType {
	return ServiceTypeDocker
}

// MockMetricsCollector implements MetricsCollector for testing
type MockMetricsCollector struct {
	metrics *ServiceMetrics
	err     error
}

func (m *MockMetricsCollector) CollectMetrics(ctx context.Context, service *Service) (*ServiceMetrics, error) {
	return m.metrics, m.err
}

func TestNewApp(t *testing.T) {
	app := NewApp()
	
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
	
	if app.services == nil {
		t.Error("services map not initialized")
	}
	
	if app.detectors == nil {
		t.Error("detectors slice not initialized")
	}
	
	if app.metricsCollectors == nil {
		t.Error("metricsCollectors map not initialized")
	}
}

func TestRegisterDetector(t *testing.T) {
	app := NewApp()
	detector := &MockDetector{}
	
	app.RegisterDetector(detector)
	
	if len(app.detectors) != 1 {
		t.Errorf("Expected 1 detector, got %d", len(app.detectors))
	}
}

func TestRegisterMetricsCollector(t *testing.T) {
	app := NewApp()
	collector := &MockMetricsCollector{}
	
	app.RegisterMetricsCollector(ServiceTypeDocker, collector)
	
	if len(app.metricsCollectors) != 1 {
		t.Errorf("Expected 1 metrics collector, got %d", len(app.metricsCollectors))
	}
}

func TestRefreshServices(t *testing.T) {
	app := NewApp()
	
	// Create mock service
	mockService := &Service{
		ID:     "test-service",
		Name:   "Test Service",
		Type:   ServiceTypeDocker,
		Status: ServiceStatusRunning,
	}
	
	detector := &MockDetector{
		services: []*Service{mockService},
	}
	
	app.RegisterDetector(detector)
	app.RefreshServices()
	
	services := app.GetServices()
	if len(services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(services))
	}
	
	if services[0].ID != "test-service" {
		t.Errorf("Expected service ID 'test-service', got '%s'", services[0].ID)
	}
}

func TestGetService(t *testing.T) {
	app := NewApp()
	
	mockService := &Service{
		ID:     "test-service",
		Name:   "Test Service",
		Type:   ServiceTypeDocker,
		Status: ServiceStatusRunning,
	}
	
	detector := &MockDetector{
		services: []*Service{mockService},
	}
	
	app.RegisterDetector(detector)
	app.RefreshServices()
	
	service := app.GetService("test-service")
	if service == nil {
		t.Error("GetService returned nil for existing service")
	}
	
	if service.Name != "Test Service" {
		t.Errorf("Expected service name 'Test Service', got '%s'", service.Name)
	}
	
	// Test non-existent service
	nonExistent := app.GetService("non-existent")
	if nonExistent != nil {
		t.Error("GetService should return nil for non-existent service")
	}
}

func TestAppStop(t *testing.T) {
	app := NewApp()
	
	// Start the app
	app.Start()
	
	// Stop the app
	app.Stop()
	
	// Verify context is cancelled
	select {
	case <-app.ctx.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("App context was not cancelled after Stop()")
	}
}
