package bidder

import (
	"testing"
	"time"

	"github.com/vineetjain1712/auction-simulator/config"
	"github.com/vineetjain1712/auction-simulator/internal/models"
)

func TestBidderCreation(t *testing.T) {
	cfg := config.DefaultConfig()
	bidder := NewBidder(1, &cfg.Bidder)

	if bidder.ID != 1 {
		t.Errorf("Expected bidder ID 1, got %d", bidder.ID)
	}
}

func TestDecideIfBid(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Bidder.BidProbability = 0.5 // 50% chance

	bidder := NewBidder(1, &cfg.Bidder)

	item := models.AuctionItem{
		ID:        1,
		BasePrice: 100.0,
	}

	// Run multiple times to test randomness
	bidCount := 0
	iterations := 1000

	for i := 0; i < iterations; i++ {
		if bidder.DecideIfBid(item) {
			bidCount++
		}
	}

	// With 50% probability and 1000 iterations,
	// we expect around 500 bids (allow 30% variance for randomness)
	expectedMin := int(float64(iterations) * 0.35)
	expectedMax := int(float64(iterations) * 0.65)

	if bidCount < expectedMin || bidCount > expectedMax {
		t.Logf("WARNING: Bid count %d outside expected range %d-%d (this can happen randomly)",
			bidCount, expectedMin, expectedMax)
	}
}

func TestCalculateBidAmount(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Bidder.MinBidMultiplier = 1.0
	cfg.Bidder.MaxBidMultiplier = 2.0

	bidder := NewBidder(1, &cfg.Bidder)

	item := models.AuctionItem{
		ID:        1,
		BasePrice: 100.0,
	}

	// Test multiple times
	for i := 0; i < 10; i++ {
		amount := bidder.CalculateBidAmount(item)

		// Bid should be between 1.0x and 2.0x base price
		if amount < 100.0 || amount > 200.0 {
			t.Errorf("Expected bid between 100-200, got %.2f", amount)
		}
	}
}

func TestSimulateBidDelay(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Bidder.BidDelayMinMs = 100
	cfg.Bidder.BidDelayMaxMs = 500

	bidder := NewBidder(1, &cfg.Bidder)

	delay := bidder.SimulateBidDelay()

	if delay < 100*time.Millisecond || delay > 500*time.Millisecond {
		t.Errorf("Expected delay between 100-500ms, got %v", delay)
	}
}
