package stats

import (
	"fmt"
	"math"
	"sort"

	"github.com/vineetjain1712/auction-simulator/internal/models"
)

// Statistics holds calculated statistics from simulation
type Statistics struct {
	// Bid Statistics
	TotalBids   int
	AverageBids float64
	MinBids     int
	MaxBids     int
	MedianBids  float64
	StdDevBids  float64

	// Amount Statistics
	TotalRevenue     float64
	AverageWinAmount float64
	MinWinAmount     float64
	MaxWinAmount     float64
	MedianWinAmount  float64

	// Bidder Statistics
	UniqueBidders        int
	UniqueWinners        int
	MostActiveBidder     int
	MostSuccessfulBidder int

	// Performance Statistics
	BidsPerSecond     float64
	AuctionsPerSecond float64

	// Success Metrics
	SuccessRate     float64
	AuctionsFailed  int
	AuctionsSuccess int
}

// Analyzer analyzes simulation results
type Analyzer struct{}

// NewAnalyzer creates a new statistics analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// Analyze performs comprehensive analysis on simulation results
func (a *Analyzer) Analyze(result models.SimulationResult) Statistics {
	stats := Statistics{
		TotalBids:       result.TotalBids,
		AuctionsSuccess: result.SuccessfulAuctions,
		AuctionsFailed:  result.FailedAuctions,
	}

	// Calculate bid statistics
	a.analyzeBidCounts(result.AuctionResults, &stats)

	// Calculate amount statistics
	a.analyzeWinningAmounts(result.AuctionResults, &stats)

	// Calculate bidder statistics
	a.analyzeBidders(result.AuctionResults, &stats)

	// Calculate performance metrics
	a.analyzePerformance(result, &stats)

	// Calculate success rate
	if result.TotalAuctions > 0 {
		stats.SuccessRate = float64(stats.AuctionsSuccess) / float64(result.TotalAuctions) * 100
	}

	return stats
}

// analyzeBidCounts calculates statistics about bid counts
func (a *Analyzer) analyzeBidCounts(results []models.AuctionResult, stats *Statistics) {
	if len(results) == 0 {
		return
	}

	bidCounts := make([]int, len(results))
	sum := 0
	stats.MinBids = results[0].TotalBids
	stats.MaxBids = results[0].TotalBids

	for i, result := range results {
		bidCounts[i] = result.TotalBids
		sum += result.TotalBids

		if result.TotalBids < stats.MinBids {
			stats.MinBids = result.TotalBids
		}
		if result.TotalBids > stats.MaxBids {
			stats.MaxBids = result.TotalBids
		}
	}

	// Average
	stats.AverageBids = float64(sum) / float64(len(results))

	// Median
	sort.Ints(bidCounts)
	mid := len(bidCounts) / 2
	if len(bidCounts)%2 == 0 {
		stats.MedianBids = float64(bidCounts[mid-1]+bidCounts[mid]) / 2.0
	} else {
		stats.MedianBids = float64(bidCounts[mid])
	}

	// Standard Deviation
	variance := 0.0
	for _, count := range bidCounts {
		diff := float64(count) - stats.AverageBids
		variance += diff * diff
	}
	variance /= float64(len(bidCounts))
	stats.StdDevBids = math.Sqrt(variance)
}

// analyzeWinningAmounts calculates statistics about winning bid amounts
func (a *Analyzer) analyzeWinningAmounts(results []models.AuctionResult, stats *Statistics) {
	amounts := make([]float64, 0)

	for _, result := range results {
		if result.WinningBid != nil {
			amounts = append(amounts, result.WinningBid.Amount)
			stats.TotalRevenue += result.WinningBid.Amount
		}
	}

	if len(amounts) == 0 {
		return
	}

	// Min/Max
	stats.MinWinAmount = amounts[0]
	stats.MaxWinAmount = amounts[0]

	for _, amount := range amounts {
		if amount < stats.MinWinAmount {
			stats.MinWinAmount = amount
		}
		if amount > stats.MaxWinAmount {
			stats.MaxWinAmount = amount
		}
	}

	// Average
	stats.AverageWinAmount = stats.TotalRevenue / float64(len(amounts))

	// Median
	sort.Float64s(amounts)
	mid := len(amounts) / 2
	if len(amounts)%2 == 0 {
		stats.MedianWinAmount = (amounts[mid-1] + amounts[mid]) / 2.0
	} else {
		stats.MedianWinAmount = amounts[mid]
	}
}

// analyzeBidders calculates statistics about bidder activity
func (a *Analyzer) analyzeBidders(results []models.AuctionResult, stats *Statistics) {
	// bidderBids := make(map[int]int) // bidderID -> total bids
	bidderWins := make(map[int]int) // bidderID -> total wins

	for _, result := range results {
		// Count wins
		if result.WinningBid != nil {
			bidderWins[result.WinningBid.BidderID]++
		}

		// Count bids (we'd need to track this in auction, for now approximate)
		// This is a simplified version
	}

	stats.UniqueWinners = len(bidderWins)

	// Find most successful bidder
	maxWins := 0
	for bidderID, wins := range bidderWins {
		if wins > maxWins {
			maxWins = wins
			stats.MostSuccessfulBidder = bidderID
		}
	}
}

// analyzePerformance calculates performance metrics
func (a *Analyzer) analyzePerformance(result models.SimulationResult, stats *Statistics) {
	durationSeconds := result.TotalDuration.Seconds()

	if durationSeconds > 0 {
		stats.BidsPerSecond = float64(result.TotalBids) / durationSeconds
		stats.AuctionsPerSecond = float64(result.TotalAuctions) / durationSeconds
	}
}

// FormatReport generates a formatted text report
func (a *Analyzer) FormatReport(stats Statistics) string {
	report := "\n📈 DETAILED STATISTICS\n"
	report += "════════════════════════════════════════════════════════\n\n"

	// Bid Statistics
	report += "💰 Bid Statistics:\n"
	report += fmt.Sprintf("   ├─ Total Bids: %d\n", stats.TotalBids)
	report += fmt.Sprintf("   ├─ Average per Auction: %.1f\n", stats.AverageBids)
	report += fmt.Sprintf("   ├─ Median: %.1f\n", stats.MedianBids)
	report += fmt.Sprintf("   ├─ Min/Max: %d / %d\n", stats.MinBids, stats.MaxBids)
	report += fmt.Sprintf("   └─ Std Deviation: %.2f\n\n", stats.StdDevBids)

	// Amount Statistics
	if stats.TotalRevenue > 0 {
		report += "💵 Revenue Statistics:\n"
		report += fmt.Sprintf("   ├─ Total Revenue: $%.2f\n", stats.TotalRevenue)
		report += fmt.Sprintf("   ├─ Average Win: $%.2f\n", stats.AverageWinAmount)
		report += fmt.Sprintf("   ├─ Median Win: $%.2f\n", stats.MedianWinAmount)
		report += fmt.Sprintf("   └─ Min/Max: $%.2f / $%.2f\n\n", stats.MinWinAmount, stats.MaxWinAmount)
	}

	// Bidder Statistics
	report += "👥 Bidder Statistics:\n"
	report += fmt.Sprintf("   ├─ Unique Winners: %d\n", stats.UniqueWinners)
	if stats.MostSuccessfulBidder > 0 {
		report += fmt.Sprintf("   └─ Top Bidder: #%d\n\n", stats.MostSuccessfulBidder)
	} else {
		report += "   └─ No winners\n\n"
	}

	// Performance Metrics
	report += "⚡ Performance Metrics:\n"
	report += fmt.Sprintf("   ├─ Bids/Second: %.1f\n", stats.BidsPerSecond)
	report += fmt.Sprintf("   ├─ Auctions/Second: %.2f\n", stats.AuctionsPerSecond)
	report += fmt.Sprintf("   └─ Success Rate: %.1f%%\n\n", stats.SuccessRate)

	return report
}
