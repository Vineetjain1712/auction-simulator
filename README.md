# 🎯 Concurrent Auction Simulator

A high-performance concurrent auction system built in Go that simulates 40 simultaneous auctions with 100 bidders, featuring comprehensive resource monitoring and analytics.

## ✨ Features

- ✅ **40 Concurrent Auctions** running simultaneously
- ✅ **100 Simulated Bidders** with realistic behavior
- ✅ **20-Attribute Item System** for detailed auction objects
- ✅ **Timeout-Based Closure** with context management
- ✅ **Resource Monitoring** (CPU & Memory tracking)
- ✅ **Comprehensive Analytics** (stats, metrics, reports)
- ✅ **Multiple Export Formats** (JSON, CSV, TXT)
- ✅ **Race-Condition Free** verified with Go's race detector
- ✅ **~4,000+ Goroutines** managed efficiently

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- 4GB RAM minimum (8GB recommended)
- 4 CPU cores recommended

### Installation
```bash
# Clone the repository
git clone https://github.com/vineetjain1712/auction-simulator
cd auction-simulator

# Download dependencies
go mod tidy

# Run the simulator
go run cmd/simulator/main.go
```

### Expected Output

The simulator will:
1. Generate 40 unique items with 20 attributes each
2. Start 40 concurrent auctions
3. Activate 100 bidders across all auctions
4. Track resource usage in real-time
5. Export results to `./output` directory

## 📊 Performance Metrics

Typical results on a 4-core system:
- **Duration**: ~10 seconds (based on timeout)
- **Total Bids**: 1,000-1,500 bids
- **Peak Memory**: 15-25 MB
- **Peak Goroutines**: ~4,000
- **Throughput**: 100-150 bids/second

## 🏗️ Architecture
```
auction-simulator/
├── cmd/
│   ├── simulator/          # Main application
│   └── compare/            # Comparison tool
├── internal/
│   ├── auction/           # Auction logic & management
│   ├── bidder/            # Bidder simulation
│   ├── models/            # Data structures
│   ├── stats/             # Statistical analysis
│   ├── export/            # Export functionality
│   └── monitor/           # Resource monitoring
├── config/                # Configuration
├── test/                  # Integration tests
└── output/                # Generated results
```

## ⚙️ Configuration

Default settings in `config/config.go`:
```go
Auctions:      40
Bidders:       100
Timeout:       10 seconds
Bid Probability: 30%
CPU Cores:     4
```

## 🧪 Testing
```bash
# Run all tests
go test ./...  -v

# Run with race detector
go test ./... -race -v

# Run integration tests
go test ./test -v

# Run benchmarks
go test ./test -bench=. -benchmem
```

## 📈 Resource Standardization

The simulator standardizes resources for consistent benchmarking:

1. **CPU**: Uses `GOMAXPROCS(4)` to limit to 4 cores
2. **Memory**: Forces GC before execution
3. **Monitoring**: Samples every 500ms
4. **Metrics**: Exports detailed resource usage

See `RESOURCE_STANDARDIZATION.md` for details.

## 📁 Output Files

The simulator generates:
- `simulation_*.json` - Complete results in JSON
- `simulation_*.csv` - Auction data in CSV
- `resources_*.csv` - Resource metrics
- `summary_*.txt` - Human-readable summary

## 🔍 Key Components

### Auction Flow
1. Generate items with 20 attributes
2. Create auction instances
3. Start all auctions concurrently
4. Collect bids via channels
5. Determine winners on timeout
6. Aggregate results

### Bidder Behavior
- Each bidder independently decides whether to bid (30% probability)
- Random delay before bidding (100-2000ms)
- Bid amount varies based on base price (1.0x - 2.5x)
- Can participate in multiple auctions

### Concurrency Model
- **Main goroutine**: Orchestration
- **40 auction goroutines**: One per auction
- **4,000 bidder goroutines**: 100 bidders × 40 auctions
- **Synchronization**: WaitGroups, channels, mutexes

## 💡 Design Decisions

### Why Channels?
- Thread-safe communication
- Natural fit for concurrent bid collection
- Prevents race conditions

### Why Context?
- Clean timeout management
- Graceful cancellation
- Prevents goroutine leaks

### Why 4 Cores?
- Sufficient parallelism
- Commonly available
- Consistent benchmarking

## 📊 Example Results
```
Total Duration:     10.1 seconds
Success Rate:       95%
Total Bids:         1,247
Peak Memory:        18.9 MB
Peak Goroutines:    4,042
Bids/Second:        123.4
```

## 🛠️ Development

### Project Structure
- **cmd/**: Application entry points
- **internal/**: Private application code
- **config/**: Configuration management
- **test/**: Integration tests

### Code Style
- Clear, descriptive naming
- Comments on exported functions
- Small, focused functions
- Comprehensive error handling

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📝 License

MIT License - see LICENSE file

## 👤 Author

Vineet Jain 

## 🙏 Acknowledgments

- Built with Go's excellent concurrency primitives
- Inspired by real-world auction systems
- Designed for learning concurrent programming

## 📚 Learn More

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Memory Model](https://go.dev/ref/mem)

---

**Happy Auctioning!** 🎉