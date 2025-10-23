package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/vineetjain1712/auction-simulator/config"
)

func main() {
	fmt.Println("🎯 Auction Simulator Starting...")
	fmt.Println("=" + string(make([]byte, 50)) + "=")

	// Load configuration
	cfg := config.DefaultConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Set CPU cores for consistent resource usage
	runtime.GOMAXPROCS(cfg.System.MaxCPUCores)

	fmt.Printf("📊 Configuration:\n")
	fmt.Printf("   - Auctions: %d (concurrent)\n", cfg.Auction.TotalAuctions)
	fmt.Printf("   - Bidders: %d\n", cfg.Bidder.TotalBidders)
	fmt.Printf("   - Timeout: %v per auction\n", cfg.Auction.AuctionTimeout)
	fmt.Printf("   - CPU Cores: %d\n", cfg.System.MaxCPUCores)
	fmt.Println()

	// TODO: Phase 2 - Create and run simulation
	fmt.Println("✅ Setup complete! Ready for Phase 2...")
}
