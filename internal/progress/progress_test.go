package progress

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

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

func TestProgressBar(t *testing.T) {
	total := 10
	pb := NewProgressBar(total, 10)
	progressChan := make(chan int, total)
	done := make(chan struct{})

	// Start the progress bar in a goroutine.
	go pb.Start(progressChan, done)

	// Send progress increments.
	for i := 0; i < total; i++ {
		progressChan <- 1
		time.Sleep(10 * time.Millisecond)
	}
	close(progressChan)

	// Wait for progress bar to finish.
	<-done

	out := captureOutput(func() {
		// Final print to trigger end-line
		pb.print()
	})
	if !strings.Contains(out, "100%") {
		t.Errorf("Expected progress bar to show 100%%, got: %s", out)
	}
}
