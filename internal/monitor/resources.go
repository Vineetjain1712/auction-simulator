package monitor

import (
	"fmt"
	"runtime"
	"time"
)

// ResourceSnapshot represents a point-in-time snapshot of resource usage
type ResourceSnapshot struct {
	Timestamp      time.Time
	MemoryAllocMB  float64 // Currently allocated memory in MB
	MemoryTotalMB  float64 // Total memory obtained from OS in MB
	MemorySysMB    float64 // Total memory from system in MB
	NumGoroutines  int     // Number of goroutines
	NumCPU         int     // Number of CPUs available
	GOMAXPROCS     int     // Number of CPUs being used
}

// ResourceMonitor tracks system resource usage
type ResourceMonitor struct {
	snapshots     []ResourceSnapshot
	startSnapshot ResourceSnapshot
	stopSnapshot  ResourceSnapshot
	interval      time.Duration
	stopChan      chan struct{}
}

// NewResourceMonitor creates a new resource monitor
func NewResourceMonitor(interval time.Duration) *ResourceMonitor {
	return &ResourceMonitor{
		snapshots: make([]ResourceSnapshot, 0),
		interval:  interval,
		stopChan:  make(chan struct{}),
	}
}

// Start begins monitoring resources at the specified interval
func (rm *ResourceMonitor) Start() {
	// Take initial snapshot
	rm.startSnapshot = rm.takeSnapshot()
	rm.snapshots = append(rm.snapshots, rm.startSnapshot)
	
	go func() {
		ticker := time.NewTicker(rm.interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				snapshot := rm.takeSnapshot()
				rm.snapshots = append(rm.snapshots, snapshot)
			case <-rm.stopChan:
				return
			}
		}
	}()
}

// Stop stops monitoring and takes a final snapshot
func (rm *ResourceMonitor) Stop() {
	close(rm.stopChan)
	rm.stopSnapshot = rm.takeSnapshot()
	rm.snapshots = append(rm.snapshots, rm.stopSnapshot)
}

// takeSnapshot captures current resource usage
func (rm *ResourceMonitor) takeSnapshot() ResourceSnapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return ResourceSnapshot{
		Timestamp:      time.Now(),
		MemoryAllocMB:  float64(m.Alloc) / 1024 / 1024,
		MemoryTotalMB:  float64(m.TotalAlloc) / 1024 / 1024,
		MemorySysMB:    float64(m.Sys) / 1024 / 1024,
		NumGoroutines:  runtime.NumGoroutine(),
		NumCPU:         runtime.NumCPU(),
		GOMAXPROCS:     runtime.GOMAXPROCS(0),
	}
}

// GetStats returns computed statistics from all snapshots
func (rm *ResourceMonitor) GetStats() ResourceStats {
	if len(rm.snapshots) == 0 {
		return ResourceStats{}
	}
	
	var stats ResourceStats
	
	// Calculate memory statistics
	stats.InitialMemoryMB = rm.startSnapshot.MemoryAllocMB
	stats.FinalMemoryMB = rm.stopSnapshot.MemoryAllocMB
	stats.PeakMemoryMB = rm.startSnapshot.MemoryAllocMB
	
	maxGoroutines := rm.startSnapshot.NumGoroutines
	totalMemory := 0.0
	
	for _, snapshot := range rm.snapshots {
		if snapshot.MemoryAllocMB > stats.PeakMemoryMB {
			stats.PeakMemoryMB = snapshot.MemoryAllocMB
		}
		if snapshot.NumGoroutines > maxGoroutines {
			maxGoroutines = snapshot.NumGoroutines
		}
		totalMemory += snapshot.MemoryAllocMB
	}
	
	stats.AverageMemoryMB = totalMemory / float64(len(rm.snapshots))
	stats.MemoryDeltaMB = stats.FinalMemoryMB - stats.InitialMemoryMB
	stats.PeakGoroutines = maxGoroutines
	stats.NumCPU = rm.startSnapshot.NumCPU
	stats.GOMAXPROCS = rm.startSnapshot.GOMAXPROCS
	
	return stats
}

// GetSnapshots returns all captured snapshots
func (rm *ResourceMonitor) GetSnapshots() []ResourceSnapshot {
	return rm.snapshots
}

// ResourceStats contains aggregated resource statistics
type ResourceStats struct {
	InitialMemoryMB float64 // Memory at start
	FinalMemoryMB   float64 // Memory at end
	PeakMemoryMB    float64 // Maximum memory used
	AverageMemoryMB float64 // Average memory across snapshots
	MemoryDeltaMB   float64 // Change in memory (final - initial)
	PeakGoroutines  int     // Maximum concurrent goroutines
	NumCPU          int     // Total CPUs available
	GOMAXPROCS      int     // CPUs being used
}

// FormatReport generates a formatted report of resource usage
func (rs ResourceStats) FormatReport() string {
	report := "\n💻 RESOURCE USAGE REPORT\n"
	report += "════════════════════════════════════════════════════════\n\n"
	
	report += "🧠 Memory Usage:\n"
	report += fmt.Sprintf("   ├─ Initial:     %.2f MB\n", rs.InitialMemoryMB)
	report += fmt.Sprintf("   ├─ Final:       %.2f MB\n", rs.FinalMemoryMB)
	report += fmt.Sprintf("   ├─ Peak:        %.2f MB\n", rs.PeakMemoryMB)
	report += fmt.Sprintf("   ├─ Average:     %.2f MB\n", rs.AverageMemoryMB)
	report += fmt.Sprintf("   └─ Delta:       %+.2f MB\n\n", rs.MemoryDeltaMB)
	
	report += "⚙️  CPU & Concurrency:\n"
	report += fmt.Sprintf("   ├─ Available CPUs:    %d\n", rs.NumCPU)
	report += fmt.Sprintf("   ├─ GOMAXPROCS:        %d\n", rs.GOMAXPROCS)
	report += fmt.Sprintf("   └─ Peak Goroutines:   %d\n\n", rs.PeakGoroutines)
	
	// Calculate efficiency
	cpuUtilization := float64(rs.GOMAXPROCS) / float64(rs.NumCPU) * 100
	report += "📊 Efficiency Metrics:\n"
	report += fmt.Sprintf("   ├─ CPU Utilization:   %.1f%%\n", cpuUtilization)
	report += fmt.Sprintf("   └─ Memory Efficiency: %.2f MB/goroutine (peak)\n\n", 
		rs.PeakMemoryMB/float64(rs.PeakGoroutines))
	
	return report
}

// StandardizeResources sets consistent resource limits for benchmarking
func StandardizeResources(maxCPUs int) {
	// Set maximum CPUs to use
	runtime.GOMAXPROCS(maxCPUs)
	
	// Force garbage collection to start clean
	runtime.GC()
	
	// Give GC a moment to complete
	time.Sleep(100 * time.Millisecond)
}

// GetCurrentResources returns current resource snapshot
func GetCurrentResources() ResourceSnapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return ResourceSnapshot{
		Timestamp:      time.Now(),
		MemoryAllocMB:  float64(m.Alloc) / 1024 / 1024,
		MemoryTotalMB:  float64(m.TotalAlloc) / 1024 / 1024,
		MemorySysMB:    float64(m.Sys) / 1024 / 1024,
		NumGoroutines:  runtime.NumGoroutine(),
		NumCPU:         runtime.NumCPU(),
		GOMAXPROCS:     runtime.GOMAXPROCS(0),
	}
}