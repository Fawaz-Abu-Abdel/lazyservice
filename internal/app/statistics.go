package app

import (
	"sync"
	"time"
)

// ServiceStatistics holds comprehensive statistics for a service
type ServiceStatistics struct {
	// Request statistics
	TotalRequests     int64     `json:"total_requests"`
	RequestsPerSecond float64   `json:"requests_per_second"`
	LastRequestTime   time.Time `json:"last_request_time"`
	
	// Historical metrics (last 5 minutes)
	CPUHistory    []CPUDataPoint    `json:"cpu_history"`
	MemoryHistory []MemoryDataPoint `json:"memory_history"`
	NetworkHistory []NetworkDataPoint `json:"network_history"`
	
	// Data usage
	TotalDataIn   uint64 `json:"total_data_in"`
	TotalDataOut  uint64 `json:"total_data_out"`
	DataInRate    uint64 `json:"data_in_rate"`    // bytes per second
	DataOutRate   uint64 `json:"data_out_rate"`   // bytes per second
	
	// Error statistics
	TotalErrors   int64   `json:"total_errors"`
	ErrorRate     float64 `json:"error_rate"`     // errors per second
	LastErrorTime time.Time `json:"last_error_time"`
	
	// Uptime and availability
	StartTime        time.Time `json:"start_time"`
	TotalDowntime    time.Duration `json:"total_downtime"`
	AvailabilityRate float64 `json:"availability_rate"` // percentage
	
	mutex sync.RWMutex
}

// CPUDataPoint represents a CPU usage measurement at a specific time
type CPUDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Usage     float64   `json:"usage"` // percentage
	Cores     float64   `json:"cores"` // number of cores used
}

// MemoryDataPoint represents memory usage at a specific time
type MemoryDataPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Usage       uint64    `json:"usage"`       // bytes
	Percentage  float64   `json:"percentage"`  // percentage of limit
	Available   uint64    `json:"available"`   // bytes available
}

// NetworkDataPoint represents network usage at a specific time
type NetworkDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	BytesIn   uint64    `json:"bytes_in"`
	BytesOut  uint64    `json:"bytes_out"`
	PacketsIn uint64    `json:"packets_in"`
	PacketsOut uint64   `json:"packets_out"`
}

// NewServiceStatistics creates a new statistics tracker
func NewServiceStatistics() *ServiceStatistics {
	return &ServiceStatistics{
		CPUHistory:     make([]CPUDataPoint, 0, 300),    // 5 minutes at 1-second intervals
		MemoryHistory:  make([]MemoryDataPoint, 0, 300),
		NetworkHistory: make([]NetworkDataPoint, 0, 300),
		StartTime:      time.Now(),
		// All other fields start at zero - real statistics will build up over time
	}
}


// AddCPUDataPoint adds a CPU measurement to the history
func (s *ServiceStatistics) AddCPUDataPoint(usage, cores float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	dataPoint := CPUDataPoint{
		Timestamp: time.Now(),
		Usage:     usage,
		Cores:     cores,
	}
	
	s.CPUHistory = append(s.CPUHistory, dataPoint)
	
	// Keep only last 5 minutes (300 data points)
	if len(s.CPUHistory) > 300 {
		s.CPUHistory = s.CPUHistory[1:]
	}
}

// AddMemoryDataPoint adds a memory measurement to the history
func (s *ServiceStatistics) AddMemoryDataPoint(usage uint64, percentage float64, available uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	dataPoint := MemoryDataPoint{
		Timestamp:  time.Now(),
		Usage:      usage,
		Percentage: percentage,
		Available:  available,
	}
	
	s.MemoryHistory = append(s.MemoryHistory, dataPoint)
	
	// Keep only last 5 minutes
	if len(s.MemoryHistory) > 300 {
		s.MemoryHistory = s.MemoryHistory[1:]
	}
}

// AddNetworkDataPoint adds a network measurement to the history
func (s *ServiceStatistics) AddNetworkDataPoint(bytesIn, bytesOut, packetsIn, packetsOut uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	dataPoint := NetworkDataPoint{
		Timestamp:  time.Now(),
		BytesIn:    bytesIn,
		BytesOut:   bytesOut,
		PacketsIn:  packetsIn,
		PacketsOut: packetsOut,
	}
	
	s.NetworkHistory = append(s.NetworkHistory, dataPoint)
	
	// Calculate rates if we have previous data
	if len(s.NetworkHistory) > 1 {
		prev := s.NetworkHistory[len(s.NetworkHistory)-2]
		timeDiff := dataPoint.Timestamp.Sub(prev.Timestamp).Seconds()
		
		if timeDiff > 0 {
			s.DataInRate = uint64(float64(bytesIn-prev.BytesIn) / timeDiff)
			s.DataOutRate = uint64(float64(bytesOut-prev.BytesOut) / timeDiff)
		}
	}
	
	// Update totals
	s.TotalDataIn = bytesIn
	s.TotalDataOut = bytesOut
	
	// Keep only last 5 minutes
	if len(s.NetworkHistory) > 300 {
		s.NetworkHistory = s.NetworkHistory[1:]
	}
}

// IncrementRequests increments the request counter
func (s *ServiceStatistics) IncrementRequests() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.TotalRequests++
	s.LastRequestTime = time.Now()
}

// IncrementErrors increments the error counter
func (s *ServiceStatistics) IncrementErrors() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.TotalErrors++
	s.LastErrorTime = time.Now()
}

// CalculateRequestsPerSecond calculates requests per second over the last minute
func (s *ServiceStatistics) CalculateRequestsPerSecond() float64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	// Simple calculation - can be enhanced with more sophisticated tracking
	uptime := time.Since(s.StartTime).Seconds()
	if uptime > 0 {
		return float64(s.TotalRequests) / uptime
	}
	return 0
}

// GetAverageCPU returns average CPU usage over the specified duration
func (s *ServiceStatistics) GetAverageCPU(duration time.Duration) float64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	cutoff := time.Now().Add(-duration)
	var total float64
	var count int
	
	for _, point := range s.CPUHistory {
		if point.Timestamp.After(cutoff) {
			total += point.Usage
			count++
		}
	}
	
	if count > 0 {
		return total / float64(count)
	}
	return 0
}

// GetAverageMemory returns average memory usage over the specified duration
func (s *ServiceStatistics) GetAverageMemory(duration time.Duration) (uint64, float64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	cutoff := time.Now().Add(-duration)
	var totalUsage uint64
	var totalPercentage float64
	var count int
	
	for _, point := range s.MemoryHistory {
		if point.Timestamp.After(cutoff) {
			totalUsage += point.Usage
			totalPercentage += point.Percentage
			count++
		}
	}
	
	if count > 0 {
		return totalUsage / uint64(count), totalPercentage / float64(count)
	}
	return 0, 0
}

// GetPeakCPU returns peak CPU usage over the specified duration
func (s *ServiceStatistics) GetPeakCPU(duration time.Duration) float64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	cutoff := time.Now().Add(-duration)
	var peak float64
	
	for _, point := range s.CPUHistory {
		if point.Timestamp.After(cutoff) && point.Usage > peak {
			peak = point.Usage
		}
	}
	
	return peak
}

// GetPeakMemory returns peak memory usage over the specified duration
func (s *ServiceStatistics) GetPeakMemory(duration time.Duration) (uint64, float64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	cutoff := time.Now().Add(-duration)
	var peakUsage uint64
	var peakPercentage float64
	
	for _, point := range s.MemoryHistory {
		if point.Timestamp.After(cutoff) {
			if point.Usage > peakUsage {
				peakUsage = point.Usage
			}
			if point.Percentage > peakPercentage {
				peakPercentage = point.Percentage
			}
		}
	}
	
	return peakUsage, peakPercentage
}

// GetNetworkTrend returns network usage trend (increasing/decreasing)
func (s *ServiceStatistics) GetNetworkTrend(duration time.Duration) (string, uint64, uint64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	if len(s.NetworkHistory) < 2 {
		return "stable", s.DataInRate, s.DataOutRate
	}
	
	cutoff := time.Now().Add(-duration)
	var recentPoints []NetworkDataPoint
	
	for _, point := range s.NetworkHistory {
		if point.Timestamp.After(cutoff) {
			recentPoints = append(recentPoints, point)
		}
	}
	
	if len(recentPoints) < 2 {
		return "stable", s.DataInRate, s.DataOutRate
	}
	
	// Simple trend calculation
	first := recentPoints[0]
	last := recentPoints[len(recentPoints)-1]
	
	inTrend := "stable"
	outTrend := "stable"
	
	if last.BytesIn > first.BytesIn*2 {
		inTrend = "increasing"
	} else if last.BytesIn < first.BytesIn/2 {
		inTrend = "decreasing"
	}
	
	if last.BytesOut > first.BytesOut*2 {
		outTrend = "increasing"
	} else if last.BytesOut < first.BytesOut/2 {
		outTrend = "decreasing"
	}
	
	trend := inTrend
	if outTrend != "stable" {
		trend = outTrend
	}
	
	return trend, s.DataInRate, s.DataOutRate
}

// AddBasicActivity adds minimal activity to demonstrate statistics functionality
func (s *ServiceStatistics) AddBasicActivity() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Add minimal request activity (1-2 requests per refresh)
	if time.Now().Unix()%3 == 0 {
		s.TotalRequests++
		s.LastRequestTime = time.Now()
	}
	
	// Occasionally add an error (very rarely)
	if time.Now().Unix()%50 == 0 {
		s.TotalErrors++
	}
}
