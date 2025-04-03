package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jennxsierra/cli-port-scanner/internal/scanner"
)

func main() {
	// Command-line flags
	target := flag.String("target", "scanme.nmap.org", "Target IP address or hostname")
	startPort := flag.Int("start-port", 1, "Start of port range")
	endPort := flag.Int("end-port", 1024, "End of port range")
	workers := flag.Int("workers", 100, "Number of concurrent workers")
	timeout := flag.Int("timeout", 5, "Timeout in seconds for each connection")
	flag.Parse()

	// Create a Scanner instance using the provided flags.
	s := scanner.Scanner{
		Target:    *target,
		StartPort: *startPort,
		EndPort:   *endPort,
		Workers:   *workers,
		Timeout:   time.Duration(*timeout) * time.Second,
	}

	// Run the scan.
	summary := s.Run()

	// Print the summary header.
	fmt.Println("[SCAN SUMMARY]")
	fmt.Println(scanner.FormatSummary(summary))
}
