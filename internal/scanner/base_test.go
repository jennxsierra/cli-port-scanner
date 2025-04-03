package scanner

import (
	"net"
	"slices"
	"strconv"
	"testing"
	"time"
)

// startTestServer starts a simple TCP server to simulate an open port.
func startTestServer(t *testing.T, port int) {
	l, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	// Run the server in a goroutine.
	go func() {
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				// Listener closed or error occurred.
				return
			}
			conn.Close()
		}
	}()
}

func TestFormatSummary(t *testing.T) {
	summary := Summary{
		Target:     "example.com",
		TotalPorts: 10,
		OpenPorts:  []int{22, 80},
		Duration:   250 * time.Millisecond,
	}
	formatted := FormatSummary(summary)
	expected := `[example.com]
Total Ports Scanned: 10
Open Ports Count: 2
Open Ports: [22 80]
Time Taken: 250ms
`
	if formatted != expected {
		t.Errorf("Expected summary:\n%s\nGot:\n%s", expected, formatted)
	}
}

func TestScannerRun(t *testing.T) {
	// Select a test port and start a dummy TCP server on it.
	testPort := 8000
	startTestServer(t, testPort)

	// Create a scanner instance that scans only the test port.
	s := Scanner{
		Target:    "127.0.0.1",
		StartPort: testPort,
		EndPort:   testPort,
		Workers:   1,
		Timeout:   2 * time.Second,
	}

	summary := s.Run()
	if len(summary.OpenPorts) == 0 {
		t.Errorf("Expected port %d to be open, but got none", testPort)
	}

	found := slices.Contains(summary.OpenPorts, testPort)
	if !found {
		t.Errorf("Expected port %d to be in open ports. Got: %v", testPort, summary.OpenPorts)
	}
}
