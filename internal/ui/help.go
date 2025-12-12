package ui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HelpModal creates a help modal dialog
func (d *Dashboard) showHelpModal() {
	helpText := `[yellow]LazyService - Keyboard Shortcuts[white]

[yellow]Navigation:[white]
  ↑/↓, j/k    Navigate service list
  Enter       Select service
  Tab         Switch between panels
  Esc         Close dialogs

[yellow]Actions:[white]
  r, R        Content-only refresh (borders untouched)
  t, T        Toggle statistics/metrics view
  l, L        View logs (coming soon)
  s, S        Start/Stop service (coming soon)
  e, E        Edit configuration (coming soon)
  d, D        Delete service (coming soon)
  
[yellow]View:[white]
  f, F        Toggle full screen
  c, C        Toggle color scheme
  h, H, ?     Show this help
  
[yellow]Application:[white]
  q, Q        Quit application
  Ctrl+C      Force quit

[yellow]Border-Locked Mode:[white]
  Borders are PERMANENTLY locked after creation.
  Only text content inside borders can update.
  Navigation updates details instantly.
  Press 'R' to refresh content only.

[gray]Press any key to close this help...`

	modal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			d.tviewApp.SetRoot(d.rootLayout, true)
		})

	modal.SetBorder(true).SetTitle(" Help ")
	d.tviewApp.SetRoot(modal, false)
}

// Enhanced input handling with more shortcuts
func (d *Dashboard) setupInputHandling() {
	d.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			// Handle escape key for closing dialogs
			return nil
		case tcell.KeyTab:
			// Switch focus between panels
			d.switchFocus()
			return nil
		}

		switch event.Rune() {
		case 'q', 'Q':
			d.tviewApp.Stop()
			return nil
		case 'r', 'R':
			d.app.RefreshServices()
			d.contentOnlyRefresh() // Use content-only refresh that never touches borders
			d.showStatusMessage("Content refreshed (borders untouched)", "green")
			return nil
		case 't', 'T':
			d.toggleStatisticsView()
			return nil
		case 'l', 'L':
			d.showLogsModal()
			return nil
		case 'h', 'H', '?':
			d.showHelpModal()
			return nil
		case 'f', 'F':
			d.toggleFullScreen()
			return nil
		case 'c', 'C':
			d.toggleColorScheme()
			return nil
		case 's', 'S':
			d.showServiceActionModal("start/stop")
			return nil
		case 'e', 'E':
			d.showServiceActionModal("edit")
			return nil
		case 'd', 'D':
			d.showServiceActionModal("delete")
			return nil
		case 'j':
			// Vim-style navigation
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			// Vim-style navigation
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		}
		return event
	})
}

// switchFocus switches focus between panels
func (d *Dashboard) switchFocus() {
	// Implementation for switching focus between service list and details
	// This would require tracking current focus state
}

// toggleFullScreen toggles full screen mode for the current panel
func (d *Dashboard) toggleFullScreen() {
	// Implementation for toggling full screen
	d.showStatusMessage("Full screen toggle (coming soon)", "yellow")
}

// toggleColorScheme toggles between color schemes
func (d *Dashboard) toggleColorScheme() {
	// Implementation for color scheme switching
	d.showStatusMessage("Color scheme toggle (coming soon)", "yellow")
}

// showLogsModal shows the logs modal for the selected service
func (d *Dashboard) showLogsModal() {
	if d.selectedIndex >= len(d.services) || len(d.services) == 0 {
		d.showStatusMessage("No service selected", "red")
		return
	}
	
	service := d.services[d.selectedIndex]
	d.showStatusMessage("Logs for "+service.Name+" (coming soon)", "yellow")
}

// showServiceActionModal shows action modal for service operations
func (d *Dashboard) showServiceActionModal(action string) {
	if d.selectedIndex >= len(d.services) || len(d.services) == 0 {
		d.showStatusMessage("No service selected", "red")
		return
	}
	
	service := d.services[d.selectedIndex]
	d.showStatusMessage(action+" for "+service.Name+" (coming soon)", "yellow")
}

// showStatusMessage shows a temporary status message
func (d *Dashboard) showStatusMessage(message, color string) {
	originalText := d.statusBar.GetText(false)
	d.statusBar.SetText("[" + color + "]" + message + "[white]")
	
	// Restore original text after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		d.tviewApp.QueueUpdateDraw(func() {
			d.statusBar.SetText(originalText)
		})
	}()
}
