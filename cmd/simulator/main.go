package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/vineetjain1712/auction-simulator/config"
	"github.com/vineetjain1712/auction-simulator/internal/auction"
	"github.com/vineetjain1712/auction-simulator/internal/bidder"
	"github.com/vineetjain1712/auction-simulator/internal/export"
	"github.com/vineetjain1712/auction-simulator/internal/models"
	"github.com/vineetjain1712/auction-simulator/internal/monitor"
	"github.com/vineetjain1712/auction-simulator/internal/stats"
)

func main() {
	printBanner()

	// Load configuration
	cfg := config.DefaultConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("❌ Invalid configuration: %v", err)
	}

	// Standardize resources for consistent measurements
	monitor.StandardizeResources(cfg.System.MaxCPUCores)

	printConfiguration(cfg)

	// Run the full simulation with monitoring
	result := runFullSimulation(cfg)

	// Analyze results
	analyzer := stats.NewAnalyzer()
	statistics := analyzer.Analyze(result)

	// Display results
	displayResults(result)

	// Display statistics
	fmt.Println(analyzer.FormatReport(statistics))

	// Display resource usage
	displayResourceUsage(result)

	// Export results
	exportResults(result, analyzer.FormatReport(statistics))

	// Final summary
	printFinalSummary(result, statistics)
}

// printBanner displays the application banner
func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║          🎯 CONCURRENT AUCTION SIMULATOR v1.0            ║
║                                                           ║
║              Built with Go • Phase 5 Complete            ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

// runFullSimulation orchestrates the entire auction simulation with monitoring
func runFullSimulation(cfg *config.Config) models.SimulationResult {
	fmt.Println("🎬 Starting Simulation")
	fmt.Println("════════════════════════════════════════════════════════")
	
	// Start resource monitoring
	resourceMonitor := monitor.NewResourceMonitor(500 * time.Millisecond)
	resourceMonitor.Start()
	
	ctx := context.Background()
	
	// Create manager and bidder pool
	manager := auction.NewManager(cfg)
	bidderPool := bidder.NewPool(&cfg.Bidder)
	
	// Pre-create all auctions
	items := manager.Generator.GenerateItems(cfg.Auction.TotalAuctions)
	for i, item := range items {
		auc := auction.NewAuction(i+1, item, cfg.Auction.AuctionTimeout)
		manager.Auctions = append(manager.Auctions, auc)
	}
	
	fmt.Printf("📦 Pre-generated %d auctions\n", len(manager.Auctions))
	
	var wg sync.WaitGroup
	
	// Record start time
	manager.StartTime = time.Now()
	fmt.Printf("⏱️  Start Time: %s\n\n", manager.StartTime.Format("15:04:05.000"))
	
	// Start all auctions
	fmt.Println("🔨 Starting all auctions...")
	for _, auc := range manager.Auctions {
		wg.Add(1)
		go func(auction *auction.Auction) {
			defer wg.Done()
			result := auction.Run(ctx)
			
			manager.Mu.Lock()
			manager.Results = append(manager.Results, result)
			manager.Mu.Unlock()
		}(auc)
	}
	
	// Small delay to ensure auctions are running
	time.Sleep(50 * time.Millisecond)
	
	// Activate bidders
	fmt.Println("👥 Activating bidders...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		bidderPool.ParticipateInAllAuctions(ctx, manager.Auctions)
	}()
	
	// Wait for completion
	fmt.Println("⏳ Waiting for completion...")
	wg.Wait()
	
	manager.EndTime = time.Now()
	fmt.Printf("\n⏱️  End Time: %s\n", manager.EndTime.Format("15:04:05.000"))
	
	// Stop monitoring ONCE 
	resourceMonitor.Stop()
	resourceStats := resourceMonitor.GetStats()
	
	fmt.Println("\n✅ Simulation Complete!")
	
	// Build result with resource metrics
	result := manager.AggregateResults()
	result.CPUCount = resourceStats.NumCPU
	result.CPUUsed = resourceStats.GOMAXPROCS
	result.InitialMemoryMB = resourceStats.InitialMemoryMB
	result.FinalMemoryMB = resourceStats.FinalMemoryMB
	result.PeakMemoryMB = resourceStats.PeakMemoryMB
	result.AverageMemoryMB = resourceStats.AverageMemoryMB
	result.PeakGoroutines = resourceStats.PeakGoroutines
	
	return result
}

// printConfiguration displays the simulation configuration
func printConfiguration(cfg *config.Config) {
	fmt.Printf("📊 Configuration\n")
	fmt.Printf("════════════════════════════════════════════════════════\n")
	fmt.Printf("  Concurrent Auctions:    %d\n", cfg.Auction.TotalAuctions)
	fmt.Printf("  Total Bidders:          %d\n", cfg.Bidder.TotalBidders)
	fmt.Printf("  Auction Timeout:        %v\n", cfg.Auction.AuctionTimeout)
	fmt.Printf("  Bid Probability:        %.1f%%\n", cfg.Bidder.BidProbability*100)
	fmt.Printf("  CPU Cores Available:    %d\n", runtime.NumCPU())
	fmt.Printf("  CPU Cores Used:         %d\n", cfg.System.MaxCPUCores)
	fmt.Printf("  Expected Goroutines:    ~%d\n",
		cfg.Auction.TotalAuctions+(cfg.Bidder.TotalBidders*cfg.Auction.TotalAuctions))
	fmt.Println()
}

// displayResults shows comprehensive simulation results
func displayResults(result models.SimulationResult) {
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("📊 SIMULATION RESULTS")
	fmt.Println(strings.Repeat("═", 60))

	// Timing
	fmt.Printf("\n⏱️  Timing:\n")
	fmt.Printf("   ├─ Start:      %s\n", result.StartTime.Format("15:04:05.000"))
	fmt.Printf("   ├─ End:        %s\n", result.EndTime.Format("15:04:05.000"))
	fmt.Printf("   └─ Duration:   %v\n", result.TotalDuration)

	// Auction Summary
	fmt.Printf("\n🔨 Auction Summary:\n")
	fmt.Printf("   ├─ Total:      %d\n", result.TotalAuctions)
	fmt.Printf("   ├─ Successful: %d (%.1f%%)\n",
		result.SuccessfulAuctions,
		float64(result.SuccessfulAuctions)/float64(result.TotalAuctions)*100)
	fmt.Printf("   └─ Failed:     %d\n", result.FailedAuctions)

	// Bidding Activity
	fmt.Printf("\n💰 Bidding Activity:\n")
	fmt.Printf("   ├─ Total Bids:       %d\n", result.TotalBids)
	fmt.Printf("   └─ Avg per Auction:  %.1f\n",
		float64(result.TotalBids)/float64(result.TotalAuctions))

	// Top auctions
	fmt.Printf("\n🏆 Top 5 Most Popular Auctions:\n")
	displayTopAuctions(result.AuctionResults, 5)

	// Winners
	fmt.Printf("\n🎉 Winners:\n")
	displayWinnersSummary(result.AuctionResults)
}

// displayResourceUsage shows resource utilization
func displayResourceUsage(result models.SimulationResult) {
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("💻 RESOURCE UTILIZATION")
	fmt.Println(strings.Repeat("═", 60))

	fmt.Printf("\n🧠 Memory:\n")
	fmt.Printf("   ├─ Initial:        %.2f MB\n", result.InitialMemoryMB)
	fmt.Printf("   ├─ Final:          %.2f MB\n", result.FinalMemoryMB)
	fmt.Printf("   ├─ Peak:           %.2f MB\n", result.PeakMemoryMB)
	fmt.Printf("   ├─ Average:        %.2f MB\n", result.AverageMemoryMB)
	fmt.Printf("   └─ Delta:          %+.2f MB\n", result.FinalMemoryMB-result.InitialMemoryMB)

	fmt.Printf("\n⚙️  CPU & Concurrency:\n")
	fmt.Printf("   ├─ CPUs Available:     %d\n", result.CPUCount)
	fmt.Printf("   ├─ CPUs Used:          %d (%.1f%%)\n",
		result.CPUUsed,
		float64(result.CPUUsed)/float64(result.CPUCount)*100)
	fmt.Printf("   └─ Peak Goroutines:    %d\n", result.PeakGoroutines)

	fmt.Printf("\n📊 Efficiency:\n")
	memPerGoroutine := result.PeakMemoryMB / float64(result.PeakGoroutines)
	fmt.Printf("   ├─ Memory/Goroutine:   %.3f MB\n", memPerGoroutine)

	bidsPerSecond := float64(result.TotalBids) / result.TotalDuration.Seconds()
	fmt.Printf("   ├─ Bids/Second:        %.1f\n", bidsPerSecond)

	auctionsPerSecond := float64(result.TotalAuctions) / result.TotalDuration.Seconds()
	fmt.Printf("   └─ Auctions/Second:    %.2f\n", auctionsPerSecond)
}

// displayTopAuctions shows the most popular auctions
func displayTopAuctions(results []models.AuctionResult, topN int) {
	// Sort by bid count
	sorted := make([]models.AuctionResult, len(results))
	copy(sorted, results)

	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].TotalBids < sorted[j+1].TotalBids {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	// Display top N
	for i := 0; i < topN && i < len(sorted); i++ {
		result := sorted[i]
		winnerInfo := "No winner"
		if result.WinningBid != nil {
			winnerInfo = fmt.Sprintf("Bidder #%d - $%.2f",
				result.WinningBid.BidderID, result.WinningBid.Amount)
		}

		fmt.Printf("   %d. Auction #%-3d: %3d bids → %s\n",
			i+1, result.AuctionID, result.TotalBids, winnerInfo)
	}
}

// displayWinnersSummary shows statistics about winners
func displayWinnersSummary(results []models.AuctionResult) {
	winnerMap := make(map[int]int)
	totalRevenue := 0.0

	for _, result := range results {
		if result.WinningBid != nil {
			winnerMap[result.WinningBid.BidderID]++
			totalRevenue += result.WinningBid.Amount
		}
	}

	fmt.Printf("   ├─ Unique Winners:  %d\n", len(winnerMap))
	fmt.Printf("   ├─ Total Revenue:   $%.2f\n", totalRevenue)

	if len(winnerMap) > 0 {
		avgWin := totalRevenue / float64(len(winnerMap))
		fmt.Printf("   └─ Avg Win Amount:  $%.2f\n", avgWin)

		// Find top winner
		maxWins := 0
		topBidder := 0
		for bidderID, wins := range winnerMap {
			if wins > maxWins {
				maxWins = wins
				topBidder = bidderID
			}
		}

		if maxWins > 1 {
			fmt.Printf("\n   🌟 Top Winner: Bidder #%d (%d auctions won)\n",
				topBidder, maxWins)
		}
	}
}

// exportResults exports simulation results to files
func exportResults(result models.SimulationResult, statsReport string) {
	fmt.Println("\n💾 Exporting Results")
	fmt.Println("════════════════════════════════════════════════════════")

	exporter := export.NewExporter("./output")

	// Export JSON
	if jsonFile, err := exporter.ExportToJSON(result); err != nil {
		fmt.Printf("   ✗ JSON export failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ JSON exported: %s\n", jsonFile)
	}

	// Export CSV
	if csvFile, err := exporter.ExportToCSV(result); err != nil {
		fmt.Printf("   ✗ CSV export failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ CSV exported: %s\n", csvFile)
	}

	// Export Summary
	if summaryFile, err := exporter.ExportSummary(result, statsReport); err != nil {
		fmt.Printf("   ✗ Summary export failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ Summary exported: %s\n", summaryFile)
	}

	// Export Resource Metrics
	if resourceFile, err := exporter.ExportResourceMetrics(result); err != nil {
		fmt.Printf("   ✗ Resource export failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ Resources exported: %s\n", resourceFile)
	}
}

// printFinalSummary displays final performance summary
func printFinalSummary(result models.SimulationResult, stats stats.Statistics) {
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("✨ FINAL SUMMARY")
	fmt.Println(strings.Repeat("═", 60))

	fmt.Printf("\n⚡ Performance:\n")
	fmt.Printf("   ├─ Total Time:           %v\n", result.TotalDuration)
	fmt.Printf("   ├─ Bids/Second:          %.1f\n", stats.BidsPerSecond)
	fmt.Printf("   ├─ Success Rate:         %.1f%%\n", stats.SuccessRate)
	fmt.Printf("   ├─ Peak Memory:          %.2f MB\n", result.PeakMemoryMB)
	fmt.Printf("   └─ Peak Goroutines:      %d\n", result.PeakGoroutines)

	fmt.Printf("\n✅ Simulation completed successfully!\n")
	fmt.Printf("📁 Results saved to ./output directory\n\n")
}
