package ui

import (
	"crypto/md5"
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
	lastDetailsText string
	lastMetricsText string
	lastStatsText   string
	lastSelectedID  string
	lastUpdateTime  time.Time
	lastServiceHash string // Hash of service list to detect changes
	skipUpdates     bool   // Flag to temporarily skip updates
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

	for {
		select {
		case <-ticker.C:
			// Only update internal data - never touch UI
			d.services = d.app.GetServices()
			// UI updates only happen on user interaction or manual refresh
		}
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

// refreshLoop - DISABLED automatic updates to prevent border redraws
func (d *Dashboard) refreshLoop() {
	// Completely disable automatic refresh to prevent any border updates
	// Updates only happen on manual refresh (R key)

	// Keep the goroutine alive but do nothing
	select {}
}

// updateServiceList - lazydocker-style: NEVER modify UI structure, only internal data
func (d *Dashboard) updateServiceList() {
	if !d.staticMode {
		return // Only update in static mode
	}

	newServices := d.app.GetServices()

	// In lazydocker-style mode, we NEVER call any tview methods that could affect borders
	// We only update our internal data and let the user manually trigger UI updates

	// Just update internal data - no UI modifications whatsoever
	d.services = newServices

	// The UI will only update when user explicitly requests it via keyboard shortcuts
}

// manualRefreshUI - lazydocker-style: Only update UI when explicitly requested
func (d *Dashboard) manualRefreshUI() {
	if !d.staticMode {
		return
	}

	// First update internal data
	d.services = d.app.GetServices()

	// Then update UI in the most conservative way possible
	d.rebuildServiceListStatic()
	d.updateDetailViewStatic()
}

// rebuildServiceListStatic - lazydocker-style: minimal UI rebuild
func (d *Dashboard) rebuildServiceListStatic() {
	// Store current selection
	currentSelection := d.serviceList.GetCurrentItem()
	selectedID := ""
	if currentSelection >= 0 && currentSelection < len(d.services) {
		selectedID = d.services[currentSelection].ID
	}

	// Clear and rebuild - this is the only time we modify UI structure
	d.serviceList.Clear()

	if len(d.services) == 0 {
		d.serviceList.AddItem("[gray]No services detected", "[gray]Scanning for Docker, Systemd, and processes...", 0, nil)
		return
	}

	for i, service := range d.services {
		mainText := fmt.Sprintf("%s %s", service.GetStatusEmoji(), service.Name)
		secondaryText := d.formatSecondaryText(service)
		d.serviceList.AddItem(mainText, secondaryText, rune(48+i%10), nil)
	}

	// Restore selection by ID if possible
	if selectedID != "" {
		for i, service := range d.services {
			if service.ID == selectedID {
				d.serviceList.SetCurrentItem(i)
				d.selectedIndex = i
				break
			}
		}
	}
}

// updateDetailViewStatic - lazydocker-style: only update when explicitly called
func (d *Dashboard) updateDetailViewStatic() {
	if d.selectedIndex >= len(d.services) || len(d.services) == 0 {
		d.detailsView.SetText("\n [gray]No service selected")
		d.metricsView.SetText("\n [gray]No metrics available")
		return
	}

	service := d.services[d.selectedIndex]

	// Build details text
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

	d.detailsView.SetText(details.String())

	// Update metrics if available
	if service.Metrics != nil {
		d.updateMetricsViewStatic(service)
	} else {
		d.metricsView.SetText("\n [gray]Collecting metrics...")
	}
}

// updateMetricsViewStatic - lazydocker-style: only update when explicitly called
func (d *Dashboard) updateMetricsViewStatic(service *app.Service) {
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

	d.metricsView.SetText(metricsText.String())
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

// updateExistingServicesInPlace updates only the text content of existing items
func (d *Dashboard) updateExistingServicesInPlace(newServices []*app.Service) {
	// Create a map for quick lookup
	serviceMap := make(map[string]*app.Service)
	for _, service := range newServices {
		serviceMap[service.ID] = service
	}

	// Update only items that have changed, preserving exact same positions
	for i := 0; i < len(d.services) && i < len(newServices); i++ {
		oldService := d.services[i]

		// Try to find the service in the same position first
		var newService *app.Service
		if i < len(newServices) && newServices[i].ID == oldService.ID {
			newService = newServices[i]
		} else {
			// If not in same position, look it up by ID
			newService = serviceMap[oldService.ID]
		}

		if newService != nil && d.serviceNeedsUpdate(oldService, newService) {
			// Only update the text content - no structural changes
			mainText := fmt.Sprintf("%s %s", newService.GetStatusEmoji(), newService.Name)
			secondaryText := d.formatSecondaryText(newService)
			d.serviceList.SetItemText(i, mainText, secondaryText)
		}
	}

	// Update our internal services slice
	d.services = newServices
}

// handleServiceCountChange handles when services are added/removed with minimal disruption
func (d *Dashboard) handleServiceCountChange(newServices []*app.Service) {
	currentCount := len(d.services)
	newCount := len(newServices)

	// Store current selection to restore it
	currentSelection := d.serviceList.GetCurrentItem()
	selectedID := ""
	if currentSelection >= 0 && currentSelection < len(d.services) {
		selectedID = d.services[currentSelection].ID
	}

	// Only rebuild if absolutely necessary - when count changes significantly
	// For small changes, try to update existing items first
	if newCount > currentCount {
		// Services added - update existing ones first, then add new ones
		d.updateExistingServicesInPlace(d.services) // Update existing

		// Add new services
		for i := currentCount; i < newCount; i++ {
			if i < len(newServices) {
				service := newServices[i]
				mainText := fmt.Sprintf("%s %s", service.GetStatusEmoji(), service.Name)
				secondaryText := d.formatSecondaryText(service)
				d.serviceList.AddItem(mainText, secondaryText, rune(48+i%10), nil)
			}
		}
	} else {
		// Services removed or reordered - minimal rebuild required
		d.rebuildServiceListMinimal(newServices)
	}

	// Restore selection by ID if possible
	if selectedID != "" {
		for i, service := range newServices {
			if service.ID == selectedID {
				d.serviceList.SetCurrentItem(i)
				d.selectedIndex = i
				break
			}
		}
	}

	d.services = newServices
	d.lastUpdateTime = time.Now()
}

// serviceNeedsUpdate checks if a service needs UI update
func (d *Dashboard) serviceNeedsUpdate(old, new *app.Service) bool {
	return old.Status != new.Status ||
		old.Name != new.Name ||
		!d.slicesEqual(old.Ports, new.Ports) ||
		old.Image != new.Image
}

// slicesEqual compares two string slices
func (d *Dashboard) slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// servicesOrderChanged checks if the order of services changed
func (d *Dashboard) servicesOrderChanged(newServices []*app.Service) bool {
	if len(newServices) != len(d.services) {
		return true
	}
	for i, service := range newServices {
		if i >= len(d.services) || d.services[i].ID != service.ID {
			return true
		}
	}
	return false
}

// rebuildServiceListMinimal rebuilds the list with minimal disruption
func (d *Dashboard) rebuildServiceListMinimal(newServices []*app.Service) {
	// Store current selection
	currentSelection := d.serviceList.GetCurrentItem()
	selectedID := ""
	if currentSelection < len(d.services) && currentSelection >= 0 {
		selectedID = d.services[currentSelection].ID
	}

	// Clear and rebuild
	d.serviceList.Clear()

	for i, service := range newServices {
		mainText := fmt.Sprintf("%s %s", service.GetStatusEmoji(), service.Name)
		secondaryText := d.formatSecondaryText(service)
		d.serviceList.AddItem(mainText, secondaryText, rune(48+i%10), nil)
	}

	// Restore selection by ID if possible
	if selectedID != "" {
		for i, service := range newServices {
			if service.ID == selectedID {
				d.serviceList.SetCurrentItem(i)
				d.selectedIndex = i
				break
			}
		}
	}

	d.services = newServices
}

// formatSecondaryText formats the secondary text for a service
func (d *Dashboard) formatSecondaryText(service *app.Service) string {
	secondaryText := fmt.Sprintf("[gray]%s", service.Type)
	if len(service.Ports) > 0 {
		secondaryText += fmt.Sprintf(" • port %s", strings.Join(service.Ports, ","))
	}
	return secondaryText
}

// computeServiceHash creates a hash of current services to detect changes
func (d *Dashboard) computeServiceHash(services []*app.Service) string {
	var hashData strings.Builder
	for _, service := range services {
		hashData.WriteString(fmt.Sprintf("%s|%s|%s|%v|",
			service.ID, service.Name, service.Status, service.Ports))
	}
	hash := md5.Sum([]byte(hashData.String()))
	return fmt.Sprintf("%x", hash)
}

// updateDetailView updates only the text content of detail panels
func (d *Dashboard) updateDetailView() {
	if d.selectedIndex >= len(d.services) || len(d.services) == 0 {
		emptyText := "\n [gray]No service selected"
		emptyMetrics := "\n [gray]No metrics available"
		// Only update text content if it actually changed
		if d.lastDetailsText != emptyText {
			d.detailsView.SetText(emptyText)
			d.lastDetailsText = emptyText
		}
		if d.lastMetricsText != emptyMetrics {
			d.metricsView.SetText(emptyMetrics)
			d.lastMetricsText = emptyMetrics
		}
		return
	}

	service := d.services[d.selectedIndex]

	// Only rebuild details if selected service changed
	if d.lastSelectedID != service.ID {
		d.lastSelectedID = service.ID

		// Build details text
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

		if len(service.Warnings) > 0 {
			details.WriteString("\n [red]⚠ Warnings:[white]\n")
			for _, warning := range service.Warnings {
				details.WriteString(fmt.Sprintf("   • %s\n", warning))
			}
		}

		if len(service.Environment) > 0 {
			details.WriteString("\n [yellow]Environment:[white]\n")
			count := 0
			for key, value := range service.Environment {
				if count >= 5 {
					details.WriteString(fmt.Sprintf("   ... and %d more\n", len(service.Environment)-5))
					break
				}
				// Truncate long values
				if len(value) > 45 {
					value = value[:45] + "..."
				}
				details.WriteString(fmt.Sprintf("   • %s: %s\n", key, value))
				count++
			}
		}

		detailsText := details.String()
		if d.lastDetailsText != detailsText {
			d.detailsView.SetText(detailsText)
			d.lastDetailsText = detailsText
		}
	}

	// Update metrics only if they changed or are new
	if service.Metrics != nil {
		d.updateMetricsView(service)
	} else {
		noMetrics := "\n [gray]Collecting metrics..."
		if d.lastMetricsText != noMetrics {
			d.metricsView.SetText(noMetrics)
			d.lastMetricsText = noMetrics
		}
	}
}

// updateMetricsView updates the metrics view
func (d *Dashboard) updateMetricsView(service *app.Service) {
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

	// Disk I/O
	if metrics.DiskIORead > 0 || metrics.DiskIOWrite > 0 {
		metricsText.WriteString(fmt.Sprintf(" [yellow]Disk I/O:[white] R: %s  W: %s\n",
			components.FormatBytes(metrics.DiskIORead),
			components.FormatBytes(metrics.DiskIOWrite)))
	}

	// Uptime
	if metrics.Uptime > 0 {
		metricsText.WriteString(fmt.Sprintf(" [yellow]Uptime:[white]  %s\n",
			components.FormatDuration(int64(metrics.Uptime.Seconds()))))
	}

	// History sparklines
	if metrics.History != nil && len(metrics.History.CPU) > 3 {
		metricsText.WriteString("\n [yellow]History (5min):[white]\n")
		cpuSparkline := components.CreateSparkline(metrics.History.CPU, 28)
		metricsText.WriteString(fmt.Sprintf("  CPU:    %s\n", cpuSparkline))

		if len(metrics.History.Memory) > 3 {
			memSparkline := components.CreateSparkline(metrics.History.Memory, 28)
			metricsText.WriteString(fmt.Sprintf("  Memory: %s\n", memSparkline))
		}
	}

	newMetricsText := metricsText.String()
	// Only update if metrics text changed
	if d.lastMetricsText != newMetricsText {
		d.metricsView.SetText(newMetricsText)
		d.lastMetricsText = newMetricsText
	}
}

// Stop stops the dashboard
func (d *Dashboard) Stop() {
	d.tviewApp.Stop()
}
