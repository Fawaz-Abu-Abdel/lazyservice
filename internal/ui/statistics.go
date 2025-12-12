package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/lazyservice/lazyservice/internal/app"
	"github.com/lazyservice/lazyservice/internal/ui/components"
)

// FormatStatistics formats comprehensive statistics for display
func FormatStatistics(service *app.Service) string {
	if service.Statistics == nil {
		return "\n [gray]No statistics available - Statistics object is nil"
	}
	
	// Show building message only if absolutely no data
	if service.Statistics.TotalRequests == 0 && len(service.Statistics.CPUHistory) == 0 && len(service.Statistics.MemoryHistory) == 0 {
		uptime := time.Since(service.Statistics.StartTime)
		return fmt.Sprintf("\n [yellow]📊 Building Statistics...[white]\n\n"+
			" [gray]Service started: %s\n"+
			" Uptime: %s\n\n"+
			" Statistics will appear as the service\n"+
			" processes requests and generates metrics.\n\n"+
			" Press 'R' to refresh data.",
			service.Statistics.StartTime.Format("15:04:05"),
			components.FormatDuration(int64(uptime.Seconds())))
	}

	stats := service.Statistics
	var output strings.Builder

	output.WriteString("\n")
	output.WriteString(" [yellow]═══ SERVICE STATISTICS ═══[white]\n\n")

	// Request Statistics
	output.WriteString(" [yellow]📊 Request Statistics[white]\n")
	output.WriteString(fmt.Sprintf("   Total Requests: [cyan]%d[white]\n", stats.TotalRequests))
	output.WriteString(fmt.Sprintf("   Requests/sec:   [cyan]%.2f[white]\n", stats.CalculateRequestsPerSecond()))
	if !stats.LastRequestTime.IsZero() {
		output.WriteString(fmt.Sprintf("   Last Request:   [cyan]%s[white]\n", 
			formatTimeAgo(stats.LastRequestTime)))
	}
	output.WriteString(fmt.Sprintf("   Total Errors:   [red]%d[white]\n", stats.TotalErrors))
	if stats.TotalRequests > 0 {
		errorRate := float64(stats.TotalErrors) / float64(stats.TotalRequests) * 100
		output.WriteString(fmt.Sprintf("   Error Rate:     [red]%.2f%%[white]\n", errorRate))
	}
	output.WriteString("\n")

	// CPU Statistics (Last 5 minutes)
	output.WriteString(" [yellow]🖥️  CPU Usage (Last 5m)[white]\n")
	avgCPU := stats.GetAverageCPU(5 * time.Minute)
	peakCPU := stats.GetPeakCPU(5 * time.Minute)
	output.WriteString(fmt.Sprintf("   Average:        [cyan]%.1f%%[white]\n", avgCPU))
	output.WriteString(fmt.Sprintf("   Peak:           [cyan]%.1f%%[white]\n", peakCPU))
	
	// CPU trend visualization
	if len(stats.CPUHistory) > 0 {
		cpuTrend := createMiniChart(stats.CPUHistory, 20)
		output.WriteString(fmt.Sprintf("   Trend:          %s\n", cpuTrend))
	}
	output.WriteString("\n")

	// Memory Statistics (Last 5 minutes)
	output.WriteString(" [yellow]💾 Memory Usage (Last 5m)[white]\n")
	avgMem, avgMemPct := stats.GetAverageMemory(5 * time.Minute)
	peakMem, peakMemPct := stats.GetPeakMemory(5 * time.Minute)
	output.WriteString(fmt.Sprintf("   Average:        [cyan]%s (%.1f%%)[white]\n", 
		components.FormatBytes(avgMem), avgMemPct))
	output.WriteString(fmt.Sprintf("   Peak:           [cyan]%s (%.1f%%)[white]\n", 
		components.FormatBytes(peakMem), peakMemPct))
	
	// Memory trend visualization
	if len(stats.MemoryHistory) > 0 {
		memTrend := createMiniChartFromMemory(stats.MemoryHistory, 20)
		output.WriteString(fmt.Sprintf("   Trend:          %s\n", memTrend))
	}
	output.WriteString("\n")

	// Network Statistics
	output.WriteString(" [yellow]🌐 Network Usage[white]\n")
	output.WriteString(fmt.Sprintf("   Total In:       [cyan]%s[white]\n", 
		components.FormatBytes(stats.TotalDataIn)))
	output.WriteString(fmt.Sprintf("   Total Out:      [cyan]%s[white]\n", 
		components.FormatBytes(stats.TotalDataOut)))
	output.WriteString(fmt.Sprintf("   Rate In:        [cyan]%s/s[white]\n", 
		components.FormatBytes(stats.DataInRate)))
	output.WriteString(fmt.Sprintf("   Rate Out:       [cyan]%s/s[white]\n", 
		components.FormatBytes(stats.DataOutRate)))
	
	trend, _, _ := stats.GetNetworkTrend(5 * time.Minute)
	trendEmoji := getTrendEmoji(trend)
	output.WriteString(fmt.Sprintf("   Trend:          %s [cyan]%s[white]\n", trendEmoji, trend))
	output.WriteString("\n")

	// Uptime and Availability
	output.WriteString(" [yellow]⏱️  Uptime & Availability[white]\n")
	uptime := time.Since(stats.StartTime)
	output.WriteString(fmt.Sprintf("   Uptime:         [green]%s[white]\n", 
		components.FormatDuration(int64(uptime.Seconds()))))
	output.WriteString(fmt.Sprintf("   Started:        [cyan]%s[white]\n", 
		stats.StartTime.Format("2006-01-02 15:04:05")))
	
	if stats.TotalDowntime > 0 {
		output.WriteString(fmt.Sprintf("   Downtime:       [red]%s[white]\n", 
			components.FormatDuration(int64(stats.TotalDowntime.Seconds()))))
		availability := (1 - float64(stats.TotalDowntime)/float64(uptime)) * 100
		output.WriteString(fmt.Sprintf("   Availability:   [green]%.2f%%[white]\n", availability))
	} else {
		output.WriteString("   Availability:   [green]100.00%%[white]\n")
	}

	return output.String()
}

// createMiniChart creates a mini ASCII chart from CPU history
func createMiniChart(history []app.CPUDataPoint, width int) string {
	if len(history) == 0 {
		return "[gray]No data[white]"
	}

	// Get recent data points
	start := 0
	if len(history) > width {
		start = len(history) - width
	}
	
	data := make([]float64, 0, width)
	for i := start; i < len(history); i++ {
		data = append(data, history[i].Usage)
	}

	return createASCIIChart(data, width)
}

// createMiniChartFromMemory creates a mini ASCII chart from memory history
func createMiniChartFromMemory(history []app.MemoryDataPoint, width int) string {
	if len(history) == 0 {
		return "[gray]No data[white]"
	}

	// Get recent data points
	start := 0
	if len(history) > width {
		start = len(history) - width
	}
	
	data := make([]float64, 0, width)
	for i := start; i < len(history); i++ {
		data = append(data, history[i].Percentage)
	}

	return createASCIIChart(data, width)
}

// createASCIIChart creates a simple ASCII chart
func createASCIIChart(data []float64, width int) string {
	if len(data) == 0 {
		return "[gray]No data[white]"
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
		return strings.Repeat("─", width)
	}

	// Create chart
	var chart strings.Builder
	chars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	
	for _, value := range data {
		// Normalize to 0-7 range
		normalized := (value - min) / (max - min) * 7
		index := int(normalized)
		if index >= len(chars) {
			index = len(chars) - 1
		}
		
		// Color based on value
		if value > 80 {
			chart.WriteString("[red]" + chars[index] + "[white]")
		} else if value > 60 {
			chart.WriteString("[yellow]" + chars[index] + "[white]")
		} else {
			chart.WriteString("[green]" + chars[index] + "[white]")
		}
	}

	return chart.String()
}

// getTrendEmoji returns an emoji for the trend
func getTrendEmoji(trend string) string {
	switch trend {
	case "increasing":
		return "📈"
	case "decreasing":
		return "📉"
	default:
		return "📊"
	}
}

// formatTimeAgo formats a time as "X ago"
func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)
	
	if duration < time.Minute {
		return fmt.Sprintf("%.0fs ago", duration.Seconds())
	} else if duration < time.Hour {
		return fmt.Sprintf("%.0fm ago", duration.Minutes())
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%.0fh ago", duration.Hours())
	} else {
		return fmt.Sprintf("%.0fd ago", duration.Hours()/24)
	}
}

// FormatQuickStats formats quick statistics for the metrics panel
func FormatQuickStats(service *app.Service) string {
	if service.Statistics == nil {
		return "\n [gray]No statistics available"
	}

	stats := service.Statistics
	var output strings.Builder

	output.WriteString("\n")
	
	// Always show uptime
	uptime := time.Since(stats.StartTime)
	output.WriteString(fmt.Sprintf(" [yellow]Uptime:[white]   %s\n", 
		components.FormatDuration(int64(uptime.Seconds()))))
	
	// Show request stats (even if zero)
	output.WriteString(fmt.Sprintf(" [yellow]Requests:[white] %d total, %.1f/s\n", 
		stats.TotalRequests, stats.CalculateRequestsPerSecond()))
	
	if stats.TotalErrors > 0 {
		errorRate := float64(stats.TotalErrors) / float64(stats.TotalRequests) * 100
		output.WriteString(fmt.Sprintf(" [yellow]Errors:[white]   %d total, %.1f%% rate\n", 
			stats.TotalErrors, errorRate))
	}
	
	// 5-minute averages (show even if zero)
	avgCPU := stats.GetAverageCPU(5 * time.Minute)
	_, avgMemPct := stats.GetAverageMemory(5 * time.Minute)
	
	output.WriteString(fmt.Sprintf(" [yellow]5m Avg:[white]   CPU %.1f%%, Mem %.1f%%\n", 
		avgCPU, avgMemPct))
	
	// Network rates (show even if zero)
	output.WriteString(fmt.Sprintf(" [yellow]Network:[white]  ↓%s/s ↑%s/s\n", 
		components.FormatBytes(stats.DataInRate),
		components.FormatBytes(stats.DataOutRate)))

	return output.String()
}
