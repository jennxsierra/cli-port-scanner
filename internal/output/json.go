// Filename: json.go
// Description: Contains functions for writing scan results to a JSON file.
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jennxsierra/cli-port-scanner/internal/scanner"
)

// ScanOutput wraps overall scan metadata and results.
type ScanOutput struct {
	Timestamp  string           `json:"timestamp"`
	Targets    []string         `json:"targets"`
	TotalPorts int              `json:"total_ports"`
	Results    []scanner.Result `json:"results"`
}

// NewScanOutput creates a new ScanOutput instance.
func NewScanOutput(targets []string, totalPorts int, results []scanner.Result) ScanOutput {
	return ScanOutput{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Targets:    targets,
		TotalPorts: totalPorts,
		Results:    results,
	}
}

// WriteJSON marshals the ScanOutput struct to JSON and writes it to a file.
func WriteJSON(outputData ScanOutput) error {
	bytes, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Filename format: DDMMYY-HHMMSS-cli-pscan.json.
	fileName := time.Now().Format("020106-150405") + "-cli-pscan.json"
	if err := os.WriteFile(fileName, bytes, 0644); err != nil {
		return fmt.Errorf("error writing JSON file: %w", err)
	}

	fmt.Printf("Scan results saved to %s\n", fileName)
	return nil
}
