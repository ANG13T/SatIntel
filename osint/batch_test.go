package osint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBatchSatelliteStruct(t *testing.T) {
	sat := BatchSatellite{
		Name:       "ISS (ZARYA)",
		NORADID:    "25544",
		Country:    "US",
		ObjectType: "PAYLOAD",
	}

	if sat.Name != "ISS (ZARYA)" {
		t.Errorf("Name = %q, want %q", sat.Name, "ISS (ZARYA)")
	}
	if sat.NORADID != "25544" {
		t.Errorf("NORADID = %q, want %q", sat.NORADID, "25544")
	}
}

func TestCompareSatellites(t *testing.T) {
	results := []BatchTLEResult{
		{
			Satellite: BatchSatellite{Name: "Sat1", NORADID: "12345"},
			Success:   true,
			TLE: TLE{
				OrbitInclination: 51.6,
				MeanMotion:        15.5,
				Eccentrcity:       0.0001,
			},
		},
		{
			Satellite: BatchSatellite{Name: "Sat2", NORADID: "12346"},
			Success:   true,
			TLE: TLE{
				OrbitInclination: 98.2,
				MeanMotion:        14.2,
				Eccentrcity:       0.0002,
			},
		},
		{
			Satellite: BatchSatellite{Name: "Sat3", NORADID: "12347"},
			Success:   false,
			Error:     fmt.Errorf("test error"),
		},
	}

	comparison := CompareSatellites(results)

	if comparison.Summary.TotalProcessed != 3 {
		t.Errorf("TotalProcessed = %d, want %d", comparison.Summary.TotalProcessed, 3)
	}
	if comparison.Summary.Successful != 2 {
		t.Errorf("Successful = %d, want %d", comparison.Summary.Successful, 2)
	}
	if comparison.Summary.Failed != 1 {
		t.Errorf("Failed = %d, want %d", comparison.Summary.Failed, 1)
	}

	// Check average inclination
	expectedAvg := (51.6 + 98.2) / 2.0
	if comparison.Summary.AverageInclination != expectedAvg {
		t.Errorf("AverageInclination = %.2f, want %.2f", comparison.Summary.AverageInclination, expectedAvg)
	}

	// Check average mean motion
	expectedMM := (15.5 + 14.2) / 2.0
	if comparison.Summary.AverageMeanMotion != expectedMM {
		t.Errorf("AverageMeanMotion = %.4f, want %.4f", comparison.Summary.AverageMeanMotion, expectedMM)
	}
}

func TestCompareSatellitesEmpty(t *testing.T) {
	results := []BatchTLEResult{}
	comparison := CompareSatellites(results)

	if comparison.Summary.TotalProcessed != 0 {
		t.Errorf("TotalProcessed = %d, want %d", comparison.Summary.TotalProcessed, 0)
	}
}

func TestCompareSatellitesAllFailed(t *testing.T) {
	results := []BatchTLEResult{
		{
			Satellite: BatchSatellite{Name: "Sat1", NORADID: "12345"},
			Success:   false,
			Error:     fmt.Errorf("error 1"),
		},
		{
			Satellite: BatchSatellite{Name: "Sat2", NORADID: "12346"},
			Success:   false,
			Error:     fmt.Errorf("error 2"),
		},
	}

	comparison := CompareSatellites(results)

	if comparison.Summary.Successful != 0 {
		t.Errorf("Successful = %d, want %d", comparison.Summary.Successful, 0)
	}
	if comparison.Summary.Failed != 2 {
		t.Errorf("Failed = %d, want %d", comparison.Summary.Failed, 2)
	}
}

func TestExportBatchTLECSV(t *testing.T) {
	results := []BatchTLEResult{
		{
			Satellite: BatchSatellite{Name: "Test Sat", NORADID: "12345"},
			Success:   true,
			TLE: TLE{
				CommonName:              "Test Satellite",
				SatelliteCatalogNumber:  12345,
				OrbitInclination:        51.6,
				MeanMotion:             15.5,
				Eccentrcity:             0.0001,
			},
		},
		{
			Satellite: BatchSatellite{Name: "Failed Sat", NORADID: "12346"},
			Success:   false,
			Error:     fmt.Errorf("test error"),
		},
	}

	tempFile := filepath.Join(t.TempDir(), "batch_tle.csv")
	err := exportBatchTLECSV(results, tempFile)
	if err != nil {
		t.Fatalf("exportBatchTLECSV() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatal("CSV file was not created")
	}

	// Read and verify content
	data, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Test Sat") {
		t.Error("CSV should contain satellite name")
	}
	if !strings.Contains(content, "12345") {
		t.Error("CSV should contain NORAD ID")
	}
	if !strings.Contains(content, "Success") {
		t.Error("CSV should contain success status")
	}
}

func TestExportBatchTLEJSON(t *testing.T) {
	results := []BatchTLEResult{
		{
			Satellite: BatchSatellite{Name: "Test Sat", NORADID: "12345"},
			Success:   true,
			TLE: TLE{
				CommonName:              "Test Satellite",
				SatelliteCatalogNumber:  12345,
				OrbitInclination:        51.6,
			},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "batch_tle.json")
	err := exportBatchTLEJSON(results, tempFile)
	if err != nil {
		t.Fatalf("exportBatchTLEJSON() failed: %v", err)
	}

	// Verify file exists and is valid JSON
	data, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["batch_results"] == nil {
		t.Error("JSON should contain batch_results")
	}
	if result["total_count"] != float64(1) {
		t.Errorf("total_count = %v, want %d", result["total_count"], 1)
	}
}

func TestExportBatchTLEText(t *testing.T) {
	results := []BatchTLEResult{
		{
			Satellite: BatchSatellite{Name: "Test Sat", NORADID: "12345"},
			Success:   true,
			TLE: TLE{
				CommonName:              "Test Satellite",
				SatelliteCatalogNumber:  12345,
				OrbitInclination:        51.6,
				MeanMotion:             15.5,
			},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "batch_tle.txt")
	err := exportBatchTLEText(results, tempFile)
	if err != nil {
		t.Fatalf("exportBatchTLEText() failed: %v", err)
	}

	// Verify file exists
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read text file: %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "Test Sat") {
		t.Error("Text file should contain satellite name")
	}
	if !strings.Contains(text, "12345") {
		t.Error("Text file should contain NORAD ID")
	}
	if !strings.Contains(text, "Success") {
		t.Error("Text file should contain success status")
	}
}

func TestExportBatchComparisonCSV(t *testing.T) {
	comparison := BatchComparisonResult{
		Summary: BatchSummary{
			TotalProcessed: 2,
			Successful:     1,
			Failed:         1,
			AverageInclination: 51.6,
			AverageMeanMotion:  15.5,
		},
		Results: []BatchTLEResult{
			{
				Satellite: BatchSatellite{Name: "Sat1", NORADID: "12345"},
				Success:   true,
				TLE: TLE{
					OrbitInclination: 51.6,
					MeanMotion:       15.5,
					Eccentrcity:      0.0001,
				},
			},
			{
				Satellite: BatchSatellite{Name: "Sat2", NORADID: "12346"},
				Success:   false,
				Error:     fmt.Errorf("test error"),
			},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "comparison.csv")
	err := exportBatchComparisonCSV(comparison, tempFile)
	if err != nil {
		t.Fatalf("exportBatchComparisonCSV() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatal("CSV file was not created")
	}

	// Read and verify
	data, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Total Processed") {
		t.Error("CSV should contain summary")
	}
	if !strings.Contains(content, "Sat1") {
		t.Error("CSV should contain satellite names")
	}
}

func TestExportBatchComparisonJSON(t *testing.T) {
	comparison := BatchComparisonResult{
		Summary: BatchSummary{
			TotalProcessed: 1,
			Successful:     1,
		},
		Results: []BatchTLEResult{
			{
				Satellite: BatchSatellite{Name: "Sat1", NORADID: "12345"},
				Success:   true,
				TLE: TLE{OrbitInclination: 51.6},
			},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "comparison.json")
	err := exportBatchComparisonJSON(comparison, tempFile)
	if err != nil {
		t.Fatalf("exportBatchComparisonJSON() failed: %v", err)
	}

	// Verify file exists and is valid JSON
	data, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["summary"] == nil {
		t.Error("JSON should contain summary")
	}
	if result["results"] == nil {
		t.Error("JSON should contain results")
	}
}

func TestExportBatchComparisonText(t *testing.T) {
	comparison := BatchComparisonResult{
		Summary: BatchSummary{
			TotalProcessed: 1,
			Successful:     1,
			AverageInclination: 51.6,
		},
		Results: []BatchTLEResult{
			{
				Satellite: BatchSatellite{Name: "Sat1", NORADID: "12345"},
				Success:   true,
				TLE: TLE{OrbitInclination: 51.6},
			},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "comparison.txt")
	err := exportBatchComparisonText(comparison, tempFile)
	if err != nil {
		t.Fatalf("exportBatchComparisonText() failed: %v", err)
	}

	// Verify file exists
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read text file: %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "Sat1") {
		t.Error("Text file should contain satellite name")
	}
	if !strings.Contains(text, "Summary") {
		t.Error("Text file should contain summary section")
	}
}

func TestBatchTLEResultStruct(t *testing.T) {
	result := BatchTLEResult{
		Satellite: BatchSatellite{
			Name:     "Test",
			NORADID:  "12345",
			Country:  "US",
			ObjectType: "PAYLOAD",
		},
		Success: true,
		TLE: TLE{
			SatelliteCatalogNumber: 12345,
		},
	}

	if !result.Success {
		t.Error("Success should be true")
	}
	if result.Satellite.NORADID != "12345" {
		t.Errorf("NORADID = %q, want %q", result.Satellite.NORADID, "12345")
	}
}

func TestBatchSummaryStruct(t *testing.T) {
	summary := BatchSummary{
		TotalProcessed:       10,
		Successful:           8,
		Failed:              2,
		AverageInclination:  51.6,
		AverageMeanMotion:   15.5,
		LowestAltitude:      400.0,
		HighestAltitude:     500.0,
	}

	if summary.TotalProcessed != 10 {
		t.Errorf("TotalProcessed = %d, want %d", summary.TotalProcessed, 10)
	}
	if summary.Successful != 8 {
		t.Errorf("Successful = %d, want %d", summary.Successful, 8)
	}
	if summary.AverageInclination != 51.6 {
		t.Errorf("AverageInclination = %.2f, want %.2f", summary.AverageInclination, 51.6)
	}
}

// Benchmark tests
func BenchmarkCompareSatellites(b *testing.B) {
	results := []BatchTLEResult{
		{
			Satellite: BatchSatellite{Name: "Sat1", NORADID: "12345"},
			Success:   true,
			TLE: TLE{
				OrbitInclination: 51.6,
				MeanMotion:       15.5,
				Eccentrcity:      0.0001,
			},
		},
		{
			Satellite: BatchSatellite{Name: "Sat2", NORADID: "12346"},
			Success:   true,
			TLE: TLE{
				OrbitInclination: 98.2,
				MeanMotion:       14.2,
				Eccentrcity:      0.0002,
			},
		},
	}

	for i := 0; i < b.N; i++ {
		CompareSatellites(results)
	}
}

func BenchmarkExportBatchTLEJSON(b *testing.B) {
	results := []BatchTLEResult{
		{
			Satellite: BatchSatellite{Name: "Test", NORADID: "12345"},
			Success:   true,
			TLE: TLE{
				CommonName:              "Test Satellite",
				SatelliteCatalogNumber:  12345,
				OrbitInclination:       51.6,
				MeanMotion:             15.5,
			},
		},
	}

	tempDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tempDir, fmt.Sprintf("test_%d.json", i))
		exportBatchTLEJSON(results, file)
	}
}

