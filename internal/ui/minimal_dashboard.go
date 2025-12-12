package ui

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"github.com/lazyservice/lazyservice/internal/app"
	"github.com/lazyservice/lazyservice/internal/ui/components"
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
	
	// Console handles
	stdinHandle  windows.Handle
	originalMode uint32
	
	// Totals tracking
	totalRequests   uint64
	totalNetworkIn  uint64
	totalNetworkOut uint64
}

// NewMinimalDashboard creates a minimal dashboard
func NewMinimalDashboard(application *app.App) *MinimalDashboard {
	return &MinimalDashboard{
		app:      application,
		services: make([]*app.Service, 0),
		selected: 0,
		running:  true,
	}
}

// Run starts the minimal dashboard
func (d *MinimalDashboard) Run() error {
	d.getTerminalSize()
	d.enableVirtualTerminal()
	d.enableMouseInput()
	defer d.restoreTerminal()
	
	d.hideCursor()
	defer d.showCursor()
	
	// Draw borders ONCE
	d.drawBorders()
	
	// Initial data
	d.services = d.app.GetServices()
	d.calculateTotals()
	d.updateAllValues()
	
	// Background update loop
	go d.updateLoop()
	
	// Handle keyboard and mouse input
	return d.handleInput()
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
	stdout := windows.Handle(os.Stdout.Fd())
	var csbi windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(stdout, &csbi); err == nil {
		d.termWidth = int(csbi.Window.Right - csbi.Window.Left + 1)
		d.termHeight = int(csbi.Window.Bottom - csbi.Window.Top + 1)
	} else {
		d.termWidth = 120
		d.termHeight = 30
	}
}

func (d *MinimalDashboard) enableVirtualTerminal() {
	stdout := windows.Handle(os.Stdout.Fd())
	var mode uint32
	windows.GetConsoleMode(stdout, &mode)
	windows.SetConsoleMode(stdout, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

func (d *MinimalDashboard) enableMouseInput() {
	d.stdinHandle = windows.Handle(os.Stdin.Fd())
	windows.GetConsoleMode(d.stdinHandle, &d.originalMode)
	// Enable mouse input and disable line input/echo
	newMode := (d.originalMode | windows.ENABLE_MOUSE_INPUT | windows.ENABLE_EXTENDED_FLAGS) &^ 
		(windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT | windows.ENABLE_QUICK_EDIT_MODE)
	windows.SetConsoleMode(d.stdinHandle, newMode)
}

func (d *MinimalDashboard) restoreTerminal() {
	if d.stdinHandle != 0 {
		windows.SetConsoleMode(d.stdinHandle, d.originalMode)
	}
}

func (d *MinimalDashboard) hideCursor() { fmt.Print("\033[?25l") }
func (d *MinimalDashboard) showCursor() { fmt.Print("\033[?25h") }
func (d *MinimalDashboard) clearScreen() { fmt.Print("\033[2J\033[H") }
func (d *MinimalDashboard) moveCursor(row, col int) { fmt.Printf("\033[%d;%dH", row, col) }

// drawBorders draws all borders ONCE - never called again
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
	fmt.Print(cyan + "╔" + strings.Repeat("═", d.termWidth-2) + "╗" + reset)
	d.moveCursor(2, 1)
	fmt.Print(cyan + "║" + strings.Repeat(" ", d.termWidth-2) + "║" + reset)
	d.moveCursor(3, 1)
	fmt.Print(cyan + "╠" + strings.Repeat("═", d.termWidth-2) + "╣" + reset)
	
	// Totals bar (row 4)
	d.moveCursor(4, 1)
	fmt.Print(cyan + "║" + strings.Repeat(" ", d.termWidth-2) + "║" + reset)
	d.moveCursor(5, 1)
	fmt.Print(cyan + "╚" + strings.Repeat("═", d.termWidth-2) + "╝" + reset)
	
	// Services panel (left, rows 6-22)
	serviceWidth := 42
	d.moveCursor(6, 1)
	fmt.Print(yellow + "┌─ Services ─" + strings.Repeat("─", serviceWidth-13) + "┐" + reset)
	for i := 7; i <= 21; i++ {
		d.moveCursor(i, 1)
		fmt.Print(yellow + "│" + strings.Repeat(" ", serviceWidth-1) + "│" + reset)
	}
	d.moveCursor(22, 1)
	fmt.Print(yellow + "└" + strings.Repeat("─", serviceWidth-1) + "┘" + reset)
	
	// Details panel (right top, rows 6-13)
	detailStart := serviceWidth + 2
	detailWidth := d.termWidth - serviceWidth - 3
	d.moveCursor(6, detailStart)
	fmt.Print(magenta + "┌─ Details " + strings.Repeat("─", detailWidth-12) + "┐" + reset)
	for i := 7; i <= 12; i++ {
		d.moveCursor(i, detailStart)
		fmt.Print(magenta + "│" + strings.Repeat(" ", detailWidth-1) + "│" + reset)
	}
	d.moveCursor(13, detailStart)
	fmt.Print(magenta + "└" + strings.Repeat("─", detailWidth-1) + "┘" + reset)
	
	// Metrics panel (right middle, rows 14-22)
	d.moveCursor(14, detailStart)
	fmt.Print(green + "┌─ Live Metrics " + strings.Repeat("─", detailWidth-17) + "┐" + reset)
	for i := 15; i <= 21; i++ {
		d.moveCursor(i, detailStart)
		fmt.Print(green + "│" + strings.Repeat(" ", detailWidth-1) + "│" + reset)
	}
	d.moveCursor(22, detailStart)
	fmt.Print(green + "└" + strings.Repeat("─", detailWidth-1) + "┘" + reset)
	
	// Footer (row 23)
	d.moveCursor(23, 1)
	fmt.Print(blue + strings.Repeat("─", d.termWidth) + reset)
	d.moveCursor(24, 1)
	fmt.Print(dim + strings.Repeat(" ", d.termWidth) + reset)
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
	cyan := "\033[36m"
	yellow := "\033[33m"
	green := "\033[32m"
	reset := "\033[0m"
	
	d.moveCursor(2, 3)
	
	// Count running services
	running := 0
	for _, s := range d.services {
		if s.Status == app.ServiceStatusRunning {
			running++
		}
	}
	
	title := cyan + "⚡ LazyService Pro" + reset + "  │  " + 
		yellow + fmt.Sprintf("%d", len(d.services)) + reset + " services  │  " +
		green + fmt.Sprintf("%d", running) + reset + " running  │  " +
		white + time.Now().Format("15:04:05") + reset
	fmt.Print(title + strings.Repeat(" ", 40))
}

// updateTotalsBar updates the totals bar with aggregate metrics
func (d *MinimalDashboard) updateTotalsBar() {
	white := "\033[97m"
	cyan := "\033[36m"
	yellow := "\033[33m"
	green := "\033[32m"
	magenta := "\033[35m"
	reset := "\033[0m"
	
	d.moveCursor(4, 3)
	
	// Format totals
	totalNetIn := components.FormatBytes(d.totalNetworkIn)
	totalNetOut := components.FormatBytes(d.totalNetworkOut)
	totalNet := components.FormatBytes(d.totalNetworkIn + d.totalNetworkOut)
	
	totals := cyan + "📊 Totals:" + reset + "  " +
		yellow + "Requests: " + white + fmt.Sprintf("%d", d.totalRequests) + reset + "  │  " +
		green + "Net↓: " + white + totalNetIn + reset + "  " +
		magenta + "Net↑: " + white + totalNetOut + reset + "  │  " +
		cyan + "Total Net: " + white + totalNet + reset
	
	fmt.Print(totals + strings.Repeat(" ", 30))
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
			fmt.Print(line)
		} else {
			fmt.Print(strings.Repeat(" ", 37))
		}
	}
}

// updateDetails updates detail values only
func (d *MinimalDashboard) updateDetails() {
	white := "\033[97m"
	yellow := "\033[33m"
	cyan := "\033[36m"
	green := "\033[32m"
	dim := "\033[90m"
	reset := "\033[0m"
	
	detailStart := 46
	
	if d.selected >= len(d.services) || len(d.services) == 0 {
		for i := 7; i <= 12; i++ {
			d.moveCursor(i, detailStart)
			fmt.Print(dim + "No service selected" + reset + strings.Repeat(" ", 40))
		}
		return
	}
	
	service := d.services[d.selected]
	
	// Row 7: Name
	d.moveCursor(7, detailStart)
	name := service.Name
	if len(name) > 45 { name = name[:42] + "..." }
	fmt.Print(yellow + "Name: " + reset + white + fmt.Sprintf("%-55s", name) + reset)
	
	// Row 8: Type + Status
	d.moveCursor(8, detailStart)
	typeIcon := d.getTypeIcon(service.Type)
	statusColor := dim
	if service.Status == app.ServiceStatusRunning { statusColor = green }
	if service.Status == app.ServiceStatusFailed { statusColor = "\033[31m" }
	fmt.Print(yellow + "Type: " + reset + typeIcon + " " + cyan + string(service.Type) + reset + 
		"  │  " + yellow + "Status: " + reset + statusColor + string(service.Status) + reset + strings.Repeat(" ", 20))
	
	// Row 9: ID
	d.moveCursor(9, detailStart)
	id := service.ID
	if len(id) > 50 { id = id[:47] + "..." }
	fmt.Print(yellow + "ID: " + reset + dim + fmt.Sprintf("%-57s", id) + reset)
	
	// Row 10: Uptime + Ports
	d.moveCursor(10, detailStart)
	uptime := time.Since(service.CreatedAt)
	ports := "-"
	if len(service.Ports) > 0 {
		ports = strings.Join(service.Ports, ", ")
		if len(ports) > 20 { ports = ports[:17] + "..." }
	}
	fmt.Print(yellow + "Uptime: " + reset + green + components.FormatDuration(int64(uptime.Seconds())) + reset +
		"  │  " + yellow + "Ports: " + reset + cyan + ports + reset + strings.Repeat(" ", 20))
	
	// Row 11: Threads & Connections
	d.moveCursor(11, detailStart)
	threads := "0"
	conns := "0"
	if service.Metrics != nil {
		threads = fmt.Sprintf("%d", service.Metrics.ThreadCount)
		conns = fmt.Sprintf("%d", service.Metrics.Connections)
	}
	fmt.Print(yellow + "Threads: " + reset + white + threads + reset +
		"  │  " + yellow + "Connections: " + reset + white + conns + reset + strings.Repeat(" ", 30))
	
	// Row 12: Requests per sec
	d.moveCursor(12, detailStart)
	rps := "0"
	if service.Metrics != nil && service.Metrics.RequestsPerSec > 0 {
		rps = fmt.Sprintf("%.1f", service.Metrics.RequestsPerSec)
	}
	fmt.Print(yellow + "Requests/sec: " + reset + white + rps + reset + strings.Repeat(" ", 45))
}

// updateMetrics updates metric values only
func (d *MinimalDashboard) updateMetrics() {
	white := "\033[97m"
	yellow := "\033[33m"
	green := "\033[32m"
	red := "\033[31m"
	dim := "\033[90m"
	reset := "\033[0m"
	
	detailStart := 46
	
	if d.selected >= len(d.services) || len(d.services) == 0 {
		for i := 15; i <= 21; i++ {
			d.moveCursor(i, detailStart)
			fmt.Print(strings.Repeat(" ", 65))
		}
		return
	}
	
	service := d.services[d.selected]
	if service.Metrics == nil {
		d.moveCursor(15, detailStart)
		fmt.Print(dim + "Collecting metrics..." + reset + strings.Repeat(" ", 45))
		for i := 16; i <= 21; i++ {
			d.moveCursor(i, detailStart)
			fmt.Print(strings.Repeat(" ", 65))
		}
		return
	}
	
	m := service.Metrics
	
	// Row 15: CPU
	d.moveCursor(15, detailStart)
	cpuBar := d.makeBar(m.CPUPercent, 30)
	cpuColor := green
	if m.CPUPercent > 70 { cpuColor = yellow }
	if m.CPUPercent > 90 { cpuColor = red }
	fmt.Print(yellow + "CPU:    " + reset + cpuBar + " " + cpuColor + fmt.Sprintf("%6.1f%%", m.CPUPercent) + reset + strings.Repeat(" ", 10))
	
	// Row 16: Memory
	d.moveCursor(16, detailStart)
	memBar := d.makeBar(m.MemoryPercent, 30)
	memColor := green
	if m.MemoryPercent > 70 { memColor = yellow }
	if m.MemoryPercent > 90 { memColor = red }
	memStr := components.FormatBytes(m.MemoryUsage)
	fmt.Print(yellow + "Memory: " + reset + memBar + " " + memColor + fmt.Sprintf("%6.1f%%", m.MemoryPercent) + reset + " (" + dim + memStr + reset + ")")
	
	// Row 17: Disk I/O
	d.moveCursor(17, detailStart)
	diskRead := components.FormatBytes(m.DiskIORead)
	diskWrite := components.FormatBytes(m.DiskIOWrite)
	fmt.Print(yellow + "Disk:   " + reset + green + "Read: " + white + fmt.Sprintf("%-10s", diskRead) + reset + 
		red + " Write: " + white + fmt.Sprintf("%-10s", diskWrite) + reset + strings.Repeat(" ", 20))
	
	// Row 18: Network
	d.moveCursor(18, detailStart)
	netIn := components.FormatBytes(m.NetworkIn)
	netOut := components.FormatBytes(m.NetworkOut)
	fmt.Print(yellow + "Net:    " + reset + green + "↓ " + white + fmt.Sprintf("%-10s", netIn) + reset + 
		red + " ↑ " + white + fmt.Sprintf("%-10s", netOut) + reset + strings.Repeat(" ", 20))
	
	// Row 19: Uptime from metrics
	d.moveCursor(19, detailStart)
	fmt.Print(yellow + "Uptime: " + reset + white + components.FormatDuration(int64(m.Uptime.Seconds())) + reset + strings.Repeat(" ", 45))
	
	// Row 20: Threads & Connections
	d.moveCursor(20, detailStart)
	fmt.Print(yellow + "Threads: " + reset + white + fmt.Sprintf("%d", m.ThreadCount) + reset + 
		"  │  " + yellow + "Connections: " + reset + white + fmt.Sprintf("%d", m.Connections) + reset + strings.Repeat(" ", 30))
	
	// Row 21: Total network for this service
	d.moveCursor(21, detailStart)
	totalSvcNet := components.FormatBytes(m.NetworkIn + m.NetworkOut)
	fmt.Print(yellow + "Total Network: " + reset + white + totalSvcNet + reset + strings.Repeat(" ", 45))
}

// updateFooter updates footer values only
func (d *MinimalDashboard) updateFooter() {
	white := "\033[97m"
	cyan := "\033[36m"
	green := "\033[32m"
	yellow := "\033[33m"
	dim := "\033[90m"
	reset := "\033[0m"
	
	d.moveCursor(24, 1)
	footer := cyan + "↑↓" + dim + "/" + cyan + "jk" + reset + " Navigate  " +
		yellow + "🖱️ Click" + reset + " Select  " +
		green + "R" + reset + " Refresh  " +
		"\033[31mQ" + reset + " Quit  " +
		dim + "│ Updated: " + white + time.Now().Format("15:04:05") + reset
	fmt.Print(footer + strings.Repeat(" ", 20))
}

// makeBar creates a progress bar
func (d *MinimalDashboard) makeBar(percent float64, width int) string {
	green := "\033[32m"
	yellow := "\033[33m"
	red := "\033[31m"
	dim := "\033[90m"
	reset := "\033[0m"
	
	filled := int((percent / 100.0) * float64(width))
	if filled > width { filled = width }
	if filled < 0 { filled = 0 }
	
	color := green
	if percent > 60 { color = yellow }
	if percent > 85 { color = red }
	
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

// INPUT_RECORD structure for Windows console input
type inputRecord struct {
	EventType uint16
	_         uint16
	Event     [16]byte
}

type keyEventRecord struct {
	KeyDown         int32
	RepeatCount     uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	Char            uint16
	ControlKeyState uint32
}

type mouseEventRecord struct {
	MousePosition struct {
		X int16
		Y int16
	}
	ButtonState     uint32
	ControlKeyState uint32
	EventFlags      uint32
}

// handleInput handles keyboard and mouse input
func (d *MinimalDashboard) handleInput() error {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	readConsoleInput := kernel32.NewProc("ReadConsoleInputW")
	
	var inputRecords [10]inputRecord
	var numRead uint32
	
	for d.running {
		ret, _, _ := readConsoleInput.Call(
			uintptr(d.stdinHandle),
			uintptr(unsafe.Pointer(&inputRecords[0])),
			10,
			uintptr(unsafe.Pointer(&numRead)),
		)
		
		if ret == 0 {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		
		for i := uint32(0); i < numRead; i++ {
			record := inputRecords[i]
			
			switch record.EventType {
			case 0x0001: // KEY_EVENT
				keyEvent := (*keyEventRecord)(unsafe.Pointer(&record.Event[0]))
				if keyEvent.KeyDown != 0 {
					d.handleKeyEvent(keyEvent)
				}
				
			case 0x0002: // MOUSE_EVENT
				mouseEvent := (*mouseEventRecord)(unsafe.Pointer(&record.Event[0]))
				d.handleMouseEvent(mouseEvent)
			}
		}
	}
	
	d.clearScreen()
	d.moveCursor(1, 1)
	return nil
}

func (d *MinimalDashboard) handleKeyEvent(event *keyEventRecord) {
	// Virtual key codes
	const (
		VK_UP    = 0x26
		VK_DOWN  = 0x28
		VK_ESCAPE = 0x1B
	)
	
	switch event.VirtualKeyCode {
	case VK_UP:
		d.navigateUp()
	case VK_DOWN:
		d.navigateDown()
	case VK_ESCAPE:
		d.running = false
	default:
		// Check character
		char := rune(event.Char)
		switch char {
		case 'q', 'Q':
			d.running = false
		case 'r', 'R':
			d.app.RefreshServices()
			d.services = d.app.GetServices()
			d.calculateTotals()
			d.updateAllValues()
		case 'k', 'K':
			d.navigateUp()
		case 'j', 'J':
			d.navigateDown()
		}
	}
}

func (d *MinimalDashboard) handleMouseEvent(event *mouseEventRecord) {
	// Check for mouse click (left button)
	if event.ButtonState&0x0001 != 0 {
		// Check if click is in service list area (rows 7-20, cols 1-42)
		row := int(event.MousePosition.Y) + 1 // Convert 0-based to 1-based
		col := int(event.MousePosition.X) + 1
		
		if row >= 7 && row <= 20 && col >= 1 && col <= 42 {
			// Calculate which service was clicked
			serviceIndex := row - 7
			if serviceIndex >= 0 && serviceIndex < len(d.services) {
				d.selected = serviceIndex
				d.updateServiceList()
				d.updateDetails()
				d.updateMetrics()
			}
		}
	}
}

func (d *MinimalDashboard) navigateUp() {
	if len(d.services) > 0 && d.selected > 0 {
		d.selected--
		d.updateServiceList()
		d.updateDetails()
		d.updateMetrics()
	}
}

func (d *MinimalDashboard) navigateDown() {
	if len(d.services) > 0 && d.selected < len(d.services)-1 {
		d.selected++
		d.updateServiceList()
		d.updateDetails()
		d.updateMetrics()
	}
}
