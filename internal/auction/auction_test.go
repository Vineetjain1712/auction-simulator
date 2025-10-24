package auction

import (
	"context"
	"testing"
	"time"

	"github.com/vineetjain1712/auction-simulator/internal/models"
)

func TestItemGenerator(t *testing.T) {
	generator := NewItemGenerator()

	// Test single item generation
	item := generator.GenerateItem(1)

	// Verify all 20 attributes are set
	if item.ID != 1 {
		t.Errorf("Expected ID 1, got %d", item.ID)
	}

	if item.Name == "" {
		t.Error("Item name should not be empty")
	}

	if item.BasePrice <= 0 {
		t.Error("Base price should be positive")
	}

	if item.Rating < 0 || item.Rating > 10 {
		t.Errorf("Rating should be between 0-10, got %.2f", item.Rating)
	}

	// Test batch generation
	items := generator.GenerateItems(5)
	if len(items) != 5 {
		t.Errorf("Expected 5 items, got %d", len(items))
	}
}

func TestAuctionWithNoBids(t *testing.T) {
	generator := NewItemGenerator()
	item := generator.GenerateItem(1)

	auction := NewAuction(1, item, 100*time.Millisecond)

	ctx := context.Background()
	result := auction.Run(ctx)

	if result.Status != "no_bids" {
		t.Errorf("Expected status 'no_bids', got '%s'", result.Status)
	}

	if result.TotalBids != 0 {
		t.Errorf("Expected 0 bids, got %d", result.TotalBids)
	}

	if result.WinningBid != nil {
		t.Error("Expected no winning bid")
	}
}

func TestAuctionWithBids(t *testing.T) {
	generator := NewItemGenerator()
	item := generator.GenerateItem(1)

	auction := NewAuction(1, item, 200*time.Millisecond)

	// Start auction in background
	ctx := context.Background()
	done := make(chan models.AuctionResult)

	go func() {
		result := auction.Run(ctx)
		done <- result
	}()

	// Send some test bids
	bidChannel := auction.GetBidChannel()

	bid1 := models.Bid{BidderID: 1, AuctionID: 1, Amount: 100.0, Timestamp: time.Now()}
	bid2 := models.Bid{BidderID: 2, AuctionID: 1, Amount: 150.0, Timestamp: time.Now()}
	bid3 := models.Bid{BidderID: 3, AuctionID: 1, Amount: 120.0, Timestamp: time.Now()}

	bidChannel <- bid1
	bidChannel <- bid2
	bidChannel <- bid3

	// Wait for result
	result := <-done

	if result.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", result.Status)
	}

	if result.TotalBids != 3 {
		t.Errorf("Expected 3 bids, got %d", result.TotalBids)
	}

	if result.WinningBid == nil {
		t.Fatal("Expected a winning bid")
	}

	// Highest bid should win
	if result.WinningBid.Amount != 150.0 {
		t.Errorf("Expected winning bid of 150.0, got %.2f", result.WinningBid.Amount)
	}

	if result.WinningBid.BidderID != 2 {
		t.Errorf("Expected bidder 2 to win, got bidder %d", result.WinningBid.BidderID)
	}
}
