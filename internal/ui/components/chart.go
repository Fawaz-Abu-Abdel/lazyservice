package components

import (
	"fmt"
	"strings"
)

// CreateBarChart creates a simple ASCII bar chart
func CreateBarChart(value float64, max float64, width int) string {
	if max == 0 {
		max = 100
	}
	
	filled := int((value / max) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	bar := strings.Repeat("█", filled)
	empty := strings.Repeat("░", width-filled)
	return bar + empty
}

// CreateSparkline creates a simple sparkline from data points
func CreateSparkline(data []float64, width int) string {
	if len(data) == 0 {
		return strings.Repeat("▁", width)
	}

	// Find min and max
	min, max := data[0], data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	// Create sparkline
	chars := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	result := ""

	// Sample data points if we have more than width
	step := len(data) / width
	if step < 1 {
		step = 1
	}

	for i := 0; i < width; i++ {
		var value float64
		if len(data) > 0 {
			idx := (i * len(data)) / width
			if idx >= len(data) {
				idx = len(data) - 1
			}
			value = data[idx]
		}

		normalized := 0.0
		if max > min {
			normalized = (value - min) / (max - min)
		}

		charIdx := int(normalized * float64(len(chars)-1))
		if charIdx >= len(chars) {
			charIdx = len(chars) - 1
		}
		if charIdx < 0 {
			charIdx = 0
		}
		result += string(chars[charIdx])
	}

	return result
}

// FormatBytes formats bytes into human-readable format
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration formats duration into human-readable format
func FormatDuration(d int64) string {
	days := d / 86400
	hours := (d % 86400) / 3600
	minutes := (d % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// CreateMiniChart creates a mini chart for metrics
func CreateMiniChart(data []float64, height int, width int) []string {
	if len(data) == 0 || height == 0 || width == 0 {
		return []string{}
	}

	// Find max value for scaling
	max := data[0]
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	if max == 0 {
		max = 1
	}

	// Sample data points if needed
	sampledData := make([]float64, width)
	step := float64(len(data)) / float64(width)
	for i := 0; i < width; i++ {
		idx := int(float64(i) * step)
		if idx >= len(data) {
			idx = len(data) - 1
		}
		sampledData[i] = data[idx]
	}

	// Create chart rows
	lines := make([]string, height)
	for row := 0; row < height; row++ {
		line := ""
		threshold := max * float64(height-row) / float64(height)
		for _, value := range sampledData {
			if value >= threshold {
				line += "█"
			} else {
				line += " "
			}
		}
		lines[row] = line
	}

	return lines
}
