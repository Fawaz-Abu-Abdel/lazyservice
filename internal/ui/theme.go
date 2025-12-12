package ui

import (
	"strings"
	
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Theme defines the color scheme and styling for the UI
type Theme struct {
	// Background colors
	BackgroundColor     tcell.Color
	PrimaryBackground   tcell.Color
	SecondaryBackground tcell.Color
	
	// Border colors
	BorderColor         tcell.Color
	FocusedBorderColor  tcell.Color
	
	// Text colors
	PrimaryText         tcell.Color
	SecondaryText       tcell.Color
	AccentText          tcell.Color
	
	// Status colors
	SuccessColor        tcell.Color
	WarningColor        tcell.Color
	ErrorColor          tcell.Color
	InfoColor           tcell.Color
	
	// Special colors
	HighlightColor      tcell.Color
	SelectionColor      tcell.Color
}

// ModernDarkTheme returns a beautiful modern dark theme
func ModernDarkTheme() *Theme {
	return &Theme{
		BackgroundColor:     tcell.ColorBlack,
		PrimaryBackground:   tcell.NewRGBColor(24, 24, 37),   // #181825
		SecondaryBackground: tcell.NewRGBColor(30, 30, 46),   // #1e1e2e
		
		BorderColor:         tcell.NewRGBColor(88, 91, 112),  // #585b70
		FocusedBorderColor:  tcell.NewRGBColor(137, 180, 250), // #89b4fa
		
		PrimaryText:         tcell.NewRGBColor(205, 214, 244), // #cdd6f4
		SecondaryText:       tcell.NewRGBColor(166, 173, 200), // #a6adc8
		AccentText:          tcell.NewRGBColor(137, 180, 250), // #89b4fa
		
		SuccessColor:        tcell.NewRGBColor(166, 227, 161), // #a6e3a1
		WarningColor:        tcell.NewRGBColor(249, 226, 175), // #f9e2af
		ErrorColor:          tcell.NewRGBColor(243, 139, 168), // #f38ba8
		InfoColor:           tcell.NewRGBColor(137, 180, 250), // #89b4fa
		
		HighlightColor:      tcell.NewRGBColor(203, 166, 247), // #cba6f7
		SelectionColor:      tcell.NewRGBColor(49, 50, 68),    // #313244
	}
}

// CyberpunkTheme returns a vibrant cyberpunk-inspired theme
func CyberpunkTheme() *Theme {
	return &Theme{
		BackgroundColor:     tcell.ColorBlack,
		PrimaryBackground:   tcell.NewRGBColor(13, 13, 23),    // #0d0d17
		SecondaryBackground: tcell.NewRGBColor(20, 20, 35),    // #141423
		
		BorderColor:         tcell.NewRGBColor(0, 255, 255),   // #00ffff
		FocusedBorderColor:  tcell.NewRGBColor(255, 0, 255),   // #ff00ff
		
		PrimaryText:         tcell.NewRGBColor(0, 255, 255),   // #00ffff
		SecondaryText:       tcell.NewRGBColor(128, 255, 255), // #80ffff
		AccentText:          tcell.NewRGBColor(255, 0, 255),   // #ff00ff
		
		SuccessColor:        tcell.NewRGBColor(0, 255, 0),     // #00ff00
		WarningColor:        tcell.NewRGBColor(255, 255, 0),   // #ffff00
		ErrorColor:          tcell.NewRGBColor(255, 0, 0),     // #ff0000
		InfoColor:           tcell.NewRGBColor(0, 255, 255),   // #00ffff
		
		HighlightColor:      tcell.NewRGBColor(255, 0, 255),   // #ff00ff
		SelectionColor:      tcell.NewRGBColor(40, 0, 40),     // #280028
	}
}

// ApplyTheme applies the theme to tview components
func ApplyTheme(theme *Theme) {
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    theme.BackgroundColor,
		ContrastBackgroundColor:     theme.PrimaryBackground,
		MoreContrastBackgroundColor: theme.SecondaryBackground,
		BorderColor:                 theme.BorderColor,
		TitleColor:                  theme.AccentText,
		GraphicsColor:               theme.BorderColor,
		PrimaryTextColor:            theme.PrimaryText,
		SecondaryTextColor:          theme.SecondaryText,
		TertiaryTextColor:           theme.SecondaryText,
		InverseTextColor:            theme.BackgroundColor,
		ContrastSecondaryTextColor:  theme.SecondaryText,
	}
}

// StyledBox creates a beautifully styled box with the theme
func StyledBox(title string, theme *Theme, focused bool) *tview.Box {
	box := tview.NewBox()
	box.SetTitle(" " + title + " ")
	box.SetBorder(true)
	
	if focused {
		box.SetBorderColor(theme.FocusedBorderColor)
		box.SetTitleColor(theme.AccentText)
	} else {
		box.SetBorderColor(theme.BorderColor)
		box.SetTitleColor(theme.PrimaryText)
	}
	
	return box
}

// CreateGradientBorder creates a gradient-like border effect using Unicode characters
func CreateGradientBorder(width int, color string) string {
	if width <= 0 {
		return ""
	}
	
	chars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	result := ""
	
	for i := 0; i < width; i++ {
		charIndex := (i * len(chars)) / width
		if charIndex >= len(chars) {
			charIndex = len(chars) - 1
		}
		result += "[" + color + "]" + chars[charIndex] + "[white]"
	}
	
	return result
}

// CreateSparkline creates a beautiful sparkline chart
func CreateSparkline(data []float64, width int, color string) string {
	if len(data) == 0 || width <= 0 {
		return ""
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
	
	// Avoid division by zero
	if max == min {
		return "[" + color + "]" + strings.Repeat("▄", width) + "[white]"
	}
	
	// Create sparkline
	chars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	result := "[" + color + "]"
	
	// Sample data to fit width
	step := float64(len(data)) / float64(width)
	for i := 0; i < width; i++ {
		dataIndex := int(float64(i) * step)
		if dataIndex >= len(data) {
			dataIndex = len(data) - 1
		}
		
		// Normalize to 0-7 range
		normalized := (data[dataIndex] - min) / (max - min) * 7
		charIndex := int(normalized)
		if charIndex >= len(chars) {
			charIndex = len(chars) - 1
		}
		
		result += chars[charIndex]
	}
	
	result += "[white]"
	return result
}

// StatusIcon returns a colored status icon
func StatusIcon(status string, theme *Theme) string {
	switch status {
	case "running":
		return "[green]●[white]"
	case "stopped":
		return "[red]●[white]"
	case "failed":
		return "[red]✗[white]"
	case "warning":
		return "[yellow]⚠[white]"
	default:
		return "[gray]○[white]"
	}
}

// ProgressBar creates a beautiful progress bar
func ProgressBar(percentage float64, width int, color string) string {
	if width <= 0 {
		return ""
	}
	
	filled := int((percentage / 100.0) * float64(width))
	if filled > width {
		filled = width
	}
	
	bar := "[" + color + "]"
	bar += strings.Repeat("█", filled)
	bar += "[gray]"
	bar += strings.Repeat("░", width-filled)
	bar += "[white]"
	
	return bar
}

// AnimatedSpinner returns different spinner frames for animations
func AnimatedSpinner(frame int) string {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return "[cyan]" + spinners[frame%len(spinners)] + "[white]"
}

// BoxDrawing provides Unicode box drawing characters
var BoxDrawing = struct {
	Horizontal     string
	Vertical       string
	TopLeft        string
	TopRight       string
	BottomLeft     string
	BottomRight    string
	Cross          string
	VerticalRight  string
	VerticalLeft   string
	HorizontalDown string
	HorizontalUp   string
}{
	Horizontal:     "─",
	Vertical:       "│",
	TopLeft:        "┌",
	TopRight:       "┐",
	BottomLeft:     "└",
	BottomRight:    "┘",
	Cross:          "┼",
	VerticalRight:  "├",
	VerticalLeft:   "┤",
	HorizontalDown: "┬",
	HorizontalUp:   "┴",
}
