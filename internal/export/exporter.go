package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vineetjain1712/auction-simulator/internal/models"
)

// Exporter handles exporting simulation results
type Exporter struct {
	outputDir string
}

// NewExporter creates a new exporter
func NewExporter(outputDir string) *Exporter {
	return &Exporter{
		outputDir: outputDir,
	}
}

// ExportToJSON exports simulation results to JSON file
func (e *Exporter) ExportToJSON(result models.SimulationResult) (string, error) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(e.outputDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(e.outputDir, fmt.Sprintf("simulation_%s.json", timestamp))

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0o644); err != nil {
		return "", fmt.Errorf("failed to write JSON file: %w", err)
	}

	return filename, nil
}

// ExportToCSV exports auction results to CSV file
func (e *Exporter) ExportToCSV(result models.SimulationResult) (string, error) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(e.outputDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(e.outputDir, fmt.Sprintf("simulation_%s.csv", timestamp))

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"AuctionID",
		"ItemName",
		"ItemCategory",
		"BasePrice",
		"Status",
		"TotalBids",
		"WinnerBidderID",
		"WinningAmount",
		"Duration_ms",
	}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write rows
	for _, auctionResult := range result.AuctionResults {
		row := []string{
			fmt.Sprintf("%d", auctionResult.AuctionID),
			auctionResult.Item.Name,
			auctionResult.Item.Category,
			fmt.Sprintf("%.2f", auctionResult.Item.BasePrice),
			auctionResult.Status,
			fmt.Sprintf("%d", auctionResult.TotalBids),
		}

		// Add winner info
		if auctionResult.WinningBid != nil {
			row = append(row,
				fmt.Sprintf("%d", auctionResult.WinningBid.BidderID),
				fmt.Sprintf("%.2f", auctionResult.WinningBid.Amount),
			)
		} else {
			row = append(row, "N/A", "N/A")
		}

		// Add duration
		row = append(row, fmt.Sprintf("%d", auctionResult.Duration.Milliseconds()))

		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return filename, nil
}

// ExportSummary exports a summary text file
func (e *Exporter) ExportSummary(result models.SimulationResult, statsReport string) (string, error) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(e.outputDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(e.outputDir, fmt.Sprintf("summary_%s.txt", timestamp))

	// Create summary content
	summary := fmt.Sprintf("AUCTION SIMULATION SUMMARY\n")
	summary += fmt.Sprintf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	summary += fmt.Sprintf("═══════════════════════════════════════════════════════\n\n")

	summary += fmt.Sprintf("Timing:\n")
	summary += fmt.Sprintf("  Start: %s\n", result.StartTime.Format("15:04:05.000"))
	summary += fmt.Sprintf("  End:   %s\n", result.EndTime.Format("15:04:05.000"))
	summary += fmt.Sprintf("  Total: %v\n\n", result.TotalDuration)

	summary += fmt.Sprintf("Overview:\n")
	summary += fmt.Sprintf("  Total Auctions: %d\n", result.TotalAuctions)
	summary += fmt.Sprintf("  Successful: %d\n", result.SuccessfulAuctions)
	summary += fmt.Sprintf("  Failed: %d\n", result.FailedAuctions)
	summary += fmt.Sprintf("  Total Bids: %d\n\n", result.TotalBids)

	summary += statsReport

	// Write to file
	if err := os.WriteFile(filename, []byte(summary), 0o644); err != nil {
		return "", fmt.Errorf("failed to write summary file: %w", err)
	}

	return filename, nil
}

// ExportResourceMetrics exports resource usage to a separate CSV
func (e *Exporter) ExportResourceMetrics(result models.SimulationResult) (string, error) {
	if err := os.MkdirAll(e.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}
	
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(e.outputDir, fmt.Sprintf("resources_%s.csv", timestamp))
	
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create resource CSV: %w", err)
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// Write header
	header := []string{
		"Metric",
		"Value",
		"Unit",
	}
	if err := writer.Write(header); err != nil {
		return "", err
	}
	
	// Write rows
	rows := [][]string{
		{"CPU_Available", fmt.Sprintf("%d", result.CPUCount), "cores"},
		{"CPU_Used", fmt.Sprintf("%d", result.CPUUsed), "cores"},
		{"Initial_Memory", fmt.Sprintf("%.2f", result.InitialMemoryMB), "MB"},
		{"Final_Memory", fmt.Sprintf("%.2f", result.FinalMemoryMB), "MB"},
		{"Peak_Memory", fmt.Sprintf("%.2f", result.PeakMemoryMB), "MB"},
		{"Average_Memory", fmt.Sprintf("%.2f", result.AverageMemoryMB), "MB"},
		{"Peak_Goroutines", fmt.Sprintf("%d", result.PeakGoroutines), "count"},
		{"Duration", fmt.Sprintf("%.3f", result.TotalDuration.Seconds()), "seconds"},
		{"Bids_Per_Second", fmt.Sprintf("%.1f", float64(result.TotalBids)/result.TotalDuration.Seconds()), "bids/s"},
	}
	
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}
	
	return filename, nil
}