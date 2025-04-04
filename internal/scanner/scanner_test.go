package scanner

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestResultToJSON(t *testing.T) {
	result := Result{
		Target: "example.com",
		OpenPorts: []PortResult{
			{Port: 80, Banner: "BannerData"},
		},
		TotalPorts: 2,
		OpenCount:  1,
		Duration:   "0.456s",
	}
	jsonStr, err := result.ToJSON()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Quick check that the JSON string contains expected values.
	if !strings.Contains(jsonStr, "example.com") || !strings.Contains(jsonStr, "BannerData") {
		t.Errorf("JSON output does not contain expected data: %s", jsonStr)
	}

	// Verify JSON is valid.
	var parsed Result
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}
}

func TestToJSONAlternate(t *testing.T) {
	// Create a Result with no open ports.
	result := Result{
		Target:     "localhost",
		OpenPorts:  []PortResult{},
		TotalPorts: 100,
		OpenCount:  0,
		Duration:   "1.000s",
	}
	jsonStr, err := result.ToJSON()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !strings.Contains(jsonStr, `"target": "localhost"`) {
		t.Errorf("Output JSON missing target: %s", jsonStr)
	}
}
