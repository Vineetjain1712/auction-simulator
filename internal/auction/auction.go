// Package auction provides an Auction implementation for running auctions,
// collecting and storing bids, and determining winners for auction items.
// This package exposes an Auction type that can be run with a context and timeout,
// accepts bids via a channel, and returns an AuctionResult summarizing the outcome.
package auction

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/vineetjain1712/auction-simulator/internal/models"
)

// Auction represents a single auction instance
type Auction struct {
	ID      int
	Item    models.AuctionItem
	Timeout time.Duration

	// Channel to receive bids
	bidChannel chan models.Bid

	// Store all received bids
	bids []models.Bid
	mu   sync.Mutex // Protects bids slice

	// Timing
	startTime time.Time
	endTime   time.Time
}

// NewAuction creates a new auction instance
func NewAuction(id int, item models.AuctionItem, timeout time.Duration) *Auction {
	return &Auction{
		ID:         id,
		Item:       item,
		Timeout:    timeout,
		bidChannel: make(chan models.Bid, 100), // Buffered channel for bids
		bids:       make([]models.Bid, 0),
	}
}

// GetBidChannel returns the channel where bidders send their bids
func (a *Auction) GetBidChannel() chan<- models.Bid {
	return a.bidChannel
}

// Run starts the auction and runs it until timeout
// Returns the auction result
func (a *Auction) Run(ctx context.Context) models.AuctionResult {
	a.startTime = time.Now()

	// Only log every 10th auction to reduce noise
	if a.ID%10 == 0 || a.ID == 1 {
		fmt.Printf("🔨 Auction #%d started: %s (Base: $%.2f)\n",
			a.ID, a.Item.Name, a.Item.BasePrice)
	}

	// Create a context with timeout for this auction
	auctionCtx, cancel := context.WithTimeout(ctx, a.Timeout)
	defer cancel()

	// Collect bids until timeout
	a.collectBids(auctionCtx)

	a.endTime = time.Now()

	// Determine winner
	result := a.determineWinner()

	// Only log every 10th auction
	if a.ID%10 == 0 || a.ID == 1 {
		fmt.Printf("✅ Auction #%d ended: %d bids received\n", a.ID, result.TotalBids)
	}

	return result
}

// collectBids listens for incoming bids until auction closes
func (a *Auction) collectBids(ctx context.Context) {
	for {
		select {
		case bid, ok := <-a.bidChannel:
			if !ok {
				// Channel closed externally
				return
			}

			// Received a bid
			a.mu.Lock()
			a.bids = append(a.bids, bid)
			a.mu.Unlock()

		case <-ctx.Done():
			// Timeout reached, auction is closing
			// DON'T close the channel - just stop listening
			// Drain any buffered bids
			for {
				select {
				case bid, ok := <-a.bidChannel:
					if !ok {
						return
					}
					a.mu.Lock()
					a.bids = append(a.bids, bid)
					a.mu.Unlock()
				default:
					// No more buffered bids
					return
				}
			}
		}
	}
}

// determineWinner analyzes bids and determines the auction winner
func (a *Auction) determineWinner() models.AuctionResult {
	a.mu.Lock()
	defer a.mu.Unlock()

	result := models.AuctionResult{
		AuctionID: a.ID,
		Item:      a.Item,
		TotalBids: len(a.bids),
		StartTime: a.startTime,
		EndTime:   a.endTime,
		Duration:  a.endTime.Sub(a.startTime),
	}

	// Check if we have any bids
	if len(a.bids) == 0 {
		result.Status = "no_bids"
		result.WinningBid = nil
		return result
	}

	// Sort bids by amount (descending) to find highest bid
	sortedBids := make([]models.Bid, len(a.bids))
	copy(sortedBids, a.bids)

	sort.Slice(sortedBids, func(i, j int) bool {
		// If amounts are equal, earlier bid wins
		if sortedBids[i].Amount == sortedBids[j].Amount {
			return sortedBids[i].Timestamp.Before(sortedBids[j].Timestamp)
		}
		return sortedBids[i].Amount > sortedBids[j].Amount
	})

	// Winner is the highest bid
	winningBid := sortedBids[0]
	result.WinningBid = &winningBid
	result.Status = "completed"

	// Only log winners for interesting auctions
	// (removed logging here to reduce noise)

	return result
}

// GetAllBids returns all bids received (for testing/analysis)
func (a *Auction) GetAllBids() []models.Bid {
	a.mu.Lock()
	defer a.mu.Unlock()

	bidsCopy := make([]models.Bid, len(a.bids))
	copy(bidsCopy, a.bids)
	return bidsCopy
}
