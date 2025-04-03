package scanner

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type EnrichedResult struct {
	Target string `json:"target"`
	Port   int    `json:"port"`
	Status string `json:"status"`
	Banner string `json:"banner,omitempty"`
}

type Summary struct {
	Target     string        `json:"target"`
	TotalPorts int           `json:"total_ports_scanned"`
	OpenPorts  []int         `json:"open_ports"`
	Duration   time.Duration `json:"time_taken"`
}

// progressBar returns a simple progress bar string.
func progressBar(progress, total int, width int) string {
	barCount := min(int(float64(progress)/float64(total)*float64(width)), width)
	bar := strings.Repeat("=", barCount)
	space := strings.Repeat(" ", width-barCount)
	return fmt.Sprintf("[%s%s]", bar, space)
}

type EnrichedScanner struct {
	Targets    []string
	Ports      []int
	Workers    int
	Timeout    time.Duration
	JsonOutput bool
}

func (es *EnrichedScanner) Run() ([]Summary, []EnrichedResult) {
	var summaries []Summary
	var allResults []EnrichedResult

	fmt.Println("[SCAN START]")

	// Process targets sequentially (you can also parallelize if needed).
	for _, target := range es.Targets {
		totalPorts := len(es.Ports)
		tasks := make(chan int, totalPorts)
		results := make(chan EnrichedResult, totalPorts)
		var wg sync.WaitGroup
		dialer := net.Dialer{Timeout: es.Timeout}

		// Launch worker goroutines.
		for range es.Workers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for port := range tasks {
					addr := net.JoinHostPort(target, strconv.Itoa(port))
					res := EnrichedResult{
						Target: target,
						Port:   port,
						Status: "closed",
					}
					conn, err := dialer.Dial("tcp", addr)
					if err == nil {
						res.Status = "open"
						// Attempt to read initial banner.
						conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
						buffer := make([]byte, 1024)
						n, _ := conn.Read(buffer)
						res.Banner = strings.TrimSpace(string(buffer[:n]))
						conn.Close()
					}
					results <- res
				}
			}()
		}

		// Enqueue ports.
		for _, port := range es.Ports {
			tasks <- port
		}
		close(tasks)

		startTime := time.Now()
		scanned := 0
		var targetResults []EnrichedResult

		// Collect results while updating progress.
		for scanned < totalPorts {
			res := <-results
			targetResults = append(targetResults, res)
			scanned++
			// Build and print a progress bar (width=60).
			perc := int((float64(scanned) / float64(totalPorts)) * 100)
			bar := progressBar(scanned, totalPorts, 60)
			// \r will overwrite the same line (works in most terminals).
			fmt.Printf("\r[%s] [%d/%d ports] %s %3d %%", target, scanned, totalPorts, bar, perc)
		}
		fmt.Println()

		// Ensure all workers are done.
		wg.Wait()
		close(results)

		var openPorts []int
		for _, r := range targetResults {
			if r.Status == "open" {
				openPorts = append(openPorts, r.Port)
			}
		}
		duration := time.Since(startTime)
		summary := Summary{
			Target:     target,
			TotalPorts: totalPorts,
			OpenPorts:  openPorts,
			Duration:   duration,
		}
		summaries = append(summaries, summary)
		allResults = append(allResults, targetResults...)
	}

	return summaries, allResults
}

// SaveJSON saves the given results to a JSON file whose name is based on the current timestamp.
func SaveJSON(results []EnrichedResult) (string, error) {
	filename := time.Now().Format("20060102-150405-cli-pscan.json")
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		return "", err
	}
	return filename, nil
}
