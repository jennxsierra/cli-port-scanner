// Filename: progress.go
// Description: Contains the ProgressBar struct and methods for displaying progress.
package progress

import (
	"fmt"
	"strings"
	"time"
)

// ProgressBar represents a simple dynamic progress bar.
type ProgressBar struct {
	Total   int
	Current int
	Width   int
}

// NewProgressBar returns a new progress bar.
// width is the number of blocks to display.
func NewProgressBar(total, width int) *ProgressBar {
	return &ProgressBar{
		Total: total,
		Width: width,
	}
}

// Start listens on the provided channel for progress increments.
// This function runs until the total number of increments is received.
func (pb *ProgressBar) Start(progressChan <-chan int) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case incr, ok := <-progressChan:
			if !ok {
				pb.print() // final print
				return
			}
			pb.Current += incr
		case <-ticker.C:
			pb.print()
		}
		if pb.Current >= pb.Total {
			pb.print()
			return
		}
	}
}

// print outputs the progress bar.
func (pb *ProgressBar) print() {
	percent := float64(pb.Current) / float64(pb.Total)
	filled := int(percent * float64(pb.Width))
	if filled > pb.Width {
		filled = pb.Width
	}
	empty := pb.Width - filled
	bar := strings.Repeat("█", filled) + strings.Repeat(" ", empty)
	fmt.Printf("\rProgress: [%s] %3.0f%%", bar, percent*100)
	if pb.Current >= pb.Total {
		fmt.Print("\n")
	}
}
