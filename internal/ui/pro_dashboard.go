package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/lazyservice/lazyservice/internal/app"
	"github.com/lazyservice/lazyservice/internal/ui/components"
	"github.com/rivo/tview"
)

// ProDashboard - Ultra-professional dashboard with content-only updates
// Borders are created ONCE and NEVER modified. Only content regions update.
type ProDashboard struct {
	app            *app.App
	tviewApp       *tview.Application
	
	// Main panels - created once, never recreated
	servicePanel   *tview.Table
	detailsPanel   *tview.Table
	metricsPanel   *tview.Table
	statsPanel     *tview.Table
	headerPanel    *tview.TextView
	footerPanel    *tview.TextView
	
	// State
	services       []*app.Service
	selectedIndex  int
	hoverIndex     int
	theme          *Theme
	lastUpdate     time.Time
	
	// Last rendered caches (for content-only updates)
	lastServiceListText string
	lastDetailsText     string
	lastMetricsText     string
	lastStatsText       string
	
	// Update control
	bordersCreated bool
	updateTicker   *time.Ticker
}

// NewProDashboard creates the ultimate professional dashboard
func NewProDashboard(application *app.App) *ProDashboard {
	theme := ModernDarkTheme()
	ApplyTheme(theme)
	
	return &ProDashboard{
		app:            application,
		tviewApp:       tview.NewApplication(),
		services:       make([]*app.Service, 0),
		selectedIndex:  0,
		hoverIndex:     -1,
		theme:          theme,
		bordersCreated: false,
	}
}

// Run starts the professional dashboard
func (d *ProDashboard) Run() error {
	// Step 1: Create borders ONCE - they will NEVER be touched again
	d.createBordersOnce()
	
	// Step 2: Load initial content into regions
	d.loadInitialContent()
	
	// Step 3: Start content-only background updates
	go d.contentOnlyUpdateLoop()
	
	// Step 4: Setup keyboard handling
	d.setupKeyboardHandling()
	
	// Step 5: Enable mouse support
	d.tviewApp.EnableMouse(true)
	
	// Run the application
	return d.tviewApp.Run()
}

// Stop stops the dashboard
func (d *ProDashboard) Stop() {
	if d.updateTicker != nil {
		d.updateTicker.Stop()
	}
	d.tviewApp.Stop()
}

// createBordersOnce creates all borders ONCE - they will NEVER be modified
func (d *ProDashboard) createBordersOnce() {
	// Header panel - created once
	d.headerPanel = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetRegions(true)

	// Service list panel - created once as a Table
	d.servicePanel = tview.NewTable().
		SetSelectable(true, false)
	d.servicePanel.SetBorder(true).
		SetTitle(" 🚀 Services ").
		SetBorderColor(tcell.NewRGBColor(137, 180, 250)).
		SetTitleColor(tcell.NewRGBColor(203, 166, 247)).
		SetBorderPadding(0, 0, 1, 1)

	d.servicePanel.SetSelectionChangedFunc(func(row, column int) {
		serviceIndex := row / 4
		if serviceIndex >= 0 && serviceIndex < len(d.services) && serviceIndex != d.selectedIndex {
			d.selectedIndex = serviceIndex
			d.updateServiceListContent()
			d.updateDetailsContent()
			d.updateMetricsContent()
			d.updateStatsContent()
		}
	})

	// Details panel - created once as a Table
	d.detailsPanel = tview.NewTable()
	d.detailsPanel.SetBorder(true).
		SetTitle(" 📋 Service Details ").
		SetBorderColor(tcell.NewRGBColor(137, 180, 250)).
		SetTitleColor(tcell.NewRGBColor(249, 226, 175)).
		SetBorderPadding(0, 0, 1, 1)

	// Metrics panel - created once as a Table
	d.metricsPanel = tview.NewTable()
	d.metricsPanel.SetBorder(true).
		SetTitle(" 📊 Live Metrics ").
		SetBorderColor(tcell.NewRGBColor(137, 180, 250)).
		SetTitleColor(tcell.NewRGBColor(166, 227, 161)).
		SetBorderPadding(0, 0, 1, 1)

	// Statistics panel - created once as a Table
	d.statsPanel = tview.NewTable()
	d.statsPanel.SetBorder(true).
		SetTitle(" 📈 Statistics & Insights ").
		SetBorderColor(tcell.NewRGBColor(137, 180, 250)).
		SetTitleColor(tcell.NewRGBColor(245, 194, 231)).
		SetBorderPadding(0, 0, 1, 1)
	
	// Footer panel - created once with regions
	d.footerPanel = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter)
	
	// Create layout - assembled once, never modified
	rightTop := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(d.detailsPanel, 0, 1, false).
		AddItem(d.metricsPanel, 0, 1, false)
	
	rightBottom := d.statsPanel
	
	rightPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(rightTop, 0, 1, false).
		AddItem(rightBottom, 12, 0, false)
	
	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(d.servicePanel, 45, 0, true).
		AddItem(rightPanel, 0, 1, false)
	
	rootLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.headerPanel, 3, 0, false).
		AddItem(mainLayout, 0, 1, true).
		AddItem(d.footerPanel, 1, 0, false)
	
	d.tviewApp.SetRoot(rootLayout, true)
	
	// Mark borders as created - they will NEVER be touched again
	d.bordersCreated = true
}

// loadInitialContent loads initial content into regions (borders already exist)
func (d *ProDashboard) loadInitialContent() {
	// Update header content only
	d.updateHeaderContent()
	
	// Update footer content only
	d.updateFooterContent()
	
	// Load services
	d.services = d.app.GetServices()
	
	// Update all content regions
	d.updateServiceListContent()
	d.updateDetailsContent()
	d.updateMetricsContent()
	d.updateStatsContent()
}

// updateHeaderContent updates ONLY the header text, never the structure
func (d *ProDashboard) updateHeaderContent() {
	header := fmt.Sprintf(
		"\n[#cba6f7::b]⚡ LazyService Pro[::-] [#585b70]│[white] [#89b4fa]Professional Dashboard[white] [#585b70]│[white] [#fab387]%d Services[white] [#585b70]│[white] [#a6e3a1]%s[white]\n",
		len(d.services),
		time.Now().Format("15:04:05 Mon Jan 2"),
	)
	if d.headerPanel.GetText(false) != header {
		d.headerPanel.SetText(header)
	}
}

// updateFooterContent updates ONLY the footer text, never the structure
func (d *ProDashboard) updateFooterContent() {
	footer := "[#89b4fa]↑↓/Click[white] Navigate [#585b70]│[white] [#a6e3a1]R[white] Refresh [#585b70]│[white] [#cba6f7]T[white] Toggle View [#585b70]│[white] [#f38ba8]Q[white] Quit [#585b70]│[white] [#585b70::i]Content-Only Updates[::-][white]"
	d.footerPanel.SetText(footer)
}

func (d *ProDashboard) updateTableCell(table *tview.Table, row, col int, text string) {
	cell := table.GetCell(row, col)
	if cell.Text != text {
		cell.SetText(text)
	}
}

// updateServiceListContent updates ONLY service list content, never borders
func (d *ProDashboard) updateServiceListContent() {
	if len(d.services) == 0 {
		d.updateTableCell(d.servicePanel, 0, 0, "[#585b70::i]No services detected[white]")
		return
	}

	for i, service := range d.services {
		row := i * 4

		// Selection indicator
		indicator := "   "
		if i == d.selectedIndex {
			indicator = " [#89b4fa]❯[white] "
		}

		// Icons and name
		statusIcon := d.getStatusIcon(service.Status)
		typeIcon := d.getServiceTypeIcon(service.Type)
		name := service.Name
		if len(name) > 23 {
			name = name[:20] + "..."
		}

		// First line (name/type)
		style := "[white]"
		if i == d.selectedIndex {
			style = "[white:#313244:b]"
		}
		d.updateTableCell(d.servicePanel, row, 0, fmt.Sprintf("%s%s%s %s %-23s [#585b70]%s[white::]",
			indicator, style, statusIcon, typeIcon, name, service.Type))

		// Background tag for subsequent lines
		bgTag := ""
		if i == d.selectedIndex {
			bgTag = ":#313244"
		}

		// Second line: Ports or uptime
		line2 := ""
		if len(service.Ports) > 0 && len(service.Ports) <= 3 {
			portStr := strings.Join(service.Ports, ",")
			if len(portStr) > 30 {
				portStr = portStr[:27] + "..."
			}
			line2 = fmt.Sprintf("      [#585b70%s:]↳[white%s:] [#fab387%s:]:%s[white::]", bgTag, bgTag, bgTag, portStr)
		} else {
			uptime := time.Since(service.CreatedAt)
			line2 = fmt.Sprintf("      [#585b70%s:]↳ up:[white%s:] [#a6adc8%s:]%s[white::]", bgTag, bgTag, bgTag, components.FormatDuration(int64(uptime.Seconds())))
		}
		d.updateTableCell(d.servicePanel, row+1, 0, line2)

		// Third line: Compact metrics
		line3 := ""
		if service.Metrics != nil {
			cpuBar := d.compactProgressBar(service.Metrics.CPUPercent, 8)
			memBar := d.compactProgressBar(service.Metrics.MemoryPercent, 8)
			line3 = fmt.Sprintf("      [#585b70%s:]cpu[white::]%s [#585b70%s:]mem[white::]%s", bgTag, cpuBar, bgTag, memBar)
		}
		d.updateTableCell(d.servicePanel, row+2, 0, line3)

		// Fourth line: Spacer
		d.updateTableCell(d.servicePanel, row+3, 0, "")
	}
}

// updateDetailsContent updates ONLY details content, never borders (cached)
func (d *ProDashboard) updateDetailsContent() {
	if d.selectedIndex >= len(d.services) || len(d.services) == 0 {
		d.updateTableCell(d.detailsPanel, 0, 0, "[#585b70::i]No service selected[white]")
		return
	}

	service := d.services[d.selectedIndex]

	// Service information rows
	d.updateTableCell(d.detailsPanel, 0, 0, "[#cba6f7]╭─ Service Information ────────────╮[white]")
	d.updateTableCell(d.detailsPanel, 1, 0, fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#cdd6f4]%s[white]", "Name", service.Name))
	d.updateTableCell(d.detailsPanel, 2, 0, fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#89b4fa]%s[white] %s", "Type", service.Type, d.getServiceTypeIcon(service.Type)))
	d.updateTableCell(d.detailsPanel, 3, 0, fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] %s %s", "Status", d.getStatusIcon(service.Status), service.Status))
	d.updateTableCell(d.detailsPanel, 4, 0, fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#585b70]%s[white]", "ID", truncate(service.ID, 20)))

	row := 5
	if service.Image != "" {
		d.updateTableCell(d.detailsPanel, row, 0, fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#a6adc8]%s[white]", "Image", truncate(service.Image, 20)))
		row++
	}
	if service.Version != "" {
		d.updateTableCell(d.detailsPanel, row, 0, fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#89b4fa]%s[white]", "Version", service.Version))
		row++
	}
	if len(service.Ports) > 0 {
		portStr := strings.Join(service.Ports, ", ")
		if len(portStr) > 25 {
			portStr = portStr[:22] + "..."
		}
		d.updateTableCell(d.detailsPanel, row, 0, fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#a6e3a1]%s[white]", "Ports", portStr))
		row++
	}
	d.updateTableCell(d.detailsPanel, row, 0, "[#cba6f7]╰──────────────────────────────────╯[white]")
	row++

	// Runtime information
	d.updateTableCell(d.detailsPanel, row, 0, "")
	row++
	d.updateTableCell(d.detailsPanel, row, 0, " [#89b4fa]╭─ Runtime ────────────────────────╮[white]")
	row++
	uptime := time.Since(service.CreatedAt)
	d.updateTableCell(d.detailsPanel, row, 0, fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af::b]%-15s[::-][white] [#a6e3a1]%s[white]", "Uptime", components.FormatDuration(int64(uptime.Seconds()))))
	row++
	d.updateTableCell(d.detailsPanel, row, 0, fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af::b]%-15s[::-][white] [#585b70]%s[white]", "Created", service.CreatedAt.Format("Jan 2, 15:04")))
	row++
	if service.HealthCheck != "" {
		d.updateTableCell(d.detailsPanel, row, 0, fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af::b]%-15s[::-][white] %s", "Health", service.HealthCheck))
		row++
	}
	d.updateTableCell(d.detailsPanel, row, 0, " [#89b4fa]╰──────────────────────────────────╯[white]")
}

// updateMetricsContent updates ONLY metrics content, never borders
func (d *ProDashboard) updateMetricsContent() {
	if d.selectedIndex >= len(d.services) || len(d.services) == 0 {
		d.updateTableCell(d.metricsPanel, 0, 0, "[#585b70::i]No service selected[white]")
		return
	}

	service := d.services[d.selectedIndex]
	if service.Metrics == nil {
		d.updateTableCell(d.metricsPanel, 0, 0, "[#585b70::i]Collecting metrics...[white]")
		return
	}

	metrics := service.Metrics
	d.updateTableCell(d.metricsPanel, 0, 0, "")
	d.updateTableCell(d.metricsPanel, 1, 0, " [#f9e2af::b]CPU[white]")
	d.updateTableCell(d.metricsPanel, 2, 0, fmt.Sprintf(" %s", d.fullProgressBar(metrics.CPUPercent, 20)))
	d.updateTableCell(d.metricsPanel, 3, 0, fmt.Sprintf(" [#585b70]%5.1f%%[white]", metrics.CPUPercent))

	d.updateTableCell(d.metricsPanel, 4, 0, "")
	d.updateTableCell(d.metricsPanel, 5, 0, " [#f9e2af::b]Memory[white]")
	if metrics.MemoryLimit > 0 {
		d.updateTableCell(d.metricsPanel, 6, 0, fmt.Sprintf(" %s", d.fullProgressBar(metrics.MemoryPercent, 20)))
		d.updateTableCell(d.metricsPanel, 7, 0, fmt.Sprintf(" [#585b70]%5.1f%% ─ %s / %s[white]",
			metrics.MemoryPercent,
			components.FormatBytes(metrics.MemoryUsage),
			components.FormatBytes(metrics.MemoryLimit)))
	} else {
		d.updateTableCell(d.metricsPanel, 6, 0, fmt.Sprintf(" [#a6adc8]%s[white]", components.FormatBytes(metrics.MemoryUsage)))
		d.updateTableCell(d.metricsPanel, 7, 0, "")
	}

	d.updateTableCell(d.metricsPanel, 8, 0, "")
	d.updateTableCell(d.metricsPanel, 9, 0, " [#f9e2af::b]Network[white]")
	d.updateTableCell(d.metricsPanel, 10, 0, fmt.Sprintf(" [#a6e3a1]↓[white] %s/s", components.FormatBytes(metrics.NetworkIn)))
	d.updateTableCell(d.metricsPanel, 11, 0, fmt.Sprintf(" [#f38ba8]↑[white] %s/s", components.FormatBytes(metrics.NetworkOut)))

	d.updateTableCell(d.metricsPanel, 12, 0, "")
	d.updateTableCell(d.metricsPanel, 13, 0, " [#f9e2af::b]Disk I/O[white]")
	d.updateTableCell(d.metricsPanel, 14, 0, fmt.Sprintf(" [#a6e3a1]R:[white] %s", components.FormatBytes(metrics.DiskIORead)))
	d.updateTableCell(d.metricsPanel, 15, 0, fmt.Sprintf(" [#f38ba8]W:[white] %s", components.FormatBytes(metrics.DiskIOWrite)))
}

// updateStatsContent updates ONLY statistics content, never borders
func (d *ProDashboard) updateStatsContent() {
	if d.selectedIndex >= len(d.services) || len(d.services) == 0 {
		d.updateTableCell(d.statsPanel, 0, 0, "[#585b70::i]No service selected[white]")
		return
	}

	service := d.services[d.selectedIndex]
	statsText := FormatStatistics(service)
	lines := strings.Split(statsText, "\n")
	for i, line := range lines {
		d.updateTableCell(d.statsPanel, i, 0, line)
	}
}

// contentOnlyUpdateLoop updates ONLY content at a slow rate, NEVER borders
func (d *ProDashboard) contentOnlyUpdateLoop() {
	// Update metrics every 2 seconds for a "live" feel
	d.updateTicker = time.NewTicker(2 * time.Second)
	defer d.updateTicker.Stop()

	for range d.updateTicker.C {
		d.tviewApp.QueueUpdate(func() {
			d.services = d.app.GetServices()
			d.updateServiceListContent()
			d.updateDetailsContent()
			d.updateMetricsContent()
			d.updateStatsContent()
			d.updateHeaderContent()
		})
	}
}

// setupKeyboardHandling sets up keyboard shortcuts
func (d *ProDashboard) setupKeyboardHandling() {
	d.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp, tcell.KeyCtrlP:
			d.navigateUp()
			d.servicePanel.Select(d.selectedIndex*4, 0)
			return nil
		case tcell.KeyDown, tcell.KeyCtrlN:
			d.navigateDown()
			d.servicePanel.Select(d.selectedIndex*4, 0)
			return nil
		case tcell.KeyEnter:
			// Service already selected, just refresh details
			d.refreshContentOnly()
			return nil
		}
		
		switch event.Rune() {
		case 'q', 'Q':
			d.Stop()
			return nil
		case 'r', 'R':
			d.manualRefresh()
			return nil
		case 'k':
			d.navigateUp()
			return nil
		case 'j':
			d.navigateDown()
			return nil
		case 't', 'T':
			d.toggleView()
			return nil
		case 'h', 'H', '?':
			d.showHelp()
			return nil
		}
		
		return event
	})
}

// navigateUp moves selection up (content-only update)
func (d *ProDashboard) navigateUp() {
    if len(d.services) == 0 {
        return
    }
    d.selectedIndex--
    if d.selectedIndex < 0 {
        d.selectedIndex = len(d.services) - 1
    }
    d.hoverIndex = -1
    d.updateServiceListContent()
    d.updateDetailsContent()
    d.updateMetricsContent()
    d.updateStatsContent()
}

// navigateDown moves selection down (content-only update)
func (d *ProDashboard) navigateDown() {
    if len(d.services) == 0 {
        return
    }
    d.selectedIndex++
    if d.selectedIndex >= len(d.services) {
        d.selectedIndex = 0
    }
    d.hoverIndex = -1
    d.updateServiceListContent()
    d.updateDetailsContent()
    d.updateMetricsContent()
    d.updateStatsContent()
}

// refreshContentOnly refreshes ONLY content, NEVER borders
func (d *ProDashboard) refreshContentOnly() {
    d.updateServiceListContent()
    d.updateDetailsContent()
    d.updateMetricsContent()
    d.updateStatsContent()
}

// manualRefresh forces a refresh of data and content
func (d *ProDashboard) manualRefresh() {
	d.app.RefreshServices()
	d.services = d.app.GetServices()
	d.updateServiceListContent()
	d.updateDetailsContent()
	d.updateMetricsContent()
	d.updateStatsContent()
	d.updateHeaderContent()

	// Show temporary status message
	originalFooter := d.footerPanel.GetText(false)
	d.footerPanel.SetText("[#a6e3a1]✨ Refreshed - Values Only[white]")
	
	go func() {
		time.Sleep(2 * time.Second)
		d.tviewApp.QueueUpdate(func() {
			d.footerPanel.SetText(originalFooter)
		})
	}()
}

// toggleView toggles between different views
func (d *ProDashboard) toggleView() {
	// Future: Toggle between different stat views
	d.refreshContentOnly()
}

// showHelp shows help information
func (d *ProDashboard) showHelp() {
	originalFooter := d.footerPanel.GetText(false)
	d.footerPanel.SetText("[#89b4fa]Help: ↑↓/jk=Navigate | R=Refresh | T=Toggle | Q=Quit | Press any key to dismiss[white]")
	
	go func() {
		time.Sleep(5 * time.Second)
		d.tviewApp.QueueUpdate(func() {
			d.footerPanel.SetText(originalFooter)
		})
	}()
}

// Helper functions

func (d *ProDashboard) getStatusIcon(status app.ServiceStatus) string {
	switch status {
	case app.ServiceStatusRunning:
		return "[#a6e3a1]●[white]"
	case app.ServiceStatusStopped:
		return "[#f38ba8]●[white]"
	case app.ServiceStatusFailed:
		return "[#f38ba8]✗[white]"
	default:
		return "[#585b70]○[white]"
	}
}

func (d *ProDashboard) getServiceTypeIcon(serviceType app.ServiceType) string {
	switch serviceType {
	case app.ServiceTypeDocker:
		return "[#89b4fa]🐳[white]"
	case app.ServiceTypeKubernetes:
		return "[#89b4fa]☸️[white]"
	case app.ServiceTypeSystemd:
		return "[#fab387]⚙️[white]"
	case app.ServiceTypeProcess:
		return "[#a6e3a1]⚡[white]"
	default:
		return "[#585b70]📦[white]"
	}
}

func (d *ProDashboard) miniProgressBar(percent float64, width int) string {
	filled := int((percent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	
	bar := "[#a6e3a1]"
	bar += strings.Repeat("▮", filled)
	bar += "[#585b70]"
	bar += strings.Repeat("▯", width-filled)
	bar += fmt.Sprintf("[#a6adc8] %4.1f%%[white]", percent)
	
	return bar
}

func (d *ProDashboard) compactProgressBar(percent float64, width int) string {
	filled := int((percent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	
	// Color based on usage
	color := "#a6e3a1" // green
	if percent > 60 {
		color = "#f9e2af" // yellow
	}
	if percent > 85 {
		color = "#f38ba8" // red
	}
	
	bar := "[" + color + "]"
	bar += strings.Repeat("█", filled)
	bar += "[#313244]"
	bar += strings.Repeat("▒", width-filled)
	bar += fmt.Sprintf("[#585b70] %2.0f%%[white]", percent)
	
	return bar
}

func (d *ProDashboard) fullProgressBar(percent float64, width int) string {
	filled := int((percent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	
	// Choose color based on percentage
	color := "#a6e3a1" // green
	if percent > 70 {
		color = "#f9e2af" // yellow
	}
	if percent > 90 {
		color = "#f38ba8" // red
	}
	
	bar := "[" + color + "]"
	bar += strings.Repeat("█", filled)
	bar += "[#313244]"
	bar += strings.Repeat("░", width-filled)
	bar += "[white]"
	
	return bar
}

func (d *ProDashboard) servicesChanged(newServices []*app.Service) bool {
	if len(d.services) != len(newServices) {
		return true
	}
	
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
		
		// Check metrics changes
		if service.Metrics != nil && newService.Metrics != nil {
			if absFloat(service.Metrics.CPUPercent - newService.Metrics.CPUPercent) > 2.0 ||
			   absFloat(service.Metrics.MemoryPercent - newService.Metrics.MemoryPercent) > 2.0 {
				return true
			}
		}
	}
	
	return false
}

// absFloat returns absolute value of a float64
func absFloat(f float64) float64 {
    if f < 0 {
        return -f
    }
    return f
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
