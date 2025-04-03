// Filename: ui.go
// Description: Contains functions for printing the scan results and progress.
package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jennxsierra/cli-port-scanner/internal/scanner"
)

// PrintHeader prints a formatted header for the scan.
func PrintHeader(target string, ports []int) {
	line := "===================================================="
	fmt.Println(line)
	fmt.Printf("               SCAN START: %s\n", target)
	fmt.Println(line)
	fmt.Printf("Scanning ports: %s\n", formatPorts(ports))
}

// PrintSummary prints the scan summary with a simple title.
func PrintSummary(result scanner.Result, ports []int) {
	scanned := formatPorts(ports)
	openPorts := []string{}
	for _, pr := range result.OpenPorts {
		openPorts = append(openPorts, strconv.Itoa(pr.Port))
	}

	fmt.Println("\n[SCAN SUMMARY]")
	fmt.Printf("Target         : %s\n", result.Target)
	fmt.Printf("Ports Scanned  : %s\n", scanned)
	fmt.Printf("Open Ports     : %s\n", strings.Join(openPorts, ", "))
	fmt.Printf("Time Taken     : %s\n", result.Duration)
}

// PrintBanners prints the banners received on open ports with a simple title.
func PrintBanners(result scanner.Result) {
	fmt.Println("\n[BANNERS]")
	for _, pr := range result.OpenPorts {
		banner := pr.Banner
		if banner == "" {
			banner = "<no banner>"
		}
		fmt.Printf("%-10s: %q\n", fmt.Sprintf("Port %d", pr.Port), banner)
	}
}

func formatPorts(ports []int) string {
	strPorts := []string{}
	for _, p := range ports {
		strPorts = append(strPorts, strconv.Itoa(p))
	}
	return strings.Join(strPorts, ", ")
}
