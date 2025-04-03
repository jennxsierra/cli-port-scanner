// Filename: scanner.go
// Description: Contains the port scanning logic and data structures.
package scanner

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/jennxsierra/cli-port-scanner/internal/banner"
)

// Config holds scanning parameters.
type Config struct {
	Target     string
	Ports      []int
	Workers    int
	Timeout    time.Duration
	MaxRetries int
}

// PortResult holds the scanning result for a single port.
type PortResult struct {
	Port   int    `json:"port"`
	Banner string `json:"banner,omitempty"`
}

// Result holds the overall scan result for a target.
type Result struct {
	Target     string       `json:"target"`
	OpenPorts  []PortResult `json:"open_ports"`
	TotalPorts int          `json:"total_ports"`
	OpenCount  int          `json:"open_count"`
	Duration   string       `json:"duration"`
}

// ScanTarget scans the given target based on the provided config.
// It sends an integer (1) on the provided progress channel for every port processed.
func ScanTarget(cfg Config, progress chan<- int) Result {
	var mu sync.Mutex
	openPorts := []PortResult{}
	totalPorts := len(cfg.Ports)

	startTime := time.Now()

	tasks := make(chan int, totalPorts)
	var wg sync.WaitGroup

	dialer := net.Dialer{
		Timeout: cfg.Timeout,
	}

	// Worker function that processes a port from the tasks channel.
	worker := func() {
		defer wg.Done()
		for port := range tasks {
			addr := net.JoinHostPort(cfg.Target, strconv.Itoa(port))
			var bannerResult string
			success := false
			for i := 0; i < cfg.MaxRetries; i++ {
				conn, err := dialer.Dial("tcp", addr)
				if err == nil {
					bannerResult = banner.GrabBanner(conn, port)
					conn.Close()
					success = true
					break
				}
				backoff := time.Duration(1<<i) * time.Second
				time.Sleep(backoff)
			}
			if success {
				mu.Lock()
				openPorts = append(openPorts, PortResult{Port: port, Banner: bannerResult})
				mu.Unlock()
			}
			// Report progress.
			progress <- 1
		}
	}

	// Launch worker goroutines.
	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go worker()
	}

	// Send ports to the tasks channel.
	go func() {
		for _, port := range cfg.Ports {
			tasks <- port
		}
		close(tasks)
	}()

	wg.Wait()

	// Sort open ports in ascending order.
	sort.Slice(openPorts, func(i, j int) bool {
		return openPorts[i].Port < openPorts[j].Port
	})

	duration := time.Since(startTime)

	return Result{
		Target:     cfg.Target,
		OpenPorts:  openPorts,
		TotalPorts: totalPorts,
		OpenCount:  len(openPorts),
		Duration:   fmt.Sprintf("%.3fs", duration.Seconds()),
	}
}

// ToJSON returns the JSON string of the scan Result.
func (r Result) ToJSON() (string, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
