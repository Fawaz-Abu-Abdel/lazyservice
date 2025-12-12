package components

import (
	"fmt"
	"strings"
	"time"
)

// CreateAdvancedBarChart creates a beautiful bar chart with gradients
func CreateAdvancedBarChart(percentage float64, max float64, width int, colors []string) string {
	if width <= 0 {
		return ""
	}
	
	filled := int((percentage / max) * float64(width))
	if filled > width {
		filled = width
	}
	
	var result strings.Builder
	
	// Create gradient effect
	for i := 0; i < width; i++ {
		if i < filled {
			// Determine color based on position and value
			colorIndex := 0
			if percentage > 80 {
				colorIndex = 2 // Red
			} else if percentage > 60 {
				colorIndex = 1 // Yellow
			} else {
				colorIndex = 0 // Green
			}
			
			if colorIndex < len(colors) {
				result.WriteString("[" + colors[colorIndex] + "]█[white]")
			} else {
				result.WriteString("█")
			}
		} else {
			result.WriteString("[gray]░[white]")
		}
	}
	
	return result.String()
}

// CreateSparklineChart creates a beautiful sparkline with colors
func CreateSparklineChart(data []float64, width int, baseColor string) string {
	if len(data) == 0 || width <= 0 {
		return "[gray]" + strings.Repeat("▁", width) + "[white]"
	}
	
	// Find min and max for scaling
	min, max := data[0], data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	
	// Avoid division by zero
	if max == min {
		return "[" + baseColor + "]" + strings.Repeat("▄", width) + "[white]"
	}
	
	// Create sparkline with color coding
	chars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	var result strings.Builder
	
	// Sample data to fit width
	step := float64(len(data)) / float64(width)
	for i := 0; i < width; i++ {
		dataIndex := int(float64(i) * step)
		if dataIndex >= len(data) {
			dataIndex = len(data) - 1
		}
		
		value := data[dataIndex]
		
		// Normalize to 0-7 range
		normalized := (value - min) / (max - min) * 7
		charIndex := int(normalized)
		if charIndex >= len(chars) {
			charIndex = len(chars) - 1
		}
		
		// Color based on value
		color := baseColor
		if value > (max * 0.8) {
			color = "red"
		} else if value > (max * 0.6) {
			color = "yellow"
		} else {
			color = "green"
		}
		
		result.WriteString("[" + color + "]" + chars[charIndex] + "[white]")
	}
	
	return result.String()
}

// CreateGaugeChart creates a beautiful gauge visualization
func CreateGaugeChart(value, max float64, width int) string {
	if width <= 6 {
		return "▓"
	}
	
	percentage := (value / max) * 100
	
	// Create gauge
	filled := int((percentage / 100.0) * float64(width))
	
	var gauge strings.Builder
	gauge.WriteString("[blue]╢[white]")
	
	for i := 0; i < width-2; i++ {
		if i < filled {
			color := "green"
			if percentage > 80 {
				color = "red"
			} else if percentage > 60 {
				color = "yellow"
			}
			gauge.WriteString("[" + color + "]█[white]")
		} else {
			gauge.WriteString("[gray]░[white]")
		}
	}
	
	gauge.WriteString("[blue]╟[white]")
	return gauge.String()
}

// CreateStatusBadge creates a beautiful status badge
func CreateStatusBadge(status string, text string) string {
	var color, icon string
	
	switch strings.ToLower(status) {
	case "running", "active", "up":
		color = "green"
		icon = "●"
	case "stopped", "inactive", "down":
		color = "red"
		icon = "●"
	case "warning", "degraded":
		color = "yellow"
		icon = "⚠"
	case "error", "failed":
		color = "red"
		icon = "✗"
	case "unknown":
		color = "gray"
		icon = "?"
	default:
		color = "blue"
		icon = "○"
	}
	
	return fmt.Sprintf("[%s]%s %s[white]", color, icon, text)
}

// CreateProgressRing creates a circular progress indicator
func CreateProgressRing(percentage float64) string {
	// Unicode circle characters for progress ring
	rings := []string{"○", "◔", "◐", "◕", "●"}
	
	index := int((percentage / 100.0) * float64(len(rings)-1))
	if index >= len(rings) {
		index = len(rings) - 1
	}
	
	color := "green"
	if percentage > 80 {
		color = "red"
	} else if percentage > 60 {
		color = "yellow"
	}
	
	return fmt.Sprintf("[%s]%s[white]", color, rings[index])
}

// CreateTrendIndicator creates a trend arrow indicator
func CreateTrendIndicator(current, previous float64) string {
	if current > previous*1.1 {
		return "[green]↗[white]" // Increasing
	} else if current < previous*0.9 {
		return "[red]↘[white]" // Decreasing
	} else {
		return "[blue]→[white]" // Stable
	}
}

// CreateAnimatedDots creates animated loading dots
func CreateAnimatedDots(frame int) string {
	dots := []string{"   ", ".  ", ".. ", "..."}
	return "[cyan]" + dots[frame%len(dots)] + "[white]"
}

// CreateBoxedText creates text with a beautiful box around it
func CreateBoxedText(text string, width int) string {
	if width < 4 {
		return text
	}
	
	lines := strings.Split(text, "\n")
	var result strings.Builder
	
	// Top border
	result.WriteString("╭" + strings.Repeat("─", width-2) + "╮\n")
	
	// Content lines
	for _, line := range lines {
		padding := width - 2 - len(line)
		if padding < 0 {
			line = line[:width-2]
			padding = 0
		}
		
		result.WriteString("│" + line + strings.Repeat(" ", padding) + "│\n")
	}
	
	// Bottom border
	result.WriteString("╰" + strings.Repeat("─", width-2) + "╯")
	
	return result.String()
}

// CreateMetricCard creates a beautiful metric display card
func CreateMetricCard(title, value, unit string, percentage float64, trend string) string {
	var card strings.Builder
	
	// Title with icon
	card.WriteString(fmt.Sprintf(" [cyan]%s[white]\n", title))
	
	// Value with unit
	card.WriteString(fmt.Sprintf(" [white]%s[gray] %s[white]", value, unit))
	
	// Trend indicator
	if trend != "" {
		card.WriteString(" " + trend)
	}
	
	card.WriteString("\n")
	
	// Progress bar
	if percentage >= 0 {
		bar := CreateAdvancedBarChart(percentage, 100, 15, []string{"green", "yellow", "red"})
		card.WriteString(" " + bar + fmt.Sprintf(" %.1f%%", percentage))
	}
	
	return card.String()
}

// CreateTimestamp creates a formatted timestamp with relative time
func CreateTimestamp(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	
	var relative string
	if diff < time.Minute {
		relative = "just now"
	} else if diff < time.Hour {
		relative = fmt.Sprintf("%dm ago", int(diff.Minutes()))
	} else if diff < 24*time.Hour {
		relative = fmt.Sprintf("%dh ago", int(diff.Hours()))
	} else {
		relative = fmt.Sprintf("%dd ago", int(diff.Hours()/24))
	}
	
	return fmt.Sprintf("[gray]%s[white] ([cyan]%s[white])", 
		t.Format("15:04:05"), relative)
}

// CreateNetworkIndicator creates a network activity indicator
func CreateNetworkIndicator(bytesIn, bytesOut uint64) string {
	var indicator strings.Builder
	
	// Incoming traffic
	if bytesIn > 0 {
		indicator.WriteString("[green]↓[white]")
		indicator.WriteString(FormatBytes(bytesIn))
	} else {
		indicator.WriteString("[gray]↓[white]0")
	}
	
	indicator.WriteString(" ")
	
	// Outgoing traffic
	if bytesOut > 0 {
		indicator.WriteString("[red]↑[white]")
		indicator.WriteString(FormatBytes(bytesOut))
	} else {
		indicator.WriteString("[gray]↑[white]0")
	}
	
	return indicator.String()
}

// CreateHealthIndicator creates a health status indicator
func CreateHealthIndicator(healthy bool, lastCheck time.Time) string {
	var indicator strings.Builder
	
	if healthy {
		indicator.WriteString("[green]●[white] Healthy")
	} else {
		indicator.WriteString("[red]●[white] Unhealthy")
	}
	
	// Add last check time
	if !lastCheck.IsZero() {
		diff := time.Since(lastCheck)
		if diff < time.Minute {
			indicator.WriteString(" [gray](checked now)[white]")
		} else {
			indicator.WriteString(fmt.Sprintf(" [gray](checked %s ago)[white]", 
				FormatDuration(int64(diff.Seconds()))))
		}
	}
	
	return indicator.String()
}
