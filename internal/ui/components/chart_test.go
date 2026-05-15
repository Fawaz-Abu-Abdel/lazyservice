package components

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestCreateBarChart(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		max      float64
		width    int
		expected int // expected number of filled characters
	}{
		{"50% of 100", 50, 100, 10, 5},
		{"75% of 100", 75, 100, 20, 15},
		{"100% of 100", 100, 100, 10, 10},
		{"0% of 100", 0, 100, 10, 0},
		{"Over 100%", 150, 100, 10, 10}, // Should cap at width
		{"Zero max", 50, 0, 10, 5},      // Should default to 100
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateBarChart(tt.value, tt.max, tt.width)
			
			resultWidth := utf8.RuneCountInString(result)
			if resultWidth != tt.width {
				t.Errorf("Expected width %d, got %d", tt.width, resultWidth)
			}
			
			filled := strings.Count(result, "█")
			if filled != tt.expected {
				t.Errorf("Expected %d filled chars, got %d", tt.expected, filled)
			}
		})
	}
}

func TestCreateSparkline(t *testing.T) {
	tests := []struct {
		name     string
		data     []float64
		width    int
		expected int // expected length
	}{
		{"Normal data", []float64{1, 2, 3, 4, 5}, 5, 5},
		{"Empty data", []float64{}, 10, 10},
		{"Single value", []float64{42}, 3, 3},
		{"Width larger than data", []float64{1, 2}, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateSparkline(tt.data, tt.width)
			
			resultWidth := utf8.RuneCountInString(result)
			if resultWidth != tt.expected {
				t.Errorf("Expected length %d, got %d", tt.expected, resultWidth)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{"Bytes", 512, "512 B"},
		{"Kilobytes", 1536, "1.5 KB"},
		{"Megabytes", 1572864, "1.5 MB"},
		{"Gigabytes", 1610612736, "1.5 GB"},
		{"Zero", 0, "0 B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int64
		expected string
	}{
		{"Minutes only", 300, "5m"},
		{"Hours and minutes", 3900, "1h 5m"},
		{"Days, hours, minutes", 90300, "1d 1h 5m"},
		{"Zero", 0, "0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.seconds)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestCreateMiniChart(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	
	result := CreateMiniChart(data, 3, 5)
	
	if len(result) != 3 {
		t.Errorf("Expected height 3, got %d", len(result))
	}
	
	for i, line := range result {
		resultWidth := utf8.RuneCountInString(line)
		if resultWidth != 5 {
			t.Errorf("Line %d: expected width 5, got %d", i, resultWidth)
		}
	}
	
	// Test empty data
	emptyResult := CreateMiniChart([]float64{}, 3, 5)
	if len(emptyResult) != 0 {
		t.Errorf("Expected empty result for empty data, got %d lines", len(emptyResult))
	}
}
