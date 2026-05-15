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
	servicePanel   *tview.TextView
	detailsPanel   *tview.TextView
	metricsPanel   *tview.TextView
	statsPanel     *tview.TextView
	headerPanel    *tview.TextView
	footerPanel    *tview.TextView
	rootLayout     tview.Primitive
	
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
	
	// Service list panel - created once with regions
	d.servicePanel = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetWordWrap(false)
	d.servicePanel.SetBorder(true).
		SetTitle(" 🚀 Services ").
		SetBorderColor(tcell.NewRGBColor(137, 180, 250)).
		SetTitleColor(tcell.NewRGBColor(203, 166, 247)).
		SetBorderPadding(0, 0, 1, 1)
	
	// Enable mouse support for service selection (click only, no hover to avoid redraws)
	d.servicePanel.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftClick {
			_, y := event.Position()
			// Get the relative position within the panel (accounting for border)
			_, _, _, panelY := d.servicePanel.GetInnerRect()
			relativeY := y - panelY - 1 // Subtract panel position and border
			
			// Calculate which service is at this position (each service takes 4 lines: name, details, metrics, blank)
			serviceIndex := relativeY / 4
			
			// Select service on click
			if serviceIndex >= 0 && serviceIndex < len(d.services) {
				d.selectedIndex = serviceIndex
				d.hoverIndex = -1 // Clear hover
				d.updateServiceListContent()
				d.updateDetailsContent()
				d.updateMetricsContent()
				d.updateStatsContent()
			}
			return action, nil
		}
		
		return action, event
	})
	
	// Details panel - created once with regions
	d.detailsPanel = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetWordWrap(true)
	d.detailsPanel.SetBorder(true).
		SetTitle(" 📋 Service Details ").
		SetBorderColor(tcell.NewRGBColor(137, 180, 250)).
		SetTitleColor(tcell.NewRGBColor(249, 226, 175)).
		SetBorderPadding(0, 0, 1, 1)
	
	// Metrics panel - created once with regions
	d.metricsPanel = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(false).
		SetWordWrap(false)
	d.metricsPanel.SetBorder(true).
		SetTitle(" 📊 Live Metrics ").
		SetBorderColor(tcell.NewRGBColor(137, 180, 250)).
		SetTitleColor(tcell.NewRGBColor(166, 227, 161)).
		SetBorderPadding(0, 0, 1, 1)
	
	// Statistics panel - created once with regions
	d.statsPanel = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetWordWrap(false)
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
	
	d.rootLayout = rootLayout
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
	d.headerPanel.SetText(header)
}

// updateFooterContent updates ONLY the footer text, never the structure
func (d *ProDashboard) updateFooterContent() {
	footer := "[#89b4fa]↑↓/Click[white] Navigate [#585b70]│[white] [#a6e3a1]R[white] Refresh [#585b70]│[white] [#cba6f7]T[white] Toggle View [#585b70]│[white] [#f9e2af]?[white] Help [#585b70]│[white] [#f38ba8]Q[white] Quit [#585b70]│[white] [#585b70::i]Content-Only Updates[::-][white]"
	d.footerPanel.SetText(footer)
}

// updateServiceListContent updates ONLY service list content, never borders
func (d *ProDashboard) updateServiceListContent() {
	if len(d.services) == 0 {
		emptyText := "\n\n  [#585b70::i]No services detected[white]\n  [#585b70]Scanning for services...[white]"
		if emptyText != d.lastServiceListText {
			d.servicePanel.SetText(emptyText)
			d.lastServiceListText = emptyText
		}
		return
	}

	var content strings.Builder
	content.WriteString("\n")

	for i, service := range d.services {
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

		// First line (name/type) with background per state
		content.WriteString(fmt.Sprintf("[\"service_%d\"]", i))
		if i == d.selectedIndex {
			content.WriteString(fmt.Sprintf("%s[white:#313244:b]%s %s %-23s [#585b70:#313244:]%-10s[white::]\n",
				indicator, statusIcon, typeIcon, name, service.Type))
		} else if i == d.hoverIndex {
			content.WriteString(fmt.Sprintf("%s[#cdd6f4:#1e1e2e:]%s %s %-23s [#585b70:#1e1e2e:]%-10s[white::]\n",
				indicator, statusIcon, typeIcon, name, service.Type))
		} else {
			content.WriteString(fmt.Sprintf("%s%s %s [#cdd6f4]%-23s[white] [#585b70]%-10s[white]\n",
				indicator, statusIcon, typeIcon, name, service.Type))
		}

		// Background tag reused for subsequent lines
		bgTag := ""
		if i == d.selectedIndex {
			bgTag = ":#313244"
		} else if i == d.hoverIndex {
			bgTag = ":#1e1e2e"
		}

		// Ports or uptime
		if len(service.Ports) > 0 && len(service.Ports) <= 3 {
			portStr := strings.Join(service.Ports, ",")
			if len(portStr) > 30 {
				portStr = portStr[:27] + "..."
			}
			content.WriteString(fmt.Sprintf("      [#585b70%s:]↳[white%s:] [#fab387%s:]:%s[white::]\n",
				bgTag, bgTag, bgTag, portStr))
		} else {
			uptime := time.Since(service.CreatedAt)
			content.WriteString(fmt.Sprintf("      [#585b70%s:]↳ up:[white%s:] [#a6adc8%s:]%s[white::]\n",
				bgTag, bgTag, bgTag, components.FormatDuration(int64(uptime.Seconds()))))
		}

		// Compact metrics
		if service.Metrics != nil {
			cpuBar := d.compactProgressBar(service.Metrics.CPUPercent, 8)
			memBar := d.compactProgressBar(service.Metrics.MemoryPercent, 8)
			content.WriteString(fmt.Sprintf("      [#585b70%s:]cpu[white::]%s [#585b70%s:]mem[white::]%s\n",
				bgTag, cpuBar, bgTag, memBar))
		}

		content.WriteString("[\"\"]")
		if i < len(d.services)-1 {
			content.WriteString("\n")
		}
	}

	// Update content ONLY - border never touched (only when changed)
	slt := content.String()
	if slt != d.lastServiceListText {
		d.servicePanel.SetText(slt)
		d.lastServiceListText = slt
		d.scrollToSelection()
	}
}

// scrollToSelection ensures the selected service is visible
func (d *ProDashboard) scrollToSelection() {
	tag := fmt.Sprintf("service_%d", d.selectedIndex)
	d.servicePanel.Highlight(tag)
	d.servicePanel.ScrollToHighlight()
}

// updateDetailsContent updates ONLY details content, never borders (cached)
func (d *ProDashboard) updateDetailsContent() {
    if d.selectedIndex >= len(d.services) || len(d.services) == 0 {
        txt := "\n  [#585b70::i]No service selected[white]"
        if txt != d.lastDetailsText {
            d.detailsPanel.SetText(txt)
            d.lastDetailsText = txt
        }
        return
    }

    service := d.services[d.selectedIndex]
    var details strings.Builder
    details.WriteString("\n")

    // Service information box
    details.WriteString(" [#cba6f7]╭─ Service Information ────────────╮[white]\n")
    details.WriteString(fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#cdd6f4]%s[white]\n", "Name", service.Name))
    details.WriteString(fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#89b4fa]%s[white] %s\n", "Type", service.Type, d.getServiceTypeIcon(service.Type)))
    details.WriteString(fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] %s %s\n", "Status", d.getStatusIcon(service.Status), service.Status))
    details.WriteString(fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#585b70]%s[white]\n", "ID", truncate(service.ID, 20)))

    if service.Image != "" {
        details.WriteString(fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#a6adc8]%s[white]\n", "Image", truncate(service.Image, 20)))
    }
    if service.Version != "" {
        details.WriteString(fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#89b4fa]%s[white]\n", "Version", service.Version))
    }
    if len(service.Ports) > 0 {
        portStr := strings.Join(service.Ports, ", ")
        if len(portStr) > 25 {
            portStr = portStr[:22] + "..."
        }
        details.WriteString(fmt.Sprintf(" [#cba6f7]│[white] [#f9e2af::b]%-15s[::-][white] [#a6e3a1]%s[white]\n", "Ports", portStr))
    }
    details.WriteString(" [#cba6f7]╰──────────────────────────────────╯[white]\n\n")

    // Runtime information
    details.WriteString(" [#89b4fa]╭─ Runtime ────────────────────────╮[white]\n")
    uptime := time.Since(service.CreatedAt)
    details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af::b]%-15s[::-][white] [#a6e3a1]%s[white]\n", "Uptime", components.FormatDuration(int64(uptime.Seconds()))))
    details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af::b]%-15s[::-][white] [#585b70]%s[white]\n", "Created", service.CreatedAt.Format("Jan 2, 15:04")))
    if service.HealthCheck != "" {
        details.WriteString(fmt.Sprintf(" [#89b4fa]│[white] [#f9e2af::b]%-15s[::-][white] %s\n", "Health", service.HealthCheck))
    }
    details.WriteString(" [#89b4fa]╰──────────────────────────────────╯[white]\n")

    // Apply only when changed
    dt := details.String()
    if dt != d.lastDetailsText {
        d.detailsPanel.SetText(dt)
        d.lastDetailsText = dt
    }
}

// updateMetricsContent updates ONLY metrics content, never borders
func (d *ProDashboard) updateMetricsContent() {
	if d.selectedIndex >= len(d.services) || len(d.services) == 0 {
		txt := "\n  [#585b70::i]No service selected[white]"
		if txt != d.lastMetricsText {
			d.metricsPanel.SetText(txt)
			d.lastMetricsText = txt
		}
		return
	}

	service := d.services[d.selectedIndex]
	if service.Metrics == nil {
		txt := "\n  [#585b70::i]Collecting metrics...[white]"
		if txt != d.lastMetricsText {
			d.metricsPanel.SetText(txt)
			d.lastMetricsText = txt
		}
		return
	}

	metrics := service.Metrics
	var content strings.Builder
	content.WriteString("\n")

	// CPU
	cpuBar := d.fullProgressBar(metrics.CPUPercent, 20)
	content.WriteString(" [#f9e2af::b]CPU[white]\n")
	content.WriteString(fmt.Sprintf(" %s\n", cpuBar))
	content.WriteString(fmt.Sprintf(" [#585b70]%5.1f%%[white]\n\n", metrics.CPUPercent))

	// Memory
	if metrics.MemoryLimit > 0 {
		memBar := d.fullProgressBar(metrics.MemoryPercent, 20)
		content.WriteString(" [#f9e2af::b]Memory[white]\n")
		content.WriteString(fmt.Sprintf(" %s\n", memBar))
		content.WriteString(fmt.Sprintf(" [#585b70]%5.1f%% ─ %s / %s[white]\n\n",
			metrics.MemoryPercent,
			components.FormatBytes(metrics.MemoryUsage),
			components.FormatBytes(metrics.MemoryLimit)))
	} else {
		content.WriteString(" [#f9e2af::b]Memory[white]\n")
		content.WriteString(fmt.Sprintf(" [#a6adc8]%s[white]\n\n", components.FormatBytes(metrics.MemoryUsage)))
	}

	// Network
	if metrics.NetworkIn > 0 || metrics.NetworkOut > 0 {
		content.WriteString(" [#f9e2af::b]Network[white]\n")
		content.WriteString(fmt.Sprintf(" [#a6e3a1]↓[white] %s/s\n", components.FormatBytes(metrics.NetworkIn)))
		content.WriteString(fmt.Sprintf(" [#f38ba8]↑[white] %s/s\n\n", components.FormatBytes(metrics.NetworkOut)))
	}

	// Disk I/O
	if metrics.DiskIORead > 0 || metrics.DiskIOWrite > 0 {
		content.WriteString(" [#f9e2af::b]Disk I/O[white]\n")
		content.WriteString(fmt.Sprintf(" [#a6e3a1]R:[white] %s\n", components.FormatBytes(metrics.DiskIORead)))
		content.WriteString(fmt.Sprintf(" [#f38ba8]W:[white] %s\n", components.FormatBytes(metrics.DiskIOWrite)))
	}

	// Update content ONLY - border never touched (only when changed)
	mt := content.String()
	if mt != d.lastMetricsText {
		d.metricsPanel.SetText(mt)
		d.lastMetricsText = mt
	}
}

// updateStatsContent updates ONLY statistics content, never borders
func (d *ProDashboard) updateStatsContent() {
	if d.selectedIndex >= len(d.services) || len(d.services) == 0 {
		txt := "\n  [#585b70::i]No service selected[white]"
		if txt != d.lastStatsText {
			d.statsPanel.SetText(txt)
			d.lastStatsText = txt
		}
		return
	}

	service := d.services[d.selectedIndex]
	statsText := FormatStatistics(service)
	// Update content ONLY - border never touched (only when changed)
	if statsText != d.lastStatsText {
		d.statsPanel.SetText(statsText)
		d.lastStatsText = statsText
	}
}

// contentOnlyUpdateLoop updates ONLY content at a slow rate, NEVER borders
func (d *ProDashboard) contentOnlyUpdateLoop() {
	d.updateTicker = time.NewTicker(5 * time.Second)
	defer d.updateTicker.Stop()

	for range d.updateTicker.C {
		newServices := d.app.GetServices()
		if d.servicesChanged(newServices) {
			d.tviewApp.QueueUpdate(func() {
				d.services = newServices
				d.updateServiceListContent()
				d.updateDetailsContent()
				d.updateMetricsContent()
				d.updateStatsContent()
				d.updateHeaderContent()
			})
		}
	}
}

// setupKeyboardHandling sets up keyboard shortcuts
func (d *ProDashboard) setupKeyboardHandling() {
	d.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp, tcell.KeyCtrlP:
			d.navigateUp()
			return nil
		case tcell.KeyDown, tcell.KeyCtrlN:
			d.navigateDown()
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
			d.showHelpModal()
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
	d.updateDetailsContent()
	d.updateMetricsContent()
	d.updateStatsContent()
	
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

// showHelpModal shows a help modal with shortcuts
func (d *ProDashboard) showHelpModal() {
	helpText := `[#cba6f7::b]LazyService Pro - Keyboard Shortcuts[white]

[#89b4fa]Navigation:[white]
  ↑/↓, j/k    Navigate service list
  Click       Select service
  Enter       Refresh selection

[#a6e3a1]Actions:[white]
  R, r        Manual refresh
  T, t        Toggle dashboard view (Coming soon)

[#f9e2af]General:[white]
  ?, H, h     Show this help
  Q, q        Quit application

[#585b70::i]Press any key to close help...[::-]`

	modal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			d.tviewApp.SetRoot(d.rootLayout, true)
		})

	modal.SetBorder(true).
		SetTitle(" Help ").
		SetTitleColor(tcell.NewRGBColor(249, 226, 175)).
		SetBorderColor(tcell.NewRGBColor(137, 180, 250))

	d.tviewApp.SetRoot(modal, true)
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
