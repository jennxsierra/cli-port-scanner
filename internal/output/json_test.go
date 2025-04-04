package output

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jennxsierra/cli-port-scanner/internal/scanner"
)

func TestNewScanOutput(t *testing.T) {
	targets := []string{"example.com"}
	totalPorts := 80
	results := []scanner.Result{
		{Target: "example.com", TotalPorts: 80, OpenCount: 2, Duration: "0.789s"},
	}

	outputData := NewScanOutput(targets, totalPorts, results)
	if outputData.TotalPorts != totalPorts {
		t.Errorf("Expected total_ports %d, got %d", totalPorts, outputData.TotalPorts)
	}
	if !strings.Contains(outputData.Timestamp, "T") {
		t.Errorf("Expected a proper timestamp, got %q", outputData.Timestamp)
	}
}

func TestWriteJSON(t *testing.T) {
	// Create a dummy ScanOutput.
	outputData := NewScanOutput(
		[]string{"test.com"},
		2,
		[]scanner.Result{
			{Target: "test.com", TotalPorts: 2, OpenCount: 1, Duration: "0.123s"},
		},
	)

	// Write JSON. This will create a file in the "scan-results" folder.
	err := WriteJSON(outputData)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	// Check that at least one file exists in the scan-results directory.
	files, err := os.ReadDir("scan-results")
	if err != nil {
		t.Fatalf("Error reading scan-results directory: %v", err)
	}
	if len(files) == 0 {
		t.Errorf("Expected JSON file to be created, but directory is empty")
	}

	// Clean up created JSON files.
	for _, file := range files {
		os.Remove(filepath.Join("scan-results", file.Name()))
	}
}
