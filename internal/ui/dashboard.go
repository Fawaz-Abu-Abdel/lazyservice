package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/lazyservice/lazyservice/internal/app"
	"github.com/lazyservice/lazyservice/internal/ui/components"
	"github.com/rivo/tview"
)

// Dashboard represents the main TUI dashboard with lazydocker-style static borders
type Dashboard struct {
	app             *app.App
	tviewApp        *tview.Application
	serviceList     *tview.List
	detailsView     *tview.TextView
	metricsView     *tview.TextView
	statisticsView  *tview.TextView
	statusBar       *tview.TextView
	rootLayout      tview.Primitive
	selectedIndex   int
	services        []*app.Service
	lastStatsText   string
	staticMode      bool   // Lazydocker-style static mode
	bordersLocked   bool   // Lock borders from any modifications
	contentOnly     bool   // Only update content, never structure
	showStatistics  bool   // Toggle between metrics and statistics view
}

// NewDashboard creates a new dashboard with lazydocker-style static borders
func NewDashboard(application *app.App) *Dashboard {
	return &Dashboard{
		app:           application,
		tviewApp:      tview.NewApplication(),
		selectedIndex: 0,
		services:      make([]*app.Service, 0),
		staticMode:    true, // Enable lazydocker-style static mode
		bordersLocked: true, // Lock borders completely
		contentOnly:   true, // Only update content, never structure
	}
}

// Run starts the dashboard with lazydocker-style border protection
func (d *Dashboard) Run() error {
	// Create UI components ONCE - borders are now locked
	d.createUI()

	// Lock borders permanently after creation
	d.bordersLocked = true

	// Initial data load using content-only updates
	d.loadInitialData()

	// Start background data collection (no UI updates)
	go d.backgroundDataLoop()

	// Run the application
	return d.tviewApp.Run()
}

// loadInitialData loads initial data without touching borders
func (d *Dashboard) loadInitialData() {
	d.services = d.app.GetServices()
	d.populateInitialContent()
}

// populateInitialContent fills the UI with initial content (borders already locked)
func (d *Dashboard) populateInitialContent() {
	if len(d.services) == 0 {
		d.serviceList.AddItem("[gray]No services detected", "[gray]Scanning for Docker, Systemd, and processes...", 0, nil)
		d.detailsView.SetText("\n [gray]No service selected")
		d.metricsView.SetText("\n [gray]No metrics available")
		return
	}

	// Populate service list once
	for i, service := range d.services {
		mainText := fmt.Sprintf("%s %s", service.GetStatusEmoji(), service.Name)
		secondaryText := d.formatSecondaryText(service)
		d.serviceList.AddItem(mainText, secondaryText, rune(48+i%10), nil)
	}

	// Set initial selection
	if len(d.services) > 0 {
		d.serviceList.SetCurrentItem(0)
		d.selectedIndex = 0
		d.updateContentOnly()
	}
}

// backgroundDataLoop continuously updates internal data without touching UI
func (d *Dashboard) backgroundDataLoop() {
	ticker := time.NewTicker(10 * time.Second) // Reduced frequency
	defer ticker.Stop()

	for range ticker.C {
		// Only update internal data - never touch UI
		d.services = d.app.GetServices()
		// UI updates only happen on user interaction or manual refresh
	}
}

// createUI creates the UI layout
func (d *Dashboard) createUI() {
	// Service list (left panel)
	d.serviceList = tview.NewList()
	d.serviceList.SetBorder(true)
	d.serviceList.SetTitle(" Services (0) ")
	d.serviceList.ShowSecondaryText(true)
	d.serviceList.SetBorderPadding(0, 0, 1, 1)

	d.serviceList.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		d.selectedIndex = index
		// In content-only mode, update only the text content
		if d.contentOnly {
			d.updateContentOnly()
		}
	})

	// Details view (top right panel)
	d.detailsView = tview.NewTextView()
	d.detailsView.SetBorder(true)
	d.detailsView.SetTitle(" Service Details ")
	d.detailsView.SetDynamicColors(true)
	d.detailsView.SetScrollable(true)
	d.detailsView.SetBorderPadding(0, 0, 1, 1)
	d.detailsView.SetWordWrap(true)

	// Metrics view (bottom right panel)
	d.metricsView = tview.NewTextView()
	d.metricsView.SetBorder(true)
	d.metricsView.SetTitle(" Resource Usage ")
	d.metricsView.SetDynamicColors(true)
	d.metricsView.SetBorderPadding(0, 0, 1, 1)

	// Statistics view (alternative to metrics view)
	d.statisticsView = tview.NewTextView()
	d.statisticsView.SetBorder(true)
	d.statisticsView.SetTitle(" Service Statistics ")
	d.statisticsView.SetDynamicColors(true)
	d.statisticsView.SetScrollable(true)
	d.statisticsView.SetBorderPadding(0, 0, 1, 1)

	// Status bar (bottom)
	d.statusBar = tview.NewTextView()
	d.statusBar.SetDynamicColors(true)
	d.statusBar.SetTextAlign(tview.AlignCenter)
	d.statusBar.SetText("[yellow]Metrics Mode:[white] [T]oggle Statistics [R]efresh [H]elp [Q]uit | [gray]Borders NEVER affected")

	// Create bottom panel with both metrics and statistics views side by side
	bottomPanel := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(d.metricsView, 0, 1, false).
		AddItem(d.statisticsView, 0, 1, false)

	// Create right panel layout
	rightPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.detailsView, 0, 1, false).
		AddItem(bottomPanel, 15, 0, false)

	// Main layout
	mainLayout := tview.NewFlex().
		AddItem(d.serviceList, 0, 1, true).
		AddItem(rightPanel, 0, 2, false)

	// Root layout with status bar
	rootLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainLayout, 0, 1, true).
		AddItem(d.statusBar, 1, 0, false)

	// Set up enhanced input handling
	d.setupInputHandling()

	d.rootLayout = rootLayout
	d.tviewApp.SetRoot(rootLayout, true)
}


// updateContentOnly - NEVER touches borders, only updates text content
func (d *Dashboard) updateContentOnly() {
	if !d.contentOnly || d.selectedIndex >= len(d.services) || len(d.services) == 0 {
		return
	}

	service := d.services[d.selectedIndex]

	// Update detail content only - no border operations
	var details strings.Builder
	details.WriteString("\n")
	details.WriteString(fmt.Sprintf(" [yellow]Name:[white]   %s\n", service.Name))
	details.WriteString(fmt.Sprintf(" [yellow]Type:[white]   %s\n", service.Type))
	details.WriteString(fmt.Sprintf(" [yellow]Status:[white] %s %s\n", service.GetStatusEmoji(), service.Status))
	details.WriteString(fmt.Sprintf(" [yellow]ID:[white]     %s\n", service.ID))

	if service.Image != "" {
		details.WriteString(fmt.Sprintf(" [yellow]Image:[white]  %s\n", service.Image))
	}
	if service.Version != "" {
		details.WriteString(fmt.Sprintf(" [yellow]Version:[white] %s\n", service.Version))
	}
	if service.Namespace != "" {
		details.WriteString(fmt.Sprintf(" [yellow]Namespace:[white] %s\n", service.Namespace))
	}
	if len(service.Ports) > 0 {
		details.WriteString(fmt.Sprintf(" [yellow]Ports:[white]  %s\n", strings.Join(service.Ports, ", ")))
	}
	if service.HealthCheck != "" {
		details.WriteString(fmt.Sprintf(" [yellow]Health:[white] %s\n", service.HealthCheck))
	}

	details.WriteString(fmt.Sprintf("\n [yellow]Created:[white] %s\n", service.CreatedAt.Format("2006-01-02 15:04:05")))

	// Only SetText - never any border operations
	d.detailsView.SetText(details.String())

	// Update metrics or statistics content based on current mode
	if d.showStatistics {
		// Show statistics view
		statsText := FormatStatistics(service)
		d.statisticsView.SetText(statsText)
		d.lastStatsText = statsText
		
		// Show quick metrics in metrics view
		if service.Metrics != nil {
			quickMetrics := FormatQuickStats(service)
			d.metricsView.SetText(quickMetrics)
		} else {
			d.metricsView.SetText("\n [gray]Collecting metrics...")
		}
	} else {
		// Show metrics view
		if service.Metrics != nil {
			d.updateMetricsContentOnly(service)
		} else {
			d.metricsView.SetText("\n [gray]Collecting metrics...")
		}
		// Show quick stats in statistics view
		quickStats := FormatQuickStats(service)
		d.statisticsView.SetText(quickStats)
		d.lastStatsText = quickStats
	}
}

// updateMetricsContentOnly - NEVER touches borders, only updates metrics text
func (d *Dashboard) updateMetricsContentOnly(service *app.Service) {
	metrics := service.Metrics
	var metricsText strings.Builder

	metricsText.WriteString("\n")

	// CPU
	cpuBar := components.CreateBarChart(metrics.CPUPercent, 100, 18)
	metricsText.WriteString(fmt.Sprintf(" [yellow]CPU:[white]     %5.1f%% %s\n", metrics.CPUPercent, cpuBar))

	// Memory
	if metrics.MemoryLimit > 0 {
		memBar := components.CreateBarChart(metrics.MemoryPercent, 100, 18)
		metricsText.WriteString(fmt.Sprintf(" [yellow]Memory:[white]  %5.1f%% %s\n", metrics.MemoryPercent, memBar))
		metricsText.WriteString(fmt.Sprintf("           %s / %s\n",
			components.FormatBytes(metrics.MemoryUsage),
			components.FormatBytes(metrics.MemoryLimit)))
	} else {
		metricsText.WriteString(fmt.Sprintf(" [yellow]Memory:[white]  %s\n",
			components.FormatBytes(metrics.MemoryUsage)))
	}

	// Network
	if metrics.NetworkIn > 0 || metrics.NetworkOut > 0 {
		metricsText.WriteString(fmt.Sprintf(" [yellow]Network:[white] ↓ %s  ↑ %s\n",
			components.FormatBytes(metrics.NetworkIn),
			components.FormatBytes(metrics.NetworkOut)))
	}

	// Uptime
	if metrics.Uptime > 0 {
		metricsText.WriteString(fmt.Sprintf(" [yellow]Uptime:[white]  %s\n",
			components.FormatDuration(int64(metrics.Uptime.Seconds()))))
	}

	// Only SetText - never any border operations
	d.metricsView.SetText(metricsText.String())
}

// toggleStatisticsView toggles between metrics and statistics view
func (d *Dashboard) toggleStatisticsView() {
	if !d.contentOnly {
		return
	}
	
	d.showStatistics = !d.showStatistics
	
	// Update the content based on current mode
	d.updateContentOnly()
	
	// Update status bar to show current mode
	if d.showStatistics {
		d.statusBar.SetText("[yellow]Statistics Mode:[white] [T]oggle to Metrics [R]efresh [H]elp [Q]uit")
	} else {
		d.statusBar.SetText("[yellow]Metrics Mode:[white] [T]oggle to Statistics [R]efresh [H]elp [Q]uit")
	}
}

// contentOnlyRefresh - Manual refresh that only updates content, never borders
func (d *Dashboard) contentOnlyRefresh() {
	if !d.contentOnly {
		return
	}

	// Get fresh data
	d.services = d.app.GetServices()

	// Update service list content only - use SetItemText for existing items
	for i, service := range d.services {
		if i < d.serviceList.GetItemCount() {
			mainText := fmt.Sprintf("%s %s", service.GetStatusEmoji(), service.Name)
			secondaryText := d.formatSecondaryText(service)
			d.serviceList.SetItemText(i, mainText, secondaryText)
		}
	}

	// Update detail content
	d.updateContentOnly()
}


// formatSecondaryText formats the secondary text for a service
func (d *Dashboard) formatSecondaryText(service *app.Service) string {
	secondaryText := fmt.Sprintf("[gray]%s", service.Type)
	if len(service.Ports) > 0 {
		secondaryText += fmt.Sprintf(" • port %s", strings.Join(service.Ports, ","))
	}
	return secondaryText
}


// Stop stops the dashboard
func (d *Dashboard) Stop() {
	d.tviewApp.Stop()
}
