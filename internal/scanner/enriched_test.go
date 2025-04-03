package scanner

import (
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// startTestServerWithBanner starts a TCP server on 127.0.0.1:port that writes the given banner.
func startTestServerWithBanner(t *testing.T, port int, banner string) net.Listener {
	listener, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			conn.Write([]byte(banner))
			conn.Close()
		}
	}()
	return listener
}

func TestEnrichedScanner_Run(t *testing.T) {
	// Start a test server on a specific port with a banner.
	testPort := 8081
	expectedBanner := "Test Banner"
	listener := startTestServerWithBanner(t, testPort, expectedBanner)
	defer listener.Close()

	// Create an EnrichedScanner instance scanning the test port.
	es := EnrichedScanner{
		Targets:    []string{"127.0.0.1"},
		Ports:      []int{testPort},
		Workers:    1,
		Timeout:    2 * time.Second,
		JsonOutput: false,
	}

	summaries, results := es.Run()

	// Validate summaries for the target.
	if len(summaries) != 1 {
		t.Errorf("expected one summary, got %d", len(summaries))
	}
	sum := summaries[0]
	if sum.Target != "127.0.0.1" {
		t.Errorf("expected target '127.0.0.1', got '%s'", sum.Target)
	}
	if sum.TotalPorts != 1 {
		t.Errorf("expected TotalPorts to be 1, got %d", sum.TotalPorts)
	}
	if len(sum.OpenPorts) != 1 {
		t.Errorf("expected one open port, got %d", len(sum.OpenPorts))
	}
	if sum.OpenPorts[0] != testPort {
		t.Errorf("expected open port %d, got %d", testPort, sum.OpenPorts[0])
	}

	// Validate enriched results.
	if len(results) != 1 {
		t.Errorf("expected one result, got %d", len(results))
	}
	result := results[0]
	if result.Target != "127.0.0.1" {
		t.Errorf("expected result.Target '127.0.0.1', got '%s'", result.Target)
	}
	if result.Port != testPort {
		t.Errorf("expected result.Port %d, got %d", testPort, result.Port)
	}
	if result.Status != "open" {
		t.Errorf("expected result.Status 'open', got '%s'", result.Status)
	}
	// Trim any extra whitespace from the banner.
	if !strings.HasPrefix(strings.TrimSpace(result.Banner), expectedBanner) {
		t.Errorf("expected banner to start with '%s', got '%s'", expectedBanner, result.Banner)
	}
}

func TestSaveJSON(t *testing.T) {
	// Prepare a dummy enriched result.
	results := []EnrichedResult{
		{
			Target: "dummy",
			Port:   22,
			Status: "open",
			Banner: "SSH Banner",
		},
	}

	filename, err := SaveJSON(results)
	if err != nil {
		t.Fatalf("failed to save JSON: %v", err)
	}
	if filename == "" {
		t.Error("expected a non-empty filename")
	}
	// Optionally, verify the file exists and clean up.
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", filename)
	} else {
		// Clean-up generated file.
		os.Remove(filename)
	}
}
