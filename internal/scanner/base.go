package scanner

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type Summary struct {
	Target     string
	TotalPorts int
	OpenPorts  []int
	Duration   time.Duration
}

type Scanner struct {
	Target    string
	StartPort int
	EndPort   int
	Workers   int
	Timeout   time.Duration
}

func (s *Scanner) Run() Summary {
	totalPorts := s.EndPort - s.StartPort + 1
	tasks := make(chan int, totalPorts)
	openChan := make(chan int, totalPorts)
	var wg sync.WaitGroup
	dialer := net.Dialer{Timeout: s.Timeout}

	// Start worker routines.
	for range s.Workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range tasks {
				addr := net.JoinHostPort(s.Target, strconv.Itoa(port))
				conn, err := dialer.Dial("tcp", addr)
				if err == nil {
					openChan <- port
					conn.Close()
				}
			}
		}()
	}

	startTime := time.Now()

	// Enqueue all ports.
	for port := s.StartPort; port <= s.EndPort; port++ {
		tasks <- port
	}
	close(tasks)
	wg.Wait()
	close(openChan)

	var openPorts []int
	for port := range openChan {
		openPorts = append(openPorts, port)
	}

	return Summary{
		Target:     s.Target,
		TotalPorts: totalPorts,
		OpenPorts:  openPorts,
		Duration:   time.Since(startTime),
	}
}

// FormatSummary returns the summary formatted as requested.
func FormatSummary(summary Summary) string {
	return fmt.Sprintf(`[%s]
Total Ports Scanned: %d
Open Ports Count: %d
Open Ports: %v
Time Taken: %s
`, summary.Target, summary.TotalPorts, len(summary.OpenPorts), summary.OpenPorts, summary.Duration.Round(time.Millisecond).String())
}
