package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/jennxsierra/cli-port-scanner/internal/scanner"
)

// captureOutput redirects stdout and returns printed output.
func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	f()
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old
	return buf.String()
}

func TestPrintHeader(t *testing.T) {
	target := "example.com"
	ports := []int{80, 443}
	output := captureOutput(func() {
		PrintHeader(target, ports)
	})
	if !strings.Contains(output, "SCAN START: "+target) {
		t.Errorf("Expected header to contain target name, got %q", output)
	}
}

func TestPrintSummaryAndBanners(t *testing.T) {
	// Create a dummy scan.Result.
	result := scanner.Result{
		Target: "example.com",
		OpenPorts: []scanner.PortResult{
			{Port: 80, Banner: "TestBanner"},
			{Port: 443, Banner: ""},
		},
		TotalPorts: 2,
		OpenCount:  1,
		Duration:   "0.123s",
	}
	ports := []int{80, 443}

	headerOut := captureOutput(func() {
		PrintSummary(result, ports)
	})
	if !strings.Contains(headerOut, "example.com") {
		t.Errorf("PrintSummary output did not contain target")
	}
	if !strings.Contains(headerOut, "2") {
		t.Errorf("PrintSummary output did not mention total ports")
	}

	bannerOut := captureOutput(func() {
		PrintBanners(result)
	})
	if !strings.Contains(bannerOut, "TestBanner") {
		t.Errorf("PrintBanners output did not contain expected banner")
	}
	if !strings.Contains(bannerOut, "<no banner>") {
		t.Errorf("PrintBanners output did not mark empty banner correctly")
	}
}
