// Package config provides configuration types and defaults for the auction simulator.
package config

import (
	"fmt"
	"time"
)

// Config holds all simulation configuration
type Config struct {
	Auction AuctionConfig
	Bidder  BidderConfig
	System  SystemConfig
}

// AuctionConfig holds auction-specific settings
type AuctionConfig struct {
	TotalAuctions       int           // Number of concurrent auctions (40)
	AuctionTimeout      time.Duration // How long each auction runs
	MinimumBidIncrement float64       // Minimum bid increase
}

// BidderConfig holds bidder-specific settings
type BidderConfig struct {
	TotalBidders     int     // Number of bidders (100)
	BidProbability   float64 // Chance a bidder will bid (0.0 to 1.0)
	MinBidMultiplier float64 // Min bid = BasePrice * multiplier
	MaxBidMultiplier float64 // Max bid = BasePrice * multiplier
	BidDelayMinMs    int     // Min delay before bidding (ms)
	BidDelayMaxMs    int     // Max delay before bidding (ms)
}

// SystemConfig holds system resource settings
type SystemConfig struct {
	MaxCPUCores     int    // Maximum CPU cores to use
	EnableProfiling bool   // Enable CPU/memory profiling
	LogLevel        string // "debug", "info", "warn", "error"
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Auction: AuctionConfig{
			TotalAuctions:       40,
			AuctionTimeout:      10 * time.Second, // 10 seconds per auction
			MinimumBidIncrement: 1.0,
		},
		Bidder: BidderConfig{
			TotalBidders:     100,
			BidProbability:   0.3, // 30% chance to bid
			MinBidMultiplier: 1.0, // Bid at least base price
			MaxBidMultiplier: 2.5, // Bid up to 2.5x base price
			BidDelayMinMs:    100,
			BidDelayMaxMs:    2000,
		},
		System: SystemConfig{
			MaxCPUCores:     4, // Use 4 cores for consistency
			EnableProfiling: true,
			LogLevel:        "info",
		},
	}
}

// Validate checks if configuration is valid
func (c *Config) Validate() error {
	if c.Auction.TotalAuctions <= 0 {
		return fmt.Errorf("total auctions must be positive")
	}
	if c.Bidder.TotalBidders <= 0 {
		return fmt.Errorf("total bidders must be positive")
	}
	if c.Bidder.BidProbability < 0 || c.Bidder.BidProbability > 1 {
		return fmt.Errorf("bid probability must be between 0 and 1")
	}
	return nil
}
