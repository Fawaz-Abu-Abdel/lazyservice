//go:build !windows

package ui

import (
	"fmt"
	"github.com/lazyservice/lazyservice/internal/app"
)

// MinimalDashboard is a stub for non-Windows platforms
type MinimalDashboard struct {
}

// NewMinimalDashboard creates a new stub dashboard
func NewMinimalDashboard(application *app.App) *MinimalDashboard {
	return &MinimalDashboard{}
}

// Run returns an error on non-Windows platforms
func (d *MinimalDashboard) Run() error {
	return fmt.Errorf("minimal dashboard is only supported on Windows")
}

// Stop is a no-op stub
func (d *MinimalDashboard) Stop() {
}
