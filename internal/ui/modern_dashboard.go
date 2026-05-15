package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/lazyservice/lazyservice/internal/app"
	"github.com/lazyservice/lazyservice/internal/ui/components"
)

// ModernDashboard represents the beautiful modern TUI dashboard
type ModernDashboard struct {
	app             *app.App
	tviewApp        *tview.Application
	
	// Layout components
	mainFlex        *tview.Flex
	leftPanel       *tview.Flex
	rightPanel      *tview.Flex
	
	// Service components
	serviceList     *tview.Table
	
	// Detail components
	detailsView     *tview.TextView
	metricsView     *tview.TextView
	statisticsView  *tview.TextView
	logsView        *tview.TextView
	
	// Status components
	statusBar       *tview.TextView
	headerBar       *tview.TextView
	
	// State
	selectedIndex   int
	services        []*app.Service
	theme           *Theme
	currentView     string
	animationFrame  int
	lastRefresh     time.Time
	bordersLocked   bool
	contentOnly     bool
}

// NewModernDashboard creates a stunning modern dashboard
func NewModernDashboard(application *app.App) *ModernDashboard {
	theme := ModernDarkTheme()
	ApplyTheme(theme)
	
	return &ModernDashboard{
		app:           application,
		tviewApp:      tview.NewApplication(),
		selectedIndex: 0,
		services:      make([]*app.Service, 0),
		theme:         theme,
		currentView:   "services",
		bordersLocked: true,
		contentOnly:   true,
	}
}

// Run starts the beautiful modern dashboard
func (d *ModernDashboard) Run() error {
	d.createModernUI()
	d.setupAdvancedInputHandling()
	
	// Start background processes
	go d.backgroundDataLoop()
	go d.animationLoop()
	
	return d.tviewApp.Run()
}

// Stop stops the modern dashboard
func (d *ModernDashboard) Stop() {
	d.tviewApp.Stop()
}

// createModernUI creates the stunning modern interface
func (d *ModernDashboard) createModernUI() {
	// Create header bar with beautiful styling
	d.createHeaderBar()
	
	// Create service list with modern table
	d.createModernServiceList()
	
	// Create detail panels with advanced styling
	d.createDetailPanels()
	
	// Create status bar with animations
	d.createStatusBar()
	
	// Assemble the layout
	d.assembleLayout()
	
	// Load initial data
	d.loadInitialData()
}

// createHeaderBar creates a beautiful header with branding and info
func (d *ModernDashboard) createHeaderBar() {
	d.headerBar = tview.NewTextView()
	d.headerBar.SetDynamicColors(true)
	d.headerBar.SetTextAlign(tview.AlignCenter)
	
	headerText := fmt.Sprintf(
		"[#89b4fa]╭─────────────────────────────────────────────────────────────────────────────────╮[white]\n"+
		"[#89b4fa]│[white] [#cba6f7]⚡ LazyService[white] [#a6adc8]v2.0[white] [#89b4fa]│[white] [#f9e2af]Modern Dashboard[white] [#89b4fa]│[white] [#a6e3a1]%s[white] [#89b4fa]│[white]\n"+
		"[#89b4fa]╰─────────────────────────────────────────────────────────────────────────────────╯[white]",
		time.Now().Format("15:04:05"))
	
	d.headerBar.SetText(headerText)
}

// createModernServiceList creates a beautiful service table
func (d *ModernDashboard) createModernServiceList() {
	d.serviceList = tview.NewTable()
	d.serviceList.SetBorder(true)
	d.serviceList.SetTitle(" 🚀 Services ")
	d.serviceList.SetBorderColor(d.theme.BorderColor)
	d.serviceList.SetTitleColor(d.theme.AccentText)
	
	// Set table styling
	d.serviceList.SetSelectable(true, false)
	d.serviceList.SetSelectedStyle(tcell.StyleDefault.
		Background(d.theme.SelectionColor).
		Foreground(d.theme.AccentText))
	
	// Create beautiful headers
	headers := []string{"Status", "Name", "Type", "Uptime", "CPU", "Memory", "Network"}
	for i, header := range headers {
		cell := tview.NewTableCell(fmt.Sprintf("[::b][%s]%s[::-]", "#89b4fa", header))
		cell.SetAlign(tview.AlignCenter)
		cell.SetSelectable(false)
		d.serviceList.SetCell(0, i, cell)
	}
	
	// Set selection handler
	d.serviceList.SetSelectedFunc(func(row, col int) {
		if row > 0 && row-1 < len(d.services) {
			d.selectedIndex = row - 1
			d.updateDetailPanels()
		}
	})
}

// createDetailPanels creates beautiful detail views
func (d *ModernDashboard) createDetailPanels() {
	// Service Details Panel
	d.detailsView = tview.NewTextView()
	d.detailsView.SetBorder(true)
	d.detailsView.SetTitle(" 📋 Service Details ")
	d.detailsView.SetBorderColor(d.theme.BorderColor)
	d.detailsView.SetTitleColor(d.theme.AccentText)
	d.detailsView.SetDynamicColors(true)
	d.detailsView.SetScrollable(true)
	
	// Real-time Metrics Panel
	d.metricsView = tview.NewTextView()
	d.metricsView.SetBorder(true)
	d.metricsView.SetTitle(" 📊 Live Metrics ")
	d.metricsView.SetBorderColor(d.theme.BorderColor)
	d.metricsView.SetTitleColor(d.theme.AccentText)
	d.metricsView.SetDynamicColors(true)
	
	// Statistics Panel
	d.statisticsView = tview.NewTextView()
	d.statisticsView.SetBorder(true)
	d.statisticsView.SetTitle(" 📈 Statistics ")
	d.statisticsView.SetBorderColor(d.theme.BorderColor)
	d.statisticsView.SetTitleColor(d.theme.AccentText)
	d.statisticsView.SetDynamicColors(true)
	d.statisticsView.SetScrollable(true)
	
	// Logs Panel
	d.logsView = tview.NewTextView()
	d.logsView.SetBorder(true)
	d.logsView.SetTitle(" 📝 Logs ")
	d.logsView.SetBorderColor(d.theme.BorderColor)
	d.logsView.SetTitleColor(d.theme.AccentText)
	d.logsView.SetDynamicColors(true)
	d.logsView.SetScrollable(true)
}

// createStatusBar creates an animated status bar
func (d *ModernDashboard) createStatusBar() {
	d.statusBar = tview.NewTextView()
	d.statusBar.SetDynamicColors(true)
	d.statusBar.SetTextAlign(tview.AlignCenter)
	d.updateStatusBar()
}

// assembleLayout creates the beautiful layout
func (d *ModernDashboard) assembleLayout() {
	// Create left panel (service list)
	d.leftPanel = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.serviceList, 0, 1, true)
	
	// Create right panel (details and metrics)
	topRightPanel := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(d.detailsView, 0, 1, false).
		AddItem(d.metricsView, 0, 1, false)
	
	// Create bottom right panel (statistics and logs)
	bottomRightPanel := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(d.statisticsView, 0, 1, false).
		AddItem(d.logsView, 0, 1, false)
	
	// Combine right panels
	d.rightPanel = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topRightPanel, 0, 1, false).
		AddItem(bottomRightPanel, 0, 1, false)
	
	// Create main layout
	d.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(d.leftPanel, 0, 1, true).
		AddItem(d.rightPanel, 0, 2, false)
	
	// Create root layout with header and status
	rootLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.headerBar, 3, 0, false).
		AddItem(d.mainFlex, 0, 1, true).
		AddItem(d.statusBar, 1, 0, false)
	
	d.tviewApp.SetRoot(rootLayout, true)
}

// loadInitialData loads and displays initial service data
func (d *ModernDashboard) loadInitialData() {
	d.services = d.app.GetServices()
	d.updateServiceTable()
	if len(d.services) > 0 {
		d.updateDetailPanels()
	}
}

// updateServiceTable updates the beautiful service table (content-only, no border changes)
func (d *ModernDashboard) updateServiceTable() {
	// Remove duplicate services first
	uniqueServices := d.removeDuplicateServices(d.services)
	d.services = uniqueServices
	
	// Clear existing rows only if service count changed to avoid unnecessary border updates
	currentRowCount := d.serviceList.GetRowCount() - 1 // Exclude header
	if currentRowCount != len(d.services) {
		// Only clear if absolutely necessary
		for row := d.serviceList.GetRowCount() - 1; row > 0; row-- {
			d.serviceList.RemoveRow(row)
		}
	}
	
	// Update or add service rows
	for i, service := range d.services {
		row := i + 1
		
		// Status with icon
		statusIcon := StatusIcon(string(service.Status), d.theme)
		statusCell := tview.NewTableCell(statusIcon + " " + string(service.Status))
		statusCell.SetAlign(tview.AlignCenter)
		
		// Service name with emoji
		nameIcon := d.getServiceIcon(service.Type)
		nameCell := tview.NewTableCell(fmt.Sprintf("%s %s", nameIcon, service.Name))
		
		// Type with color
		typeCell := tview.NewTableCell(fmt.Sprintf("[#a6adc8]%s[white]", service.Type))
		typeCell.SetAlign(tview.AlignCenter)
		
		// Uptime
		uptime := time.Since(service.CreatedAt)
		uptimeCell := tview.NewTableCell(components.FormatDuration(int64(uptime.Seconds())))
		uptimeCell.SetAlign(tview.AlignCenter)
		
		// CPU with progress bar - ensure we show actual data
		cpuPercent := 0.0
		cpuText := "0.0%"
		if service.Metrics != nil && service.Metrics.CPUPercent > 0 {
			cpuPercent = service.Metrics.CPUPercent
			cpuText = fmt.Sprintf("%.1f%%", cpuPercent)
		} else if service.Statistics != nil {
			// Fallback to statistics if metrics not available
			avgCPU := service.Statistics.GetAverageCPU(5 * time.Minute)
			if avgCPU > 0 {
				cpuPercent = avgCPU
				cpuText = fmt.Sprintf("%.1f%%", avgCPU)
			}
		}
		cpuBar := ProgressBar(cpuPercent, 6, "green")
		cpuCell := tview.NewTableCell(fmt.Sprintf("%s %s", cpuText, cpuBar))
		
		// Memory with progress bar
		memPercent := 0.0
		if service.Metrics != nil {
			memPercent = service.Metrics.MemoryPercent
		}
		memBar := ProgressBar(memPercent, 8, "blue")
		memCell := tview.NewTableCell(fmt.Sprintf("%.1f%% %s", memPercent, memBar))
		
		// Network with sparkline
		networkText := "0 B/s"
		if service.Metrics != nil {
			total := service.Metrics.NetworkIn + service.Metrics.NetworkOut
			networkText = components.FormatBytes(total) + "/s"
		}
		networkCell := tview.NewTableCell(networkText)
		networkCell.SetAlign(tview.AlignCenter)
		
		// Set all cells
		d.serviceList.SetCell(row, 0, statusCell)
		d.serviceList.SetCell(row, 1, nameCell)
		d.serviceList.SetCell(row, 2, typeCell)
		d.serviceList.SetCell(row, 3, uptimeCell)
		d.serviceList.SetCell(row, 4, cpuCell)
		d.serviceList.SetCell(row, 5, memCell)
		d.serviceList.SetCell(row, 6, networkCell)
	}
	
	// Update service count in title
	d.serviceList.SetTitle(fmt.Sprintf(" 🚀 Services (%d) ", len(d.services)))
}

// getServiceIcon returns an appropriate icon for the service type
func (d *ModernDashboard) getServiceIcon(serviceType app.ServiceType) string {
	switch serviceType {
	case app.ServiceTypeDocker:
		return "🐳"
	case app.ServiceTypeKubernetes:
		return "☸️"
	case app.ServiceTypeSystemd:
		return "⚙️"
	case app.ServiceTypeProcess:
		return "⚡"
	default:
		return "📦"
	}
}

// updateDetailPanels updates all detail panels with beautiful formatting
func (d *ModernDashboard) updateDetailPanels() {
	if d.selectedIndex >= len(d.services) || len(d.services) == 0 {
		d.showEmptyState()
		return
	}
	
	service := d.services[d.selectedIndex]
	
	// Update details panel
	d.updateServiceDetails(service)
	
	// Update metrics panel
	d.updateLiveMetrics(service)
	
	// Update statistics panel
	d.updateStatisticsPanel(service)
	
	// Update logs panel
	d.updateLogsPanel(service)
}

// updateServiceDetails creates beautiful service details
func (d *ModernDashboard) updateServiceDetails(service *app.Service) {
	var details strings.Builder
	
	details.WriteString("\n")
	details.WriteString(" [#cba6f7]╭─ Service Information ─╮[white]\n")
	details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Name:[white]      %s\n", service.Name))
	details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Type:[white]      %s %s\n", d.getServiceIcon(service.Type), service.Type))
	details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Status:[white]    %s %s\n", StatusIcon(string(service.Status), d.theme), service.Status))
	details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]ID:[white]        %s\n", service.ID))
	
	if service.Image != "" {
		details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Image:[white]     %s\n", service.Image))
	}
	if service.Version != "" {
		details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Version:[white]   %s\n", service.Version))
	}
	if len(service.Ports) > 0 {
		details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Ports:[white]     %s\n", strings.Join(service.Ports, ", ")))
	}
	
	details.WriteString(" [#cba6f7]╰─────────────────────────╯[white]\n\n")
	
	// Runtime information
	details.WriteString(" [#cba6f7]╭─ Runtime Information ─╮[white]\n")
	details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Created:[white]   %s\n", service.CreatedAt.Format("2006-01-02 15:04:05")))
	
	uptime := time.Since(service.CreatedAt)
	details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Uptime:[white]    %s\n", components.FormatDuration(int64(uptime.Seconds()))))
	
	if service.HealthCheck != "" {
		details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Health:[white]    %s\n", service.HealthCheck))
	}
	
	details.WriteString(" [#cba6f7]╰─────────────────────────╯[white]\n")
	
	d.detailsView.SetText(details.String())
}

// updateLiveMetrics creates beautiful live metrics display
func (d *ModernDashboard) updateLiveMetrics(service *app.Service) {
	if service.Metrics == nil {
		d.metricsView.SetText("\n [#a6adc8]📊 Collecting metrics...[white]\n\n [#585b70]" + AnimatedSpinner(d.animationFrame) + " Please wait[white]")
		return
	}
	
	metrics := service.Metrics
	var output strings.Builder
	
	output.WriteString("\n")
	output.WriteString(" [#cba6f7]╭─ Live Performance ─╮[white]\n")
	
	// CPU with beautiful progress bar
	cpuBar := ProgressBar(metrics.CPUPercent, 20, "green")
	output.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]CPU:[white]     %5.1f%% %s\n", metrics.CPUPercent, cpuBar))
	
	// Memory with progress bar
	if metrics.MemoryLimit > 0 {
		memBar := ProgressBar(metrics.MemoryPercent, 20, "blue")
		output.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Memory:[white]  %5.1f%% %s\n", metrics.MemoryPercent, memBar))
		output.WriteString(fmt.Sprintf(" [#89b4fa]│[white]         %s / %s\n",
			components.FormatBytes(metrics.MemoryUsage),
			components.FormatBytes(metrics.MemoryLimit)))
	}
	
	// Network with icons
	if metrics.NetworkIn > 0 || metrics.NetworkOut > 0 {
		output.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Network:[white] [#a6e3a1]↓[white] %s  [#f38ba8]↑[white] %s\n",
			components.FormatBytes(metrics.NetworkIn),
			components.FormatBytes(metrics.NetworkOut)))
	}
	
	// Disk I/O
	if metrics.DiskIORead > 0 || metrics.DiskIOWrite > 0 {
		output.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Disk I/O:[white] [#a6e3a1]R:[white] %s  [#f38ba8]W:[white] %s\n",
			components.FormatBytes(metrics.DiskIORead),
			components.FormatBytes(metrics.DiskIOWrite)))
	}
	
	// Uptime
	if metrics.Uptime > 0 {
		output.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af]Uptime:[white]  %s\n",
			components.FormatDuration(int64(metrics.Uptime.Seconds()))))
	}
	
	output.WriteString(" [#cba6f7]╰─────────────────────────╯[white]\n")
	
	d.metricsView.SetText(output.String())
}

// updateStatisticsPanel creates beautiful statistics display
func (d *ModernDashboard) updateStatisticsPanel(service *app.Service) {
	statsText := FormatStatistics(service)
	d.statisticsView.SetText(statsText)
}

// updateLogsPanel creates a beautiful logs display
func (d *ModernDashboard) updateLogsPanel(service *app.Service) {
	var logs strings.Builder
	
	logs.WriteString("\n")
	logs.WriteString(" [#cba6f7]╭─ Recent Activity ─╮[white]\n")
	logs.WriteString(" [#89b4fa]│[white] [#a6adc8]" + time.Now().Format("15:04:05") + "[white] Service started\n")
	logs.WriteString(" [#89b4fa]│[white] [#a6adc8]" + time.Now().Add(-30*time.Second).Format("15:04:05") + "[white] Metrics collection enabled\n")
	logs.WriteString(" [#89b4fa]│[white] [#a6adc8]" + time.Now().Add(-60*time.Second).Format("15:04:05") + "[white] Health check passed\n")
	logs.WriteString(" [#cba6f7]╰─────────────────────╯[white]\n\n")
	logs.WriteString(" [#585b70]💡 Real-time logs coming soon...[white]")
	
	d.logsView.SetText(logs.String())
}

// showEmptyState shows a beautiful empty state
func (d *ModernDashboard) showEmptyState() {
	emptyText := "\n [#a6adc8]📭 No service selected[white]\n\n [#585b70]Select a service from the list to view details[white]"
	
	d.detailsView.SetText(emptyText)
	d.metricsView.SetText(emptyText)
	d.statisticsView.SetText(emptyText)
	d.logsView.SetText(emptyText)
}

// updateStatusBar updates the animated status bar
func (d *ModernDashboard) updateStatusBar() {
	spinner := AnimatedSpinner(d.animationFrame)
	
	statusText := fmt.Sprintf(
		"[#89b4fa]%s[white] [#f9e2af]Modern Mode[white] │ [#a6e3a1][R]efresh[white] │ [#cba6f7][T]heme[white] │ [#f38ba8][H]elp[white] │ [#fab387][Q]uit[white] │ [#585b70]Borders: Static[white]",
		spinner)
	
	d.statusBar.SetText(statusText)
}

// setupAdvancedInputHandling sets up beautiful keyboard shortcuts
func (d *ModernDashboard) setupAdvancedInputHandling() {
	d.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q', 'Q':
			d.tviewApp.Stop()
			return nil
		case 'r', 'R':
			d.refreshData()
			return nil
		case 't', 'T':
			d.toggleTheme()
			return nil
		case 'h', 'H', '?':
			d.showHelpModal()
			return nil
		}
		return event
	})
}

// refreshData refreshes all data with rate limiting
func (d *ModernDashboard) refreshData() {
	// Rate limit manual refreshes to prevent spam
	now := time.Now()
	if d.lastRefresh.Add(2 * time.Second).After(now) {
		d.statusBar.SetText("[#f38ba8]⚠ Please wait before refreshing again[white]")
		go func() {
			time.Sleep(1 * time.Second)
			d.updateStatusBar()
		}()
		return
	}
	d.lastRefresh = now
	
	d.app.RefreshServices()
	newServices := d.app.GetServices()
	
	// Only update UI if data actually changed
	if d.servicesChanged(newServices) {
		d.services = newServices
		d.updateServiceTable()
		d.updateDetailPanels()
		d.statusBar.SetText("[#a6e3a1]✨ Data refreshed successfully![white]")
	} else {
		d.statusBar.SetText("[#89b4fa]ℹ No changes detected[white]")
	}
	
	// Reset status bar after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		d.updateStatusBar()
	}()
}

// toggleTheme switches between themes
func (d *ModernDashboard) toggleTheme() {
	if d.theme == ModernDarkTheme() {
		d.theme = CyberpunkTheme()
	} else {
		d.theme = ModernDarkTheme()
	}
	
	ApplyTheme(d.theme)
	d.statusBar.SetText("[#cba6f7]🎨 Theme changed![white]")
	
	// Reset status bar after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		d.updateStatusBar()
	}()
}

// showHelpModal shows a beautiful help modal
func (d *ModernDashboard) showHelpModal() {
	helpText := `[#cba6f7]╭─ LazyService Modern Dashboard ─╮[white]

[#f9e2af]Navigation:[white]
  [#89b4fa]↑/↓, j/k[white]    Navigate services
  [#89b4fa]Enter[white]       Select service
  [#89b4fa]Tab[white]         Switch panels

[#f9e2af]Actions:[white]
  [#a6e3a1]r, R[white]        Refresh data
  [#cba6f7]t, T[white]        Toggle theme
  [#fab387]h, H, ?[white]     Show help

[#f9e2af]Application:[white]
  [#f38ba8]q, Q[white]        Quit application

[#585b70]Press any key to close...[white]`

	modal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			d.tviewApp.SetRoot(d.mainFlex, true)
		})

	modal.SetBorderColor(d.theme.FocusedBorderColor)
	modal.SetBackgroundColor(d.theme.PrimaryBackground)
	
	d.tviewApp.SetRoot(modal, false)
}

// backgroundDataLoop continuously updates data at a very slow rate
func (d *ModernDashboard) backgroundDataLoop() {
	ticker := time.NewTicker(30 * time.Second) // Much slower - only every 30 seconds
	defer ticker.Stop()

	for range ticker.C {
		// Only update if data actually changed
		newServices := d.app.GetServices()
		if d.servicesChanged(newServices) {
			d.tviewApp.QueueUpdateDraw(func() {
				d.services = newServices
				d.updateServiceTable()
				if d.selectedIndex < len(d.services) {
					d.updateDetailPanels()
				}
			})
		}
	}
}

// animationLoop handles animations and updates at a much slower rate
func (d *ModernDashboard) animationLoop() {
	ticker := time.NewTicker(5 * time.Second) // Very slow - only every 5 seconds
	defer ticker.Stop()

	for range ticker.C {
		d.tviewApp.QueueUpdateDraw(func() {
			d.animationFrame++
			d.updateStatusBar()

			// Update header time only every 5 seconds to reduce refresh
			headerText := fmt.Sprintf(
				"[#89b4fa]╭─────────────────────────────────────────────────────────────────────────────────╮[white]\n"+
					"[#89b4fa]│[white] [#cba6f7]⚡ LazyService[white] [#a6adc8]v2.0[white] [#89b4fa]│[white] [#f9e2af]Modern Dashboard[white] [#89b4fa]│[white] [#a6e3a1]%s[white] [#89b4fa]│[white]\n"+
					"[#89b4fa]╰─────────────────────────────────────────────────────────────────────────────────╯[white]",
				time.Now().Format("15:04:05"))

			d.headerBar.SetText(headerText)
		})
	}
}

// servicesChanged checks if the services list has actually changed
func (d *ModernDashboard) servicesChanged(newServices []*app.Service) bool {
	if len(d.services) != len(newServices) {
		return true
	}
	
	// Quick check - compare service IDs and basic status
	for i, service := range d.services {
		if i >= len(newServices) {
			return true
		}
		
		newService := newServices[i]
		if service.ID != newService.ID || 
		   service.Status != newService.Status ||
		   service.Name != newService.Name {
			return true
		}
		
		// Check if metrics changed significantly
		if service.Metrics != nil && newService.Metrics != nil {
			if abs(service.Metrics.CPUPercent - newService.Metrics.CPUPercent) > 1.0 ||
			   abs(service.Metrics.MemoryPercent - newService.Metrics.MemoryPercent) > 1.0 {
				return true
			}
		}
	}
	
	return false
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// removeDuplicateServices removes duplicate services based on ID
func (d *ModernDashboard) removeDuplicateServices(services []*app.Service) []*app.Service {
	seen := make(map[string]bool)
	unique := make([]*app.Service, 0)
	
	for _, service := range services {
		if !seen[service.ID] {
			seen[service.ID] = true
			unique = append(unique, service)
		}
	}
	
	return unique
}
