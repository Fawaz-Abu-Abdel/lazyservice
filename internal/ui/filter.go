package ui

import (
	"strings"

	"github.com/lazyservice/lazyservice/internal/app"
)

// FilterOptions defines filtering options for services
type FilterOptions struct {
	SearchTerm   string
	ServiceTypes []app.ServiceType
	StatusFilter []app.ServiceStatus
	ShowStopped  bool
}

// ServiceFilter handles service filtering logic
type ServiceFilter struct {
	options FilterOptions
}

// NewServiceFilter creates a new service filter
func NewServiceFilter() *ServiceFilter {
	return &ServiceFilter{
		options: FilterOptions{
			ServiceTypes: []app.ServiceType{
				app.ServiceTypeDocker,
				app.ServiceTypeKubernetes,
				app.ServiceTypeSystemd,
				app.ServiceTypeProcess,
			},
			StatusFilter: []app.ServiceStatus{
				app.ServiceStatusRunning,
				app.ServiceStatusStopped,
				app.ServiceStatusFailed,
				app.ServiceStatusUnknown,
			},
			ShowStopped: true,
		},
	}
}

// SetSearchTerm sets the search term for filtering
func (sf *ServiceFilter) SetSearchTerm(term string) {
	sf.options.SearchTerm = strings.ToLower(term)
}

// ToggleServiceType toggles a service type filter
func (sf *ServiceFilter) ToggleServiceType(serviceType app.ServiceType) {
	for i, t := range sf.options.ServiceTypes {
		if t == serviceType {
			// Remove if exists
			sf.options.ServiceTypes = append(sf.options.ServiceTypes[:i], sf.options.ServiceTypes[i+1:]...)
			return
		}
	}
	// Add if doesn't exist
	sf.options.ServiceTypes = append(sf.options.ServiceTypes, serviceType)
}

// ToggleStatus toggles a status filter
func (sf *ServiceFilter) ToggleStatus(status app.ServiceStatus) {
	for i, s := range sf.options.StatusFilter {
		if s == status {
			// Remove if exists
			sf.options.StatusFilter = append(sf.options.StatusFilter[:i], sf.options.StatusFilter[i+1:]...)
			return
		}
	}
	// Add if doesn't exist
	sf.options.StatusFilter = append(sf.options.StatusFilter, status)
}

// Filter filters services based on current options
func (sf *ServiceFilter) Filter(services []*app.Service) []*app.Service {
	filtered := make([]*app.Service, 0)

	for _, service := range services {
		if sf.matchesFilter(service) {
			filtered = append(filtered, service)
		}
	}

	return filtered
}

// matchesFilter checks if a service matches the current filter
func (sf *ServiceFilter) matchesFilter(service *app.Service) bool {
	// Check search term
	if sf.options.SearchTerm != "" {
		searchFields := []string{
			strings.ToLower(service.Name),
			strings.ToLower(service.Image),
			strings.ToLower(service.Namespace),
		}
		
		found := false
		for _, field := range searchFields {
			if strings.Contains(field, sf.options.SearchTerm) {
				found = true
				break
			}
		}
		
		// Also search in ports
		for _, port := range service.Ports {
			if strings.Contains(strings.ToLower(port), sf.options.SearchTerm) {
				found = true
				break
			}
		}
		
		if !found {
			return false
		}
	}

	// Check service type
	typeMatches := false
	for _, t := range sf.options.ServiceTypes {
		if service.Type == t {
			typeMatches = true
			break
		}
	}
	if !typeMatches {
		return false
	}

	// Check status
	statusMatches := false
	for _, s := range sf.options.StatusFilter {
		if service.Status == s {
			statusMatches = true
			break
		}
	}
	if !statusMatches {
		return false
	}

	// Check if stopped services should be shown
	if !sf.options.ShowStopped && service.Status == app.ServiceStatusStopped {
		return false
	}

	return true
}

// GetFilterSummary returns a summary of active filters
func (sf *ServiceFilter) GetFilterSummary() string {
	var parts []string

	if sf.options.SearchTerm != "" {
		parts = append(parts, "Search: "+sf.options.SearchTerm)
	}

	if len(sf.options.ServiceTypes) < 4 {
		types := make([]string, len(sf.options.ServiceTypes))
		for i, t := range sf.options.ServiceTypes {
			types[i] = string(t)
		}
		parts = append(parts, "Types: "+strings.Join(types, ","))
	}

	if len(sf.options.StatusFilter) < 4 {
		statuses := make([]string, len(sf.options.StatusFilter))
		for i, s := range sf.options.StatusFilter {
			statuses[i] = string(s)
		}
		parts = append(parts, "Status: "+strings.Join(statuses, ","))
	}

	if !sf.options.ShowStopped {
		parts = append(parts, "Hide stopped")
	}

	if len(parts) == 0 {
		return "No filters"
	}

	return strings.Join(parts, " | ")
}
