# Auction Simulator

A concurrent auction system that simulates 40 simultaneous auctions with 100 bidders.

## Features

- ✅ 40 concurrent auctions
- ✅ 100 simulated bidders
- ✅ 20-attribute item system
- ✅ Timeout-based auction closure
- ✅ Resource monitoring (CPU & Memory)
- ✅ Accurate timing measurement

## Prerequisites

- Go 1.21 or higher
- 4GB RAM minimum

## Installation
```bash
# Clone the repository
git clone https://github.com/vineetjain1712/auction-simulator
cd auction-simulator

# Download dependencies
go mod tidy

# Run the simulator
go run cmd/simulator/main.go
```

## Project Structure
```
auction-simulator/
├── cmd/simulator/          # Main application entry point
├── internal/
│   ├── auction/           # Auction logic
│   ├── bidder/            # Bidder simulation
│   ├── models/            # Data models
│   └── monitor/           # Resource monitoring
├── config/                # Configuration
└── README.md
```

## Configuration

Default configuration:
- 40 concurrent auctions
- 100 bidders
- 10-second timeout per auction
- 4 CPU cores

## Development Status

- [x] Phase 1: Project Setup ✅
- [ ] Phase 2: Core Logic
- [ ] Phase 3: Concurrency
- [ ] Phase 4: Scaling & Timing
- [ ] Phase 5: Resource Monitoring
- [ ] Phase 6: Polish & Documentation

## Author

Vineet Jain

## License

MIT