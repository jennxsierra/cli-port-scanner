package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jennxsierra/cli-port-scanner/internal/scanner"
)

func main() {
	// New command-line flags for enriched scanning.
	targetsFlag := flag.String("targets", "scanme.nmap.org", "Comma-separated list of targets to scan")
	startPort := flag.Int("start-port", 1, "Start port for scanning (ignored if -ports provided)")
	endPort := flag.Int("end-port", 1024, "End port for scanning (ignored if -ports provided)")
	portsFlag := flag.String("ports", "", "Comma-separated list of specific ports to scan")
	workers := flag.Int("workers", 100, "Number of concurrent workers")
	timeout := flag.Int("timeout", 5, "Connection timeout in seconds per port")
	jsonOutput := flag.Bool("json", false, "Output results in JSON format")
	flag.Parse()

	// Process targets.
	targets := strings.Split(*targetsFlag, ",")

	// Process ports.
	var portList []int
	if *portsFlag != "" {
		for _, p := range strings.Split(*portsFlag, ",") {
			port, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil {
				fmt.Printf("Invalid port value: %s\n", p)
				return
			}
			portList = append(portList, port)
		}
	} else {
		for p := *startPort; p <= *endPort; p++ {
			portList = append(portList, p)
		}
	}

	// Create and run the enriched scanner.
	es := scanner.EnrichedScanner{
		Targets:    targets,
		Ports:      portList,
		Workers:    *workers,
		Timeout:    time.Duration(*timeout) * time.Second,
		JsonOutput: *jsonOutput,
	}

	summaries, results := es.Run()

	// Print a newline for clarity.
	fmt.Println()

	// [BANNERS] Section: print banner lines for open ports.
	fmt.Println("[BANNERS]")
	for _, r := range results {
		if r.Status == "open" && r.Banner != "" {
			fmt.Printf("[%s:%d] %s\n", r.Target, r.Port, r.Banner)
		}
	}
	fmt.Println()

	// [SCAN SUMMARY] Section.
	fmt.Println("[SCAN SUMMARY]")
	for _, sum := range summaries {
		fmt.Printf("\n[%s]\nTotal Ports Scanned: %d\nOpen Ports Count: %d\nOpen Ports: %v\nTime Taken: %s\n",
			sum.Target, sum.TotalPorts, len(sum.OpenPorts), sum.OpenPorts, sum.Duration.Round(time.Millisecond).String())
	}

	// If JSON output is requested, save results.
	if *jsonOutput {
		filename, err := scanner.SaveJSON(results)
		if err != nil {
			fmt.Printf("Error saving JSON output: %v\n", err)
		} else {
			fmt.Printf("\n[JSON OUTPUT SAVED: %s]\n", filename)
		}
	}
}
