package osint

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TwiN/go-color"
)

// Spinner provides an animated loading spinner for indeterminate operations.
type Spinner struct {
	chars    []string
	index    int
	message  string
	stopChan chan bool
	doneChan chan bool
	running  bool
}

// NewSpinner creates a new spinner with a custom message.
func NewSpinner(message string) *Spinner {
	return &Spinner{
		chars:    []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		index:    0,
		message:  message,
		stopChan: make(chan bool),
		doneChan: make(chan bool),
		running:  false,
	}
}

// Start begins the spinner animation in a goroutine.
func (s *Spinner) Start() {
	if s.running {
		return
	}
	s.running = true

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-s.stopChan:
				s.doneChan <- true
				return
			case <-ticker.C:
				fmt.Printf("\r%s %s", color.Ize(color.Cyan, s.chars[s.index]), color.Ize(color.Cyan, s.message))
				s.index = (s.index + 1) % len(s.chars)
			}
		}
	}()
}

// Stop stops the spinner and clears the line.
func (s *Spinner) Stop() {
	if !s.running {
		return
	}
	s.stopChan <- true
	<-s.doneChan
	fmt.Print("\r" + strings.Repeat(" ", len(s.message)+10) + "\r")
	s.running = false
}

// UpdateMessage updates the spinner message while it's running.
func (s *Spinner) UpdateMessage(message string) {
	s.message = message
}

// ProgressBar provides a progress bar for operations with known progress.
type ProgressBar struct {
	total     int
	current   int
	width     int
	message   string
	completed bool
}

// NewProgressBar creates a new progress bar.
func NewProgressBar(total int, message string) *ProgressBar {
	return &ProgressBar{
		total:     total,
		current:   0,
		width:     40,
		message:   message,
		completed: false,
	}
}

// Update updates the progress bar with the current progress.
func (pb *ProgressBar) Update(current int) {
	if pb.completed {
		return
	}
	pb.current = current
	if pb.current > pb.total {
		pb.current = pb.total
	}
	pb.render()
}

// Increment increments the progress by 1.
func (pb *ProgressBar) Increment() {
	if pb.completed {
		return
	}
	pb.current++
	if pb.current > pb.total {
		pb.current = pb.total
	}
	pb.render()
}

// SetTotal updates the total value.
func (pb *ProgressBar) SetTotal(total int) {
	pb.total = total
	if pb.current > pb.total {
		pb.current = pb.total
	}
	pb.render()
}

// render renders the progress bar.
func (pb *ProgressBar) render() {
	percentage := float64(pb.current) / float64(pb.total) * 100
	filled := int(float64(pb.width) * float64(pb.current) / float64(pb.total))
	empty := pb.width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	fmt.Printf("\r%s [%s] %d/%d (%.1f%%)", 
		color.Ize(color.Cyan, pb.message),
		color.Ize(color.Green, bar),
		pb.current,
		pb.total,
		percentage)
}

// Complete marks the progress bar as complete and renders the final state.
func (pb *ProgressBar) Complete() {
	if pb.completed {
		return
	}
	pb.current = pb.total
	pb.completed = true
	pb.render()
	fmt.Println() // New line after completion
}

// ShowProgress shows a simple progress message for operations.
func ShowProgress(message string) {
	fmt.Print(color.Ize(color.Cyan, "  [*] "+message+"..."))
}

// HideProgress clears the progress message.
func HideProgress() {
	fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
}

// ShowProgressWithSpinner shows a progress message with a spinner.
func ShowProgressWithSpinner(message string) *Spinner {
	spinner := NewSpinner(message)
	spinner.Start()
	return spinner
}

// ShowProgressWithBar shows a progress message with a progress bar.
func ShowProgressWithBar(total int, message string) *ProgressBar {
	return NewProgressBar(total, message)
}

// ShowAPIProgress shows a progress indicator for API calls.
func ShowAPIProgress(operation string) *Spinner {
	message := fmt.Sprintf("Fetching %s", operation)
	return ShowProgressWithSpinner(message)
}

// ShowBatchProgress shows progress for batch operations.
func ShowBatchProgress(current, total int, item string) {
	message := fmt.Sprintf("Processing batch: %s", item)
	pb := NewProgressBar(total, message)
	pb.Update(current)
	if current == total {
		pb.Complete()
	}
}

// ShowLoginProgress shows progress for login operations.
func ShowLoginProgress() *Spinner {
	return ShowProgressWithSpinner("Authenticating with Space-Track")
}

// ShowQueryProgress shows progress for query operations.
func ShowQueryProgress(endpoint string) *Spinner {
	// Extract a readable description from the endpoint
	desc := "satellite data"
	if strings.Contains(endpoint, "satcat") {
		desc = "satellite catalog"
	} else if strings.Contains(endpoint, "tle") {
		desc = "TLE data"
	} else if strings.Contains(endpoint, "gp_history") {
		desc = "orbital history"
	}
	return ShowProgressWithSpinner(fmt.Sprintf("Querying %s", desc))
}

// ShowDownloadProgress shows progress for download operations.
func ShowDownloadProgress(filename string) *Spinner {
	return ShowProgressWithSpinner(fmt.Sprintf("Downloading %s", filename))
}

// IsTerminal checks if stdout is a terminal (for progress display).
func IsTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// ShowSimpleProgress shows a simple progress message that works in all environments.
func ShowSimpleProgress(message string) {
	if IsTerminal() {
		fmt.Print(color.Ize(color.Cyan, "  [*] "+message+"..."))
	} else {
		fmt.Println(color.Ize(color.Cyan, "  [*] "+message+"..."))
	}
}

// HideSimpleProgress clears simple progress messages.
func HideSimpleProgress() {
	if IsTerminal() {
		fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
	}
}

