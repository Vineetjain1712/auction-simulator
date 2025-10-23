package models

import "time"

// AuctionItem represents an item being auctioned with 20 attributes
type AuctionItem struct {
	ID          int     // Unique identifier
	Name        string  // Item name
	Category    string  // Category (Electronics, Art, etc.)
	Brand       string  // Brand name
	Condition   string  // New, Used, Refurbished
	Color       string  // Primary color
	Size        string  // Size (Small, Medium, Large, XL)
	Weight      float64 // Weight in kg
	Material    string  // Primary material
	YearMade    int     // Manufacturing year
	Origin      string  // Country of origin
	Rarity      string  // Common, Rare, Ultra-Rare
	BasePrice   float64 // Starting/reserve price
	Description string  // Item description
	Features    string  // Key features
	Warranty    int     // Warranty in months
	ShipWeight  float64 // Shipping weight
	Dimensions  string  // L x W x H in cm
	Certification string // Any certifications
	Rating      float64 // Quality rating (1-10)
}

// Bid represents a bid placed by a bidder
type Bid struct {
	BidderID  int       // Who placed the bid
	AuctionID int       // Which auction
	Amount    float64   // Bid amount
	Timestamp time.Time // When the bid was placed
}

// AuctionResult represents the outcome of an auction
type AuctionResult struct {
	AuctionID     int           // Auction identifier
	Item          AuctionItem   // The item that was auctioned
	WinningBid    *Bid          // Winning bid (nil if no bids)
	TotalBids     int           // Total number of bids received
	Duration      time.Duration // How long the auction ran
	StartTime     time.Time     // When auction started
	EndTime       time.Time     // When auction ended
	Status        string        // "completed", "no_bids", "timeout"
}

// BidderStats represents statistics for a bidder
type BidderStats struct {
	BidderID       int // Bidder identifier
	TotalBids      int // How many bids placed
	AuctionsWon    int // How many auctions won
	TotalSpent     float64 // Total amount spent
	AverageWinBid  float64 // Average winning bid amount
}

// SimulationResult represents the overall simulation results
type SimulationResult struct {
	TotalAuctions     int                  // Number of auctions run
	TotalDuration     time.Duration        // Total time from start to finish
	StartTime         time.Time            // First auction start time
	EndTime           time.Time            // Last auction end time
	AuctionResults    []AuctionResult      // Results of all auctions
	SuccessfulAuctions int                 // Auctions with at least one bid
	FailedAuctions    int                  // Auctions with no bids
	TotalBids         int                  // Total bids across all auctions
	CPUUsage          float64              // CPU usage percentage
	MemoryUsedMB      float64              // Memory used in MB
	PeakMemoryMB      float64              // Peak memory usage
}