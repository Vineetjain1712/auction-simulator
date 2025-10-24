package bidder

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/vineetjain1712/auction-simulator/config"
	"github.com/vineetjain1712/auction-simulator/internal/models"
)

// Bidder represents a simulated bidder
type Bidder struct {
	ID     int
	config *config.BidderConfig
	rand   *rand.Rand
	mu     sync.Mutex // Protects rand for thread-safety
}

// NewBidder creates a new bidder with given ID
func NewBidder(id int, cfg *config.BidderConfig) *Bidder {
	// Each bidder gets its own random source for thread safety
	source := rand.NewSource(time.Now().UnixNano() + int64(id))
	return &Bidder{
		ID:     id,
		config: cfg,
		rand:   rand.New(source),
	}
}

// DecideIfBid determines if this bidder wants to bid on an item
// Returns true if bidder decides to bid, false otherwise
func (b *Bidder) DecideIfBid(item models.AuctionItem) bool {
	// Random decision based on bid probability
	// E.g., if BidProbability is 0.3, there's 30% chance to bid
	b.mu.Lock()
	decision := b.rand.Float64() < b.config.BidProbability
	b.mu.Unlock()
	return decision
}

// CalculateBidAmount determines how much to bid
// Based on the item's base price and configured multipliers
func (b *Bidder) CalculateBidAmount(item models.AuctionItem) float64 {
	// Random multiplier between MinBidMultiplier and MaxBidMultiplier
	b.mu.Lock()
	multiplier := b.config.MinBidMultiplier +
		b.rand.Float64()*(b.config.MaxBidMultiplier-b.config.MinBidMultiplier)
	b.mu.Unlock()

	return item.BasePrice * multiplier
}

// SimulateBidDelay simulates the time it takes for a bidder to decide and bid
// Returns the delay duration
func (b *Bidder) SimulateBidDelay() time.Duration {
	// Random delay between min and max
	b.mu.Lock()
	delayMs := b.config.BidDelayMinMs +
		b.rand.Intn(b.config.BidDelayMaxMs-b.config.BidDelayMinMs+1)
	b.mu.Unlock()

	return time.Duration(delayMs) * time.Millisecond
}

// ParticipateInAuction simulates a bidder participating in an auction
// It receives auction details, decides whether to bid, and sends bid if interested
func (b *Bidder) ParticipateInAuction(
	ctx context.Context,
	auctionID int,
	item models.AuctionItem,
	bidChannel chan<- models.Bid,
) {
	// First, decide if this bidder is interested
	if !b.DecideIfBid(item) {
		// Not interested, don't bid
		return
	}

	// Simulate thinking time
	delay := b.SimulateBidDelay()

	// Create a timer for the delay
	timer := time.NewTimer(delay)
	defer timer.Stop()

	// Wait for either delay or context cancellation
	select {
	case <-timer.C:
		// Delay complete - check if auction is still active
		select {
		case <-ctx.Done():
			// Auction closed during our delay
			return
		default:
			// Auction still active, proceed with bid
		}

		// Calculate bid amount
		amount := b.CalculateBidAmount(item)

		// Create the bid
		bid := models.Bid{
			BidderID:  b.ID,
			AuctionID: auctionID,
			Amount:    amount,
			Timestamp: time.Now(),
		}

		// Try to send the bid, but respect context
		select {
		case bidChannel <- bid:
			// Bid sent successfully
		case <-ctx.Done():
			// Auction closed while we were trying to send
			return
		}

	case <-ctx.Done():
		// Auction closed during our thinking time
		return
	}
}
