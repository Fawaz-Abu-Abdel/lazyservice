package components

import (
	"testing"
)

func BenchmarkCreateBarChart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CreateBarChart(75.5, 100, 20)
	}
}

func BenchmarkCreateSparkline(b *testing.B) {
	data := make([]float64, 100)
	for i := range data {
		data[i] = float64(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateSparkline(data, 50)
	}
}

func BenchmarkFormatBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatBytes(1572864) // 1.5 MB
	}
}

func BenchmarkCreateMiniChart(b *testing.B) {
	data := make([]float64, 60)
	for i := range data {
		data[i] = float64(i % 10)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateMiniChart(data, 5, 20)
	}
}
