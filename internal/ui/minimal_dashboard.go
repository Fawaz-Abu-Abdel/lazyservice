package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/lazyservice/lazyservice/internal/app"
	"github.com/lazyservice/lazyservice/internal/ui/components"
	"golang.org/x/term"
)

// MinimalDashboard - True values-only updates using ANSI cursor positioning
// Borders are drawn ONCE and NEVER touched again
type MinimalDashboard struct {
	app      *app.App
	services []*app.Service
	selected int
	running  bool

	// Terminal dimensions
	termWidth  int
	termHeight int

	// Totals tracking
	totalRequests   uint64
	totalNetworkIn  uint64
	totalNetworkOut uint64

	// Input handling
	screen tcell.Screen
	out    *bufio.Writer
}

// NewMinimalDashboard creates a minimal dashboard
func NewMinimalDashboard(application *app.App) *MinimalDashboard {
	return &MinimalDashboard{
		app:      application,
		services: make([]*app.Service, 0),
		selected: 0,
		running:  true,
		out:      bufio.NewWriter(os.Stdout),
	}
}

// Run starts the minimal dashboard
func (d *MinimalDashboard) Run() error {
	// Initialize tcell screen for input handling only
	var err error
	d.screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := d.screen.Init(); err != nil {
		return err
	}
	d.screen.EnableMouse()
	defer d.screen.Fini()

	d.getTerminalSize()
	d.hideCursor()
	defer d.showCursor()

	// Draw borders ONCE
	d.drawBorders()
	d.out.Flush()

	// Initial data
	d.services = d.app.GetServices()
	d.calculateTotals()
	d.updateAllValues()
	d.out.Flush()

	// Background update loop
	go d.updateLoop()

	// Handle input
	for d.running {
		ev := d.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			d.getTerminalSize()
			d.drawBorders()
			d.updateAllValues()
			d.out.Flush()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' || ev.Rune() == 'Q' {
				d.running = false
			} else if ev.Key() == tcell.KeyUp || ev.Rune() == 'k' {
				d.navigateUp()
			} else if ev.Key() == tcell.KeyDown || ev.Rune() == 'j' {
				d.navigateDown()
			} else if ev.Rune() == 'r' || ev.Rune() == 'R' {
				d.app.RefreshServices()
				d.services = d.app.GetServices()
				d.calculateTotals()
				d.updateAllValues()
				d.out.Flush()
			}
		case *tcell.EventMouse:
			if ev.Buttons()&tcell.Button1 != 0 {
				x, y := ev.Position()
				d.handleMouseClick(x, y)
			}
		}
	}

	d.clearScreen()
	d.moveCursor(1, 1)
	d.out.Flush()
	return nil
}

// Stop stops the dashboard
func (d *MinimalDashboard) Stop() {
	d.running = false
}

// calculateTotals calculates aggregate metrics across all services
func (d *MinimalDashboard) calculateTotals() {
	d.totalRequests = 0
	d.totalNetworkIn = 0
	d.totalNetworkOut = 0

	for _, service := range d.services {
		if service.Metrics != nil {
			d.totalRequests += uint64(service.Metrics.RequestsPerSec * 60) // Approximate
			d.totalNetworkIn += service.Metrics.NetworkIn
			d.totalNetworkOut += service.Metrics.NetworkOut
		}
	}
}

// getTerminalSize gets terminal dimensions
func (d *MinimalDashboard) getTerminalSize() {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		d.termWidth = width
		d.termHeight = height
	} else {
		d.termWidth = 120
		d.termHeight = 30
	}
}

func (d *MinimalDashboard) hideCursor()  { fmt.Fprint(d.out, "\033[?25l") }
func (d *MinimalDashboard) showCursor()  { fmt.Fprint(d.out, "\033[?25h") }
func (d *MinimalDashboard) clearScreen() { fmt.Fprint(d.out, "\033[2J\033[H") }
func (d *MinimalDashboard) moveCursor(row, col int) {
	fmt.Fprintf(d.out, "\033[%d;%dH", row, col)
}

// drawBorders draws all borders and static labels ONCE - never called again
func (d *MinimalDashboard) drawBorders() {
	d.clearScreen()

	// Colors
	cyan := "\033[36m"
	yellow := "\033[33m"
	magenta := "\033[35m"
	green := "\033[32m"
	blue := "\033[34m"
	reset := "\033[0m"
	dim := "\033[90m"

	// Header (row 1-3)
	d.moveCursor(1, 1)
	fmt.Fprint(d.out, cyan+"╔"+strings.Repeat("═", d.termWidth-2)+"╗"+reset)
	d.moveCursor(2, 1)
	fmt.Fprint(d.out, cyan+"║"+reset+"  "+cyan+"⚡ LazyService Pro"+reset+"  │  "+reset)
	d.moveCursor(2, 30)
	fmt.Fprint(d.out, reset+" services  │  "+reset)
	d.moveCursor(2, 50)
	fmt.Fprint(d.out, reset+" running  │  "+reset)
	d.moveCursor(2, d.termWidth)
	fmt.Fprint(d.out, cyan+"║"+reset)

	d.moveCursor(3, 1)
	fmt.Fprint(d.out, cyan+"╠"+strings.Repeat("═", d.termWidth-2)+"╣"+reset)

	// Totals bar (row 4)
	d.moveCursor(4, 1)
	fmt.Fprint(d.out, cyan+"║"+reset+"  "+cyan+"📊 Totals:"+reset+"  "+yellow+"Requests: "+reset)
	d.moveCursor(4, 40)
	fmt.Fprint(d.out, "│  "+green+"Net↓: "+reset)
	d.moveCursor(4, 60)
	fmt.Fprint(d.out, magenta+"Net↑: "+reset)
	d.moveCursor(4, 80)
	fmt.Fprint(d.out, "│  "+cyan+"Total Net: "+reset)
	d.moveCursor(4, d.termWidth)
	fmt.Fprint(d.out, cyan+"║"+reset)

	d.moveCursor(5, 1)
	fmt.Fprint(d.out, cyan+"╚"+strings.Repeat("═", d.termWidth-2)+"╝"+reset)

	// Services panel (left, rows 6-22)
	serviceWidth := 42
	d.moveCursor(6, 1)
	fmt.Fprint(d.out, yellow+"┌─ Services ─"+strings.Repeat("─", serviceWidth-13)+"┐"+reset)
	for i := 7; i <= 21; i++ {
		d.moveCursor(i, 1)
		fmt.Fprint(d.out, yellow+"│"+strings.Repeat(" ", serviceWidth-1)+"│"+reset)
	}
	d.moveCursor(22, 1)
	fmt.Fprint(d.out, yellow+"└"+strings.Repeat("─", serviceWidth-1)+"┘"+reset)

	// Details panel (right top, rows 6-13)
	detailStart := serviceWidth + 2
	detailWidth := d.termWidth - serviceWidth - 3
	d.moveCursor(6, detailStart)
	fmt.Fprint(d.out, magenta+"┌─ Details "+strings.Repeat("─", detailWidth-12)+"┐"+reset)
	for i := 7; i <= 12; i++ {
		d.moveCursor(i, detailStart)
		fmt.Fprint(d.out, magenta+"│"+strings.Repeat(" ", detailWidth-1)+"│"+reset)
	}

	// Static labels for details
	d.moveCursor(7, detailStart+2)
	fmt.Fprint(d.out, yellow+"Name: "+reset)
	d.moveCursor(8, detailStart+2)
	fmt.Fprint(d.out, yellow+"Type: "+reset)
	d.moveCursor(8, detailStart+30)
	fmt.Fprint(d.out, "│  "+yellow+"Status: "+reset)
	d.moveCursor(9, detailStart+2)
	fmt.Fprint(d.out, yellow+"ID: "+reset)
	d.moveCursor(10, detailStart+2)
	fmt.Fprint(d.out, yellow+"Uptime: "+reset)
	d.moveCursor(10, detailStart+30)
	fmt.Fprint(d.out, "│  "+yellow+"Ports: "+reset)
	d.moveCursor(11, detailStart+2)
	fmt.Fprint(d.out, yellow+"Threads: "+reset)
	d.moveCursor(11, detailStart+30)
	fmt.Fprint(d.out, "│  "+yellow+"Connections: "+reset)
	d.moveCursor(12, detailStart+2)
	fmt.Fprint(d.out, yellow+"Requests/sec: "+reset)

	d.moveCursor(13, detailStart)
	fmt.Fprint(d.out, magenta+"└"+strings.Repeat("─", detailWidth-1)+"┘"+reset)

	// Metrics panel (right middle, rows 14-22)
	d.moveCursor(14, detailStart)
	fmt.Fprint(d.out, green+"┌─ Live Metrics "+strings.Repeat("─", detailWidth-17)+"┐"+reset)
	for i := 15; i <= 21; i++ {
		d.moveCursor(i, detailStart)
		fmt.Fprint(d.out, green+"│"+strings.Repeat(" ", detailWidth-1)+"│"+reset)
	}

	// Static labels for metrics
	d.moveCursor(15, detailStart+2)
	fmt.Fprint(d.out, yellow+"CPU:    "+reset)
	d.moveCursor(16, detailStart+2)
	fmt.Fprint(d.out, yellow+"Memory: "+reset)
	d.moveCursor(17, detailStart+2)
	fmt.Fprint(d.out, yellow+"Disk:   "+reset)
	d.moveCursor(18, detailStart+2)
	fmt.Fprint(d.out, yellow+"Net:    "+reset)
	d.moveCursor(19, detailStart+2)
	fmt.Fprint(d.out, yellow+"Uptime: "+reset)
	d.moveCursor(20, detailStart+2)
	fmt.Fprint(d.out, yellow+"Threads: "+reset)
	d.moveCursor(20, detailStart+30)
	fmt.Fprint(d.out, "│  "+yellow+"Connections: "+reset)
	d.moveCursor(21, detailStart+2)
	fmt.Fprint(d.out, yellow+"Total Network: "+reset)

	d.moveCursor(22, detailStart)
	fmt.Fprint(d.out, green+"└"+strings.Repeat("─", detailWidth-1)+"┘"+reset)

	// Footer (row 23)
	d.moveCursor(23, 1)
	fmt.Fprint(d.out, blue+strings.Repeat("─", d.termWidth)+reset)
	d.moveCursor(24, 1)
	fmt.Fprint(d.out, dim+strings.Repeat(" ", d.termWidth)+reset)
	d.moveCursor(24, 3)
	fmt.Print(cyan + "↑↓" + dim + "/" + cyan + "jk" + reset + " Navigate  " +
		yellow + "🖱️ Click" + reset + " Select  " +
		green + "R" + reset + " Refresh  " +
		"\033[31mQ" + reset + " Quit  " +
		dim + "│ Updated: " + reset)
}

// updateLoop updates values every 2 seconds
func (d *MinimalDashboard) updateLoop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for d.running {
		<-ticker.C
		d.services = d.app.GetServices()
		d.calculateTotals()
		d.updateAllValues()
		d.out.Flush()
	}
}

// updateAllValues updates all value regions
func (d *MinimalDashboard) updateAllValues() {
	d.updateHeader()
	d.updateTotalsBar()
	d.updateServiceList()
	d.updateDetails()
	d.updateMetrics()
	d.updateFooter()
}

// updateHeader updates header values only
func (d *MinimalDashboard) updateHeader() {
	white := "\033[97m"
	yellow := "\033[33m"
	green := "\033[32m"
	reset := "\033[0m"

	// Count running services
	running := 0
	for _, s := range d.services {
		if s.Status == app.ServiceStatusRunning {
			running++
		}
	}

	d.moveCursor(2, 26)
	fmt.Fprint(d.out, yellow+fmt.Sprintf("%2d", len(d.services))+reset)

	d.moveCursor(2, 46)
	fmt.Fprint(d.out, green+fmt.Sprintf("%2d", running)+reset)

	d.moveCursor(2, 65)
	fmt.Fprint(d.out, white+time.Now().Format("15:04:05")+reset)
}

// updateTotalsBar updates the totals bar with aggregate metrics
func (d *MinimalDashboard) updateTotalsBar() {
	white := "\033[97m"
	reset := "\033[0m"

	// Format totals
	totalNetIn := components.FormatBytes(d.totalNetworkIn)
	totalNetOut := components.FormatBytes(d.totalNetworkOut)
	totalNet := components.FormatBytes(d.totalNetworkIn + d.totalNetworkOut)

	d.moveCursor(4, 25)
	fmt.Fprint(d.out, white+fmt.Sprintf("%-8d", d.totalRequests)+reset)

	d.moveCursor(4, 49)
	fmt.Fprint(d.out, white+fmt.Sprintf("%-10s", totalNetIn)+reset)

	d.moveCursor(4, 69)
	fmt.Fprint(d.out, white+fmt.Sprintf("%-10s", totalNetOut)+reset)

	d.moveCursor(4, 93)
	fmt.Fprint(d.out, white+fmt.Sprintf("%-10s", totalNet)+reset)
}

// updateServiceList updates service list values only
func (d *MinimalDashboard) updateServiceList() {
	white := "\033[97m"
	green := "\033[32m"
	red := "\033[31m"
	yellow := "\033[33m"
	dim := "\033[90m"
	reset := "\033[0m"

	for i := 0; i < 14; i++ {
		row := 7 + i
		d.moveCursor(row, 3)

		if i < len(d.services) {
			service := d.services[i]

			// Selection indicator
			indicator := "  "
			nameColor := dim
			if i == d.selected {
				indicator = yellow + "▸ " + reset
				nameColor = white
			}

			// Status icon
			statusIcon := dim + "○" + reset
			switch service.Status {
			case app.ServiceStatusRunning:
				statusIcon = green + "●" + reset
			case app.ServiceStatusStopped:
				statusIcon = red + "●" + reset
			case app.ServiceStatusFailed:
				statusIcon = red + "✗" + reset
			}

			// Service name (truncated)
			name := service.Name
			if len(name) > 30 {
				name = name[:27] + "..."
			}

			line := indicator + statusIcon + " " + nameColor + fmt.Sprintf("%-30s", name) + reset
			fmt.Fprint(d.out, line)
		} else {
			fmt.Fprint(d.out, strings.Repeat(" ", 37))
		}
	}
}

// updateDetails updates detail values only
func (d *MinimalDashboard) updateDetails() {
	white := "\033[97m"
	cyan := "\033[36m"
	green := "\033[32m"
	dim := "\033[90m"
	reset := "\033[0m"

	detailStart := 46

	if d.selected >= len(d.services) || len(d.services) == 0 {
		d.moveCursor(7, detailStart+8)
		fmt.Fprint(d.out, dim+"No service selected"+reset+strings.Repeat(" ", 30))
		for i := 8; i <= 12; i++ {
			d.moveCursor(i, detailStart+8)
			fmt.Fprint(d.out, strings.Repeat(" ", 40))
		}
		return
	}

	service := d.services[d.selected]

	// Row 7: Name
	d.moveCursor(7, detailStart+8)
	name := service.Name
	if len(name) > 45 {
		name = name[:42] + "..."
	}
	fmt.Fprint(d.out, white+fmt.Sprintf("%-45s", name)+reset)

	// Row 8: Type + Status
	d.moveCursor(8, detailStart+8)
	typeIcon := d.getTypeIcon(service.Type)
	fmt.Fprint(d.out, typeIcon+" "+cyan+fmt.Sprintf("%-15s", string(service.Type))+reset)

	d.moveCursor(8, detailStart+43)
	statusColor := dim
	if service.Status == app.ServiceStatusRunning {
		statusColor = green
	}
	if service.Status == app.ServiceStatusFailed {
		statusColor = "\033[31m"
	}
	fmt.Fprint(d.out, statusColor+fmt.Sprintf("%-15s", string(service.Status))+reset)

	// Row 9: ID
	d.moveCursor(9, detailStart+8)
	id := service.ID
	if len(id) > 50 {
		id = id[:47] + "..."
	}
	fmt.Fprint(d.out, dim+fmt.Sprintf("%-50s", id)+reset)

	// Row 10: Uptime + Ports
	d.moveCursor(10, detailStart+10)
	uptime := time.Since(service.CreatedAt)
	fmt.Fprint(d.out, green+fmt.Sprintf("%-18s", components.FormatDuration(int64(uptime.Seconds())))+reset)

	d.moveCursor(10, detailStart+42)
	ports := "-"
	if len(service.Ports) > 0 {
		ports = strings.Join(service.Ports, ", ")
		if len(ports) > 20 {
			ports = ports[:17] + "..."
		}
	}
	fmt.Fprint(d.out, cyan+fmt.Sprintf("%-18s", ports)+reset)

	// Row 11: Threads & Connections
	threads := "0"
	conns := "0"
	if service.Metrics != nil {
		threads = fmt.Sprintf("%d", service.Metrics.ThreadCount)
		conns = fmt.Sprintf("%d", service.Metrics.Connections)
	}
	d.moveCursor(11, detailStart+11)
	fmt.Fprint(d.out, white+fmt.Sprintf("%-15s", threads)+reset)
	d.moveCursor(11, detailStart+45)
	fmt.Fprint(d.out, white+fmt.Sprintf("%-15s", conns)+reset)

	// Row 12: Requests per sec
	rps := "0.0"
	if service.Metrics != nil && service.Metrics.RequestsPerSec > 0 {
		rps = fmt.Sprintf("%.1f", service.Metrics.RequestsPerSec)
	}
	d.moveCursor(12, detailStart+16)
	fmt.Fprint(d.out, white+fmt.Sprintf("%-10s", rps)+reset)
}

// updateMetrics updates metric values only
func (d *MinimalDashboard) updateMetrics() {
	white := "\033[97m"
	green := "\033[32m"
	red := "\033[31m"
	yellow := "\033[33m"
	dim := "\033[90m"
	reset := "\033[0m"

	detailStart := 46

	if d.selected >= len(d.services) || len(d.services) == 0 {
		return
	}

	service := d.services[d.selected]
	if service.Metrics == nil {
		d.moveCursor(15, detailStart+10)
		fmt.Fprint(d.out, dim+"Collecting metrics..."+reset+strings.Repeat(" ", 30))
		for i := 16; i <= 21; i++ {
			d.moveCursor(i, detailStart+10)
			fmt.Fprint(d.out, strings.Repeat(" ", 40))
		}
		return
	}

	m := service.Metrics

	// Row 15: CPU
	d.moveCursor(15, detailStart+10)
	cpuBar := d.makeBar(m.CPUPercent, 30)
	cpuColor := green
	if m.CPUPercent > 70 {
		cpuColor = yellow
	}
	if m.CPUPercent > 90 {
		cpuColor = red
	}
	fmt.Fprint(d.out, cpuBar+" "+cpuColor+fmt.Sprintf("%6.1f%%", m.CPUPercent)+reset+strings.Repeat(" ", 5))

	// Row 16: Memory
	d.moveCursor(16, detailStart+10)
	memBar := d.makeBar(m.MemoryPercent, 30)
	memColor := green
	if m.MemoryPercent > 70 {
		memColor = yellow
	}
	if m.MemoryPercent > 90 {
		memColor = red
	}
	memStr := components.FormatBytes(m.MemoryUsage)
	fmt.Fprint(d.out, memBar+" "+memColor+fmt.Sprintf("%6.1f%%", m.MemoryPercent)+reset+" ("+dim+memStr+reset+")"+strings.Repeat(" ", 5))

	// Row 17: Disk I/O
	d.moveCursor(17, detailStart+10)
	diskRead := components.FormatBytes(m.DiskIORead)
	diskWrite := components.FormatBytes(m.DiskIOWrite)
	fmt.Fprint(d.out, green+"Read: "+white+fmt.Sprintf("%-10s", diskRead)+reset+
		red+" Write: "+white+fmt.Sprintf("%-10s", diskWrite)+reset+strings.Repeat(" ", 10))

	// Row 18: Network
	d.moveCursor(18, detailStart+10)
	netIn := components.FormatBytes(m.NetworkIn)
	netOut := components.FormatBytes(m.NetworkOut)
	fmt.Fprint(d.out, green+"↓ "+white+fmt.Sprintf("%-10s", netIn)+reset+
		red+" ↑ "+white+fmt.Sprintf("%-10s", netOut)+reset+strings.Repeat(" ", 10))

	// Row 19: Uptime from metrics
	d.moveCursor(19, detailStart+10)
	fmt.Fprint(d.out, white+fmt.Sprintf("%-30s", components.FormatDuration(int64(m.Uptime.Seconds())))+reset)

	// Row 20: Threads & Connections
	d.moveCursor(20, detailStart+11)
	fmt.Fprint(d.out, white+fmt.Sprintf("%-15d", m.ThreadCount)+reset)
	d.moveCursor(20, detailStart+45)
	fmt.Fprint(d.out, white+fmt.Sprintf("%-15d", m.Connections)+reset)

	// Row 21: Total network for this service
	d.moveCursor(21, detailStart+17)
	totalSvcNet := components.FormatBytes(m.NetworkIn + m.NetworkOut)
	fmt.Fprint(d.out, white+fmt.Sprintf("%-20s", totalSvcNet)+reset)
}

// updateFooter updates footer values only
func (d *MinimalDashboard) updateFooter() {
	white := "\033[97m"
	reset := "\033[0m"

	d.moveCursor(24, 76)
	fmt.Fprint(d.out, white+time.Now().Format("15:04:05")+reset)
}

// makeBar creates a progress bar
func (d *MinimalDashboard) makeBar(percent float64, width int) string {
	green := "\033[32m"
	yellow := "\033[33m"
	red := "\033[31m"
	dim := "\033[90m"
	reset := "\033[0m"

	filled := int((percent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	color := green
	if percent > 60 {
		color = yellow
	}
	if percent > 85 {
		color = red
	}

	return color + strings.Repeat("█", filled) + dim + strings.Repeat("░", width-filled) + reset
}

// getTypeIcon returns icon for service type
func (d *MinimalDashboard) getTypeIcon(t app.ServiceType) string {
	switch t {
	case app.ServiceTypeDocker:
		return "🐳"
	case app.ServiceTypeKubernetes:
		return "☸️"
	case app.ServiceTypeProcess:
		return "⚡"
	case app.ServiceTypeSystemd:
		return "⚙️"
	default:
		return "📦"
	}
}

func (d *MinimalDashboard) handleMouseClick(x, y int) {
	// Check if click is in service list area (rows 7-20, cols 1-42)
	row := y + 1 // Convert 0-based to 1-based
	col := x + 1

	if row >= 7 && row <= 20 && col >= 1 && col <= 42 {
		// Calculate which service was clicked
		serviceIndex := row - 7
		if serviceIndex >= 0 && serviceIndex < len(d.services) {
			d.selected = serviceIndex
			d.updateServiceList()
			d.updateDetails()
			d.updateMetrics()
			d.out.Flush()
		}
	}
}

func (d *MinimalDashboard) navigateUp() {
	if len(d.services) > 0 && d.selected > 0 {
		d.selected--
		d.updateServiceList()
		d.updateDetails()
		d.updateMetrics()
		d.out.Flush()
	}
}

func (d *MinimalDashboard) navigateDown() {
	if len(d.services) > 0 && d.selected < len(d.services)-1 {
		d.selected++
		d.updateServiceList()
		d.updateDetails()
		d.updateMetrics()
		d.out.Flush()
	}
}
