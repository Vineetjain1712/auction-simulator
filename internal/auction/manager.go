package auction

import (
	"sync"
	"time"

	"github.com/vineetjain1712/auction-simulator/config"
	"github.com/vineetjain1712/auction-simulator/internal/models"
)

// Manager orchestrates multiple concurrent auctions
type Manager struct {
	config    *config.Config
	Generator *ItemGenerator

	// Track all auctions - EXPORTED so main can access
	Auctions []*Auction

	// Overall timing - EXPORTED
	StartTime time.Time
	EndTime   time.Time

	// Results collection
	Results []models.AuctionResult
	Mu      sync.Mutex // EXPORTED
}

// NewManager creates a new auction manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:    cfg,
		Generator: NewItemGenerator(),
		Auctions:  make([]*Auction, 0, cfg.Auction.TotalAuctions),
		Results:   make([]models.AuctionResult, 0, cfg.Auction.TotalAuctions),
	}
}

// AggregateResults compiles all auction results into a simulation result
func (m *Manager) AggregateResults() models.SimulationResult {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	totalBids := 0
	successfulAuctions := 0
	failedAuctions := 0

	for _, result := range m.Results {
		totalBids += result.TotalBids

		if result.Status == "completed" && result.WinningBid != nil {
			successfulAuctions++
		} else {
			failedAuctions++
		}
	}

	return models.SimulationResult{
		TotalAuctions:      m.config.Auction.TotalAuctions,
		TotalDuration:      m.EndTime.Sub(m.StartTime),
		StartTime:          m.StartTime,
		EndTime:            m.EndTime,
		AuctionResults:     m.Results,
		SuccessfulAuctions: successfulAuctions,
		FailedAuctions:     failedAuctions,
		TotalBids:          totalBids,
	}
}

// GetAuctions returns all auction instances (for backwards compatibility)
func (m *Manager) GetAuctions() []*Auction {
	return m.Auctions
}
