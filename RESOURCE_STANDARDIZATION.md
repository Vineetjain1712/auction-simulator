# Resource Standardization Guide

## Overview

This document explains how resource usage is standardized in the Auction Simulator for consistent benchmarking and comparison across different systems.

## CPU Standardization

### GOMAXPROCS Setting
The simulator uses `runtime.GOMAXPROCS()` to limit CPU core usage:
```go
runtime.GOMAXPROCS(4)  // Use exactly 4 cores
```

**Why 4 cores?**
- Provides enough parallelism for 40 concurrent auctions
- Commonly available on modern systems
- Allows fair comparison across different machines

### Measuring CPU Usage
- **Available CPUs**: `runtime.NumCPU()` - Total cores on the system
- **Used CPUs**: `runtime.GOMAXPROCS(0)` - Cores allocated to Go runtime
- **CPU Utilization**: `(Used / Available) * 100%`

## Memory Standardization

### Pre-Execution Cleanup
Before running benchmarks:
```go
runtime.GC()  // Force garbage collection
time.Sleep(100 * time.Millisecond)  // Allow GC to complete
```

### Memory Metrics Tracked
1. **Alloc**: Currently allocated heap memory
2. **TotalAlloc**: Cumulative bytes allocated
3. **Sys**: Total memory from OS
4. **NumGoroutine**: Active goroutines

### Measurement Points
- **Initial**: Before auctions start
- **Peak**: Maximum during execution
- **Final**: After all auctions complete
- **Average**: Mean across periodic samples

## Reproducible Results

### Factors Controlled
1. ✅ CPU cores (via GOMAXPROCS)
2. ✅ Initial memory state (via GC)
3. ✅ Configuration (auctions, bidders, timeout)
4. ✅ Random seeds (each run uses time-based seed)

### Factors NOT Controlled
1. ❌ OS background processes
2. ❌ Available RAM (system-dependent)
3. ❌ CPU frequency/throttling
4. ❌ Network conditions (not applicable here)

## Running Standardized Benchmarks

### Command
```bash
go run cmd/simulator/main.go
```

### Configuration
Default standardized settings (in `config/config.go`):
- Auctions: 40
- Bidders: 100
- Timeout: 10 seconds
- CPU Cores: 4
- Bid Probability: 30%

### Comparing Results

To compare across systems, focus on:
1. **Duration** (should be ~10 seconds due to timeout)
2. **Peak Memory** (normalized per goroutine)
3. **Bids/Second** (throughput metric)
4. **Memory/Goroutine** (efficiency metric)

## Example Comparison

| System | CPUs Used | Peak Memory | Bids/Sec | Memory/Goroutine |
|--------|-----------|-------------|----------|------------------|
| A      | 4         | 45.2 MB     | 125.3    | 0.011 MB         |
| B      | 4         | 48.7 MB     | 118.9    | 0.012 MB         |

System A is slightly more efficient in both speed and memory usage.

## Best Practices

1. **Close other applications** before benchmarking
2. **Run multiple times** and average results
3. **Use same configuration** across comparisons
4. **Note system specs** (CPU model, RAM)
5. **Check system load** before running

## Tuning Parameters

To test different configurations:
```go
cfg := config.DefaultConfig()
cfg.System.MaxCPUCores = 2  // Test with 2 cores
cfg.Auction.TotalAuctions = 80  // Double the load
```

## Resource Limits

Recommended system requirements:
- **Minimum**: 2 CPU cores, 2GB RAM
- **Recommended**: 4 CPU cores, 4GB RAM
- **Optimal**: 8 CPU cores, 8GB RAM

## Monitoring Tools

The simulator includes built-in monitoring:
- Periodic snapshots every 500ms
- Real-time goroutine tracking
- Memory allocation tracking
- Export to CSV for analysis