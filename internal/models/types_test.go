package models

import (
	"testing"
	"time"
)

// TestAuctionItemCreation tests creating an auction item
func TestAuctionItemCreation(t *testing.T) {
	item := AuctionItem{
		ID:        1,
		Name:      "Vintage Camera",
		Category:  "Electronics",
		Brand:     "Canon",
		Condition: "Used",
		BasePrice: 100.0,
	}

	if item.ID != 1 {
		t.Errorf("Expected ID 1, got %d", item.ID)
	}

	if item.BasePrice != 100.0 {
		t.Errorf("Expected BasePrice 100.0, got %f", item.BasePrice)
	}
}

// TestBidCreation tests creating a bid
func TestBidCreation(t *testing.T) {
	bid := Bid{
		BidderID:  42,
		AuctionID: 1,
		Amount:    150.0,
		Timestamp: time.Now(),
	}

	if bid.Amount <= 0 {
		t.Error("Bid amount should be positive")
	}

	if bid.BidderID != 42 {
		t.Errorf("Expected BidderID 42, got %d", bid.BidderID)
	}
}

// TestAuctionResult tests creating an auction result
func TestAuctionResult(t *testing.T) {
	result := AuctionResult{
		AuctionID: 1,
		TotalBids: 5,
		Status:    "completed",
	}

	if result.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", result.Status)
	}
}
