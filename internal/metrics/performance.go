package metrics

import (
	"runtime"
	"sync"
	"time"
)

// PerformanceMetrics tracks application performance
type PerformanceMetrics struct {
	mutex           sync.RWMutex
	StartTime       time.Time
	LastUpdateTime  time.Time
	ServiceCount    int
	RefreshCount    int64
	ErrorCount      int64
	MemoryUsage     uint64
	GoroutineCount  int
	RefreshDuration time.Duration
}

var globalMetrics *PerformanceMetrics
var once sync.Once

// GetGlobalMetrics returns the global performance metrics instance
func GetGlobalMetrics() *PerformanceMetrics {
	once.Do(func() {
		globalMetrics = &PerformanceMetrics{
			StartTime: time.Now(),
		}
	})
	return globalMetrics
}

// UpdateServiceCount updates the service count
func (p *PerformanceMetrics) UpdateServiceCount(count int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ServiceCount = count
	p.LastUpdateTime = time.Now()
}

// IncrementRefreshCount increments the refresh counter
func (p *PerformanceMetrics) IncrementRefreshCount() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.RefreshCount++
}

// IncrementErrorCount increments the error counter
func (p *PerformanceMetrics) IncrementErrorCount() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ErrorCount++
}

// UpdateRefreshDuration updates the last refresh duration
func (p *PerformanceMetrics) UpdateRefreshDuration(duration time.Duration) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.RefreshDuration = duration
}

// UpdateSystemMetrics updates system-level metrics
func (p *PerformanceMetrics) UpdateSystemMetrics() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	p.MemoryUsage = m.Alloc
	p.GoroutineCount = runtime.NumGoroutine()
}

// GetSnapshot returns a snapshot of current metrics
func (p *PerformanceMetrics) GetSnapshot() PerformanceSnapshot {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	return PerformanceSnapshot{
		Uptime:          time.Since(p.StartTime),
		ServiceCount:    p.ServiceCount,
		RefreshCount:    p.RefreshCount,
		ErrorCount:      p.ErrorCount,
		MemoryUsage:     p.MemoryUsage,
		GoroutineCount:  p.GoroutineCount,
		RefreshDuration: p.RefreshDuration,
		LastUpdate:      p.LastUpdateTime,
	}
}

// PerformanceSnapshot represents a point-in-time snapshot of metrics
type PerformanceSnapshot struct {
	Uptime          time.Duration
	ServiceCount    int
	RefreshCount    int64
	ErrorCount      int64
	MemoryUsage     uint64
	GoroutineCount  int
	RefreshDuration time.Duration
	LastUpdate      time.Time
}
