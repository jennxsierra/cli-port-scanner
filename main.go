// Filename: main.go
// Description: A simple command-line port scanner in Go.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jennxsierra/cli-port-scanner/internal/progress"
	"github.com/jennxsierra/cli-port-scanner/internal/scanner"
	"github.com/jennxsierra/cli-port-scanner/internal/ui"
)

func parsePorts(portFlag string, startPort, endPort int) []int {
	var ports []int
	if portFlag != "" {
		// Use strings.Split (not SplitSeq) to properly parse comma-separated ports.
		parts := strings.Split(portFlag, ",")
		for _, part := range parts {
			p, err := strconv.Atoi(strings.TrimSpace(part))
			if err == nil {
				ports = append(ports, p)
			}
		}
	} else {
		for p := startPort; p <= endPort; p++ {
			ports = append(ports, p)
		}
	}
	return ports
}

func parseTargets(targetFlag, targetsFlag string) []string {
	targetSet := make(map[string]struct{})
	if targetFlag != "" {
		targetSet[targetFlag] = struct{}{}
	}
	if targetsFlag != "" {
		parts := strings.Split(targetsFlag, ",")
		for _, part := range parts {
			t := strings.TrimSpace(part)
			if t != "" {
				targetSet[t] = struct{}{}
			}
		}
	}
	var targets []string
	for t := range targetSet {
		targets = append(targets, t)
	}
	return targets
}

func main() {
	// Define command-line flags.
	targetFlag := flag.String("target", "", "The hostname or IP address to be scanned.")
	targetsFlag := flag.String("targets", "", "Comma-separated list of targets (e.g., -targets=localhost,scanme.nmap.org).")
	startPort := flag.Int("start-port", 1, "The lower bound port to begin scanning. (default 1)")
	endPort := flag.Int("end-port", 1024, "The upper bound port to finish scanning. (default 1024)")
	workers := flag.Int("workers", 100, "The number of concurrent goroutines to launch per target. (default 100)")
	timeout := flag.Int("timeout", 5, "The connection timeout in seconds. (default 5)")
	portsFlag := flag.String("ports", "", "Comma-separated list of specific ports (overrides start-port and end-port).")
	jsonFlag := flag.Bool("json", false, "Output results in JSON format.")
	debugFlag := flag.Bool("debug", false, "Display flag values for debugging.")
	flag.Parse()

	if *debugFlag {
		fmt.Println("Flag values:")
		fmt.Printf("target: %s\n", *targetFlag)
		fmt.Printf("targets: %s\n", *targetsFlag)
		fmt.Printf("start-port: %d\n", *startPort)
		fmt.Printf("end-port: %d\n", *endPort)
		fmt.Printf("workers: %d\n", *workers)
		fmt.Printf("timeout: %d seconds\n", *timeout)
		fmt.Printf("ports: %s\n", *portsFlag)
		fmt.Printf("json: %v\n", *jsonFlag)
		fmt.Printf("debug: %v\n", *debugFlag)
	}

	// Validate and aggregate targets.
	targets := parseTargets(*targetFlag, *targetsFlag)
	if len(targets) == 0 {
		fmt.Println("No target specified. Please provide -target or -targets flag.")
		os.Exit(1)
	}

	// Parse port list (or use provided range).
	portList := parsePorts(*portsFlag, *startPort, *endPort)
	if len(portList) == 0 {
		fmt.Println("No valid ports provided.")
		os.Exit(1)
	}

	overallResults := make([]scanner.Result, 0)

	// Process each target.
	for _, tgt := range targets {
		// Print header using UI module.
		ui.PrintHeader(tgt, portList)

		// Create a progress channel and start one dynamic progress bar per target.
		totalPorts := len(portList)
		progressChan := make(chan int, totalPorts)
		pb := progress.NewProgressBar(totalPorts, 20)
		go pb.Start(progressChan)

		// Configure and run the scan. Pass the progress channel.
		cfg := scanner.Config{
			Target:     tgt,
			Ports:      portList,
			Workers:    *workers,
			Timeout:    time.Duration(*timeout) * time.Second,
			MaxRetries: 3,
		}
		result := scanner.ScanTarget(cfg, progressChan)
		// Close progress channel so progress bar finishes.
		close(progressChan)

		overallResults = append(overallResults, result)

		// Use UI module to print scan summary and banners.
		ui.PrintSummary(result, portList)
		ui.PrintBanners(result)
		fmt.Println()
	}

	// Write JSON file if the -json flag is set.
	if *jsonFlag {
		jsonData, err := json.MarshalIndent(overallResults, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			os.Exit(1)
		}
		// Filename format: DDMMYY-HHMMSS-cli-pscan.json
		fileName := time.Now().Format("020106-150405") + "-cli-pscan.json"
		err = os.WriteFile(fileName, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing JSON file:", err)
			os.Exit(1)
		}
		fmt.Printf("Scan results saved to %s\n", fileName)
	}
}
