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
	fmt.Printf("                 SCAN START: %s\n", target)
	fmt.Println(line)
	fmt.Printf("Scanning ports: %s\n", formatPorts(ports))
}

// PrintSummary prints the scan summary in a formatted box.
func PrintSummary(result scanner.Result, ports []int) {
	scanned := formatPorts(ports)
	openPorts := []string{}
	for _, pr := range result.OpenPorts {
		openPorts = append(openPorts, strconv.Itoa(pr.Port))
	}
	summary := `
┌────────────────────────────────────────────────────┐
│                     SCAN SUMMARY                   │
├────────────────────────────────────────────────────┤
│ Target          : %-40s│
│ Ports Scanned   : %-40s│
│ Open Ports      : %-40s│
│ Time Taken      : %-40s│
└────────────────────────────────────────────────────┘`
	fmt.Printf(summary+"\n", result.Target, scanned, strings.Join(openPorts, ", "), result.Duration)
}

// PrintBanners prints the banners received on open ports.
func PrintBanners(result scanner.Result) {
	fmt.Println(`
┌────────────────────────────────────────────────────┐
│                     BANNERS                        │
├────────────────────────────────────────────────────┤`)
	for _, pr := range result.OpenPorts {
		banner := pr.Banner
		if banner == "" {
			banner = "<no banner>"
		}
		fmt.Printf("│ Port %-13d: \"%-30s\"│\n", pr.Port, banner)
	}
	fmt.Println("└────────────────────────────────────────────────────┘")
}

func formatPorts(ports []int) string {
	strPorts := []string{}
	for _, p := range ports {
		strPorts = append(strPorts, strconv.Itoa(p))
	}
	return strings.Join(strPorts, ", ")
}
