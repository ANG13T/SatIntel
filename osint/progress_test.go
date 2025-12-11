package osint

import (
	"testing"
	"time"
)

func TestNewSpinner(t *testing.T) {
	spinner := NewSpinner("Test message")
	if spinner == nil {
		t.Fatal("NewSpinner() returned nil")
	}
	if spinner.message != "Test message" {
		t.Errorf("message = %q, want %q", spinner.message, "Test message")
	}
	if spinner.running {
		t.Error("Spinner should not be running initially")
	}
}

func TestSpinnerStartStop(t *testing.T) {
	spinner := NewSpinner("Testing")
	
	// Start spinner
	spinner.Start()
	if !spinner.running {
		t.Error("Spinner should be running after Start()")
	}
	
	// Wait a bit to ensure it's animating
	time.Sleep(150 * time.Millisecond)
	
	// Stop spinner
	spinner.Stop()
	if spinner.running {
		t.Error("Spinner should not be running after Stop()")
	}
}

func TestSpinnerUpdateMessage(t *testing.T) {
	spinner := NewSpinner("Original")
	if spinner.message != "Original" {
		t.Errorf("Initial message = %q, want %q", spinner.message, "Original")
	}
	
	spinner.UpdateMessage("Updated")
	if spinner.message != "Updated" {
		t.Errorf("Updated message = %q, want %q", spinner.message, "Updated")
	}
}

func TestNewProgressBar(t *testing.T) {
	pb := NewProgressBar(100, "Test progress")
	if pb == nil {
		t.Fatal("NewProgressBar() returned nil")
	}
	if pb.total != 100 {
		t.Errorf("total = %d, want %d", pb.total, 100)
	}
	if pb.current != 0 {
		t.Errorf("current = %d, want %d", pb.current, 0)
	}
	if pb.message != "Test progress" {
		t.Errorf("message = %q, want %q", pb.message, "Test progress")
	}
	if pb.completed {
		t.Error("ProgressBar should not be completed initially")
	}
}

func TestProgressBarUpdate(t *testing.T) {
	pb := NewProgressBar(100, "Test")
	
	pb.Update(50)
	if pb.current != 50 {
		t.Errorf("current = %d, want %d", pb.current, 50)
	}
	
	pb.Update(150) // Should cap at total
	if pb.current != 100 {
		t.Errorf("current = %d, want %d (should be capped)", pb.current, 100)
	}
}

func TestProgressBarIncrement(t *testing.T) {
	pb := NewProgressBar(10, "Test")
	
	if pb.current != 0 {
		t.Errorf("Initial current = %d, want %d", pb.current, 0)
	}
	
	pb.Increment()
	if pb.current != 1 {
		t.Errorf("After increment, current = %d, want %d", pb.current, 1)
	}
	
	// Increment multiple times
	for i := 0; i < 5; i++ {
		pb.Increment()
	}
	if pb.current != 6 {
		t.Errorf("After 6 increments, current = %d, want %d", pb.current, 6)
	}
	
	// Increment beyond total
	for i := 0; i < 10; i++ {
		pb.Increment()
	}
	if pb.current != 10 {
		t.Errorf("After exceeding total, current = %d, want %d (capped)", pb.current, 10)
	}
}

func TestProgressBarSetTotal(t *testing.T) {
	pb := NewProgressBar(100, "Test")
	pb.Update(50)
	
	pb.SetTotal(200)
	if pb.total != 200 {
		t.Errorf("total = %d, want %d", pb.total, 200)
	}
	if pb.current != 50 {
		t.Errorf("current = %d, want %d (should remain same)", pb.current, 50)
	}
	
	// Set total lower than current
	pb.SetTotal(30)
	if pb.current != 30 {
		t.Errorf("current = %d, want %d (should be capped to new total)", pb.current, 30)
	}
}

func TestProgressBarComplete(t *testing.T) {
	pb := NewProgressBar(100, "Test")
	pb.Update(50)
	
	if pb.completed {
		t.Error("ProgressBar should not be completed before Complete()")
	}
	
	pb.Complete()
	if !pb.completed {
		t.Error("ProgressBar should be completed after Complete()")
	}
	if pb.current != pb.total {
		t.Errorf("current = %d, want %d (should equal total)", pb.current, pb.total)
	}
}

func TestShowProgress(t *testing.T) {
	// This is a simple function that just prints, so we test it doesn't panic
	ShowProgress("Test message")
	HideProgress()
}

func TestShowProgressWithSpinner(t *testing.T) {
	spinner := ShowProgressWithSpinner("Test")
	if spinner == nil {
		t.Fatal("ShowProgressWithSpinner() returned nil")
	}
	if !spinner.running {
		t.Error("Spinner should be running")
	}
	spinner.Stop()
}

func TestShowProgressWithBar(t *testing.T) {
	pb := ShowProgressWithBar(100, "Test")
	if pb == nil {
		t.Fatal("ShowProgressWithBar() returned nil")
	}
	if pb.total != 100 {
		t.Errorf("total = %d, want %d", pb.total, 100)
	}
}

func TestShowAPIProgress(t *testing.T) {
	spinner := ShowAPIProgress("test operation")
	if spinner == nil {
		t.Fatal("ShowAPIProgress() returned nil")
	}
	if spinner.message != "Fetching test operation" {
		t.Errorf("message = %q, want %q", spinner.message, "Fetching test operation")
	}
	spinner.Stop()
}

func TestShowLoginProgress(t *testing.T) {
	spinner := ShowLoginProgress()
	if spinner == nil {
		t.Fatal("ShowLoginProgress() returned nil")
	}
	if spinner.message != "Authenticating with Space-Track" {
		t.Errorf("message = %q, want %q", spinner.message, "Authenticating with Space-Track")
	}
	spinner.Stop()
}

func TestShowQueryProgress(t *testing.T) {
	spinner := ShowQueryProgress("/class/satcat")
	if spinner == nil {
		t.Fatal("ShowQueryProgress() returned nil")
	}
	if spinner.message != "Querying satellite catalog" {
		t.Errorf("message = %q, want %q", spinner.message, "Querying satellite catalog")
	}
	spinner.Stop()
	
	spinner2 := ShowQueryProgress("/class/gp_history/format/tle")
	if spinner2.message != "Querying TLE data" {
		t.Errorf("message = %q, want %q", spinner2.message, "Querying TLE data")
	}
	spinner2.Stop()
}

func TestShowBatchProgress(t *testing.T) {
	// This function prints progress, so we test it doesn't panic
	ShowBatchProgress(5, 10, "Satellite 1")
	ShowBatchProgress(10, 10, "Satellite 2")
}

func TestIsTerminal(t *testing.T) {
	// This function checks if stdout is a terminal
	// We can't easily test the actual result, but we can test it doesn't panic
	result := IsTerminal()
	_ = result // Use result to avoid unused variable
}

func TestShowSimpleProgress(t *testing.T) {
	// Test it doesn't panic
	ShowSimpleProgress("Test message")
	HideSimpleProgress()
}

// Benchmark tests
func BenchmarkSpinnerStartStop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		spinner := NewSpinner("Benchmark")
		spinner.Start()
		time.Sleep(10 * time.Millisecond)
		spinner.Stop()
	}
}

func BenchmarkProgressBarUpdate(b *testing.B) {
	pb := NewProgressBar(1000, "Benchmark")
	for i := 0; i < b.N; i++ {
		pb.Update(i % 1001)
	}
}

func BenchmarkProgressBarIncrement(b *testing.B) {
	pb := NewProgressBar(b.N, "Benchmark")
	for i := 0; i < b.N; i++ {
		pb.Increment()
	}
}

