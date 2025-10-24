package bidder

import (
	"context"
	"fmt"
	"sync"

	"github.com/vineetjain1712/auction-simulator/config"
	"github.com/vineetjain1712/auction-simulator/internal/auction"
)

// Pool manages a collection of bidders
type Pool struct {
	bidders []*Bidder
	config  *config.BidderConfig
}

// NewPool creates a pool of bidders
func NewPool(cfg *config.BidderConfig) *Pool {
	bidders := make([]*Bidder, cfg.TotalBidders)

	for i := 0; i < cfg.TotalBidders; i++ {
		bidders[i] = NewBidder(i+1, cfg)
	}

	return &Pool{
		bidders: bidders,
		config:  cfg,
	}
}

// ParticipateInAllAuctions makes all bidders participate in all auctions
// Each bidder can bid on multiple auctions
func (p *Pool) ParticipateInAllAuctions(ctx context.Context, auctions []*auction.Auction) {
	fmt.Printf("👥 Activating %d bidders for %d auctions\n",
		len(p.bidders), len(auctions))

	var wg sync.WaitGroup

	// For each bidder
	for _, bidder := range p.bidders {
		// For each auction
		for _, auc := range auctions {
			wg.Add(1)

			// Launch goroutine for this bidder-auction pair
			go func(b *Bidder, auction *auction.Auction) {
				defer wg.Done()

				// Create a context with auction timeout
				auctionCtx, cancel := context.WithTimeout(ctx, auction.Timeout)
				defer cancel()

				// Bidder participates in this auction
				b.ParticipateInAuction(
					auctionCtx,
					auction.ID,
					auction.Item,
					auction.GetBidChannel(),
				)
			}(bidder, auc)
		}
	}

	// Wait for all bidder-auction interactions to complete
	wg.Wait()

	fmt.Println("✅ All bidders have finished participating")
}

// GetBidders returns all bidders in the pool
func (p *Pool) GetBidders() []*Bidder {
	return p.bidders
}

// GetBidderCount returns the number of bidders
func (p *Pool) GetBidderCount() int {
	return len(p.bidders)
}
