package osint

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportTLECSV(t *testing.T) {
	tle := TLE{
		CommonName:                 "ISS (ZARYA)",
		SatelliteCatalogNumber:     25544,
		ElsetClassificiation:       "U",
		InternationalDesignator:    "1998-067A",
		ElementSetEpoch:            24001.5,
		FirstDerivativeMeanMotion:  0.0001,
		SecondDerivativeMeanMotion: "00000+0",
		BDragTerm:                  "00000+0",
		ElementSetType:             0,
		ElementNumber:              999,
		ChecksumOne:                5,
		OrbitInclination:           51.6442,
		RightAscension:             123.4567,
		Eccentrcity:                0.0001234,
		Perigee:                    234.5678,
		MeanAnamoly:                345.6789,
		MeanMotion:                 15.49,
		RevolutionNumber:           12345,
		ChecksumTwo:                6,
	}

	tempFile := filepath.Join(t.TempDir(), "test_tle.csv")
	if err := exportTLECSV(tle, tempFile); err != nil {
		t.Fatalf("exportTLECSV() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatal("CSV file was not created")
	}

	// Read and verify CSV content
	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	if len(records) < 2 {
		t.Fatal("CSV should have at least header and one data row")
	}

	// Check header
	if records[0][0] != "Field" || records[0][1] != "Value" {
		t.Errorf("CSV header incorrect: got %v", records[0])
	}

	// Check some data rows
	found := false
	for _, record := range records[1:] {
		if record[0] == "Common Name" && record[1] == "ISS (ZARYA)" {
			found = true
			break
		}
	}
	if !found {
		t.Error("CSV data not found correctly")
	}
}

func TestExportTLEJSON(t *testing.T) {
	tle := TLE{
		CommonName:              "Test Satellite",
		SatelliteCatalogNumber:  12345,
		ElementSetEpoch:        24001.5,
		OrbitInclination:       51.6442,
		MeanMotion:             15.49,
	}

	tempFile := filepath.Join(t.TempDir(), "test_tle.json")
	if err := exportTLEJSON(tle, tempFile); err != nil {
		t.Fatalf("exportTLEJSON() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatal("JSON file was not created")
	}

	// Read and verify JSON content
	data, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["common_name"] != "Test Satellite" {
		t.Errorf("JSON common_name = %v, want %v", result["common_name"], "Test Satellite")
	}

	if result["satellite_catalog_number"] != float64(12345) {
		t.Errorf("JSON satellite_catalog_number = %v, want %v", result["satellite_catalog_number"], 12345)
	}

	if result["export_timestamp"] == nil {
		t.Error("JSON should include export_timestamp")
	}
}

func TestExportTLEText(t *testing.T) {
	tle := TLE{
		CommonName:              "Test Satellite",
		SatelliteCatalogNumber:  12345,
		ElementSetEpoch:        24001.5,
		OrbitInclination:       51.6442,
	}

	tempFile := filepath.Join(t.TempDir(), "test_tle.txt")
	if err := exportTLEText(tle, tempFile); err != nil {
		t.Fatalf("exportTLEText() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatal("Text file was not created")
	}

	// Read and verify content
	data, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read text file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Test Satellite") {
		t.Error("Text file should contain satellite name")
	}
	if !strings.Contains(content, "12345") {
		t.Error("Text file should contain catalog number")
	}
	if !strings.Contains(content, "Exported:") {
		t.Error("Text file should contain export timestamp")
	}
}

func TestExportVisualPredictionCSV(t *testing.T) {
	data := VisualPassesResponse{
		Info: Info{
			SatName:           "ISS (ZARYA)",
			SatID:             25544,
			TransactionsCount: 1,
			PassesCount:       2,
		},
		Passes: []Pass{
			{
				StartAz:        45.0,
				StartAzCompass: "NE",
				StartEl:        10.0,
				StartUTC:       1234567890,
				MaxAz:          90.0,
				MaxAzCompass:   "E",
				MaxEl:          60.0,
				MaxUTC:         1234567900,
				EndAz:          135.0,
				EndAzCompass:   "SE",
				EndEl:          10.0,
				EndUTC:         1234567910,
				Mag:            2.5,
				Duration:       600,
			},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "test_visual.csv")
	if err := exportVisualPredictionCSV(data, tempFile); err != nil {
		t.Fatalf("exportVisualPredictionCSV() failed: %v", err)
	}

	// Verify file exists and can be read
	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	// Filter out empty rows
	var nonEmptyRecords [][]string
	for _, record := range records {
		if len(record) > 0 && (len(record) == 1 && record[0] == "" || len(record) > 1) {
			if len(record) > 1 || record[0] != "" {
				nonEmptyRecords = append(nonEmptyRecords, record)
			}
		}
	}

	if len(nonEmptyRecords) < 2 {
		t.Fatalf("CSV should have info and pass data, got %d non-empty rows", len(nonEmptyRecords))
	}

	// Check satellite info
	if records[0][0] != "Satellite Name" {
		t.Errorf("First row should be satellite info header, got %v", records[0])
	}
}

func TestExportVisualPredictionJSON(t *testing.T) {
	data := VisualPassesResponse{
		Info: Info{
			SatName: "Test Satellite",
			SatID:   12345,
		},
		Passes: []Pass{
			{StartAz: 45.0, MaxEl: 60.0, Duration: 600},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "test_visual.json")
	if err := exportVisualPredictionJSON(data, tempFile); err != nil {
		t.Fatalf("exportVisualPredictionJSON() failed: %v", err)
	}

	fileData, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(fileData, &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["satellite_info"] == nil {
		t.Error("JSON should contain satellite_info")
	}

	if result["passes"] == nil {
		t.Error("JSON should contain passes")
	}
}

func TestExportRadioPredictionCSV(t *testing.T) {
	data := RadioPassResponse{
		Info: Info{
			SatName:           "Test Satellite",
			SatID:             12345,
			TransactionsCount: 1,
			PassesCount:       1,
		},
		Passes: []RadioPass{
			{
				StartAz:        45.0,
				StartAzCompass: "NE",
				StartUTC:       1234567890,
				MaxAz:          90.0,
				MaxAzCompass:   "E",
				MaxEl:          60.0,
				MaxUTC:         1234567900,
				EndAz:          135.0,
				EndAzCompass:   "SE",
				EndUTC:         1234567910,
			},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "test_radio.csv")
	if err := exportRadioPredictionCSV(data, tempFile); err != nil {
		t.Fatalf("exportRadioPredictionCSV() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatal("CSV file was not created")
	}
}

func TestExportSatellitePositionCSV(t *testing.T) {
	data := Response{
		SatelliteInfo: SatelliteInfo{
			Satname: "Test Satellite",
			Satid:   12345,
		},
		Positions: []Position{
			{
				Satlatitude:  40.7128,
				Satlongitude: -74.0060,
				Sataltitude:  400.0,
				Azimuth:      45.0,
				Dec:          30.0,
				Timestamp:    1234567890,
			},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "test_position.csv")
	if err := exportSatellitePositionCSV(data, tempFile); err != nil {
		t.Fatalf("exportSatellitePositionCSV() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatal("CSV file was not created")
	}

	// Read and verify
	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	// Filter out empty rows
	var nonEmptyRecords [][]string
	for _, record := range records {
		if len(record) > 0 && (len(record) == 1 && record[0] == "" || len(record) > 1) {
			if len(record) > 1 || record[0] != "" {
				nonEmptyRecords = append(nonEmptyRecords, record)
			}
		}
	}

	if len(nonEmptyRecords) < 2 {
		t.Fatalf("CSV should have info and position data, got %d non-empty rows", len(nonEmptyRecords))
	}
}

func TestExportSatellitePositionJSON(t *testing.T) {
	data := Response{
		SatelliteInfo: SatelliteInfo{
			Satname: "Test Satellite",
			Satid:   12345,
		},
		Positions: []Position{
			{
				Satlatitude:  40.7128,
				Satlongitude: -74.0060,
				Sataltitude:  400.0,
			},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "test_position.json")
	if err := exportSatellitePositionJSON(data, tempFile); err != nil {
		t.Fatalf("exportSatellitePositionJSON() failed: %v", err)
	}

	// Verify file exists and is valid JSON
	fileData, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(fileData, &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["satellite_info"] == nil {
		t.Error("JSON should contain satellite_info")
	}

	if result["positions"] == nil {
		t.Error("JSON should contain positions")
	}
}

func TestExportSatellitePositionText(t *testing.T) {
	data := Response{
		SatelliteInfo: SatelliteInfo{
			Satname: "Test Satellite",
			Satid:   12345,
		},
		Positions: []Position{
			{
				Satlatitude:  40.7128,
				Satlongitude: -74.0060,
				Sataltitude:  400.0,
			},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "test_position.txt")
	if err := exportSatellitePositionText(data, tempFile); err != nil {
		t.Fatalf("exportSatellitePositionText() failed: %v", err)
	}

	// Verify file exists
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read text file: %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "Test Satellite") {
		t.Error("Text file should contain satellite name")
	}
	if !strings.Contains(text, "40.7128") {
		t.Error("Text file should contain latitude")
	}
}

func TestExportTLE(t *testing.T) {
	tle := TLE{
		CommonName:              "Test",
		SatelliteCatalogNumber: 12345,
	}

	tempDir := t.TempDir()

	// Test CSV export
	csvFile := filepath.Join(tempDir, "test.csv")
	if err := ExportTLE(tle, FormatCSV, csvFile); err != nil {
		t.Fatalf("ExportTLE(CSV) failed: %v", err)
	}
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		t.Error("CSV file was not created")
	}

	// Test JSON export
	jsonFile := filepath.Join(tempDir, "test.json")
	if err := ExportTLE(tle, FormatJSON, jsonFile); err != nil {
		t.Fatalf("ExportTLE(JSON) failed: %v", err)
	}
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		t.Error("JSON file was not created")
	}

	// Test Text export
	txtFile := filepath.Join(tempDir, "test.txt")
	if err := ExportTLE(tle, FormatText, txtFile); err != nil {
		t.Fatalf("ExportTLE(Text) failed: %v", err)
	}
	if _, err := os.Stat(txtFile); os.IsNotExist(err) {
		t.Error("Text file was not created")
	}

	// Test invalid format
	invalidFile := filepath.Join(tempDir, "test.invalid")
	if err := ExportTLE(tle, ExportFormat("INVALID"), invalidFile); err == nil {
		t.Error("ExportTLE should fail for invalid format")
	}
}

func TestExportFormatConstants(t *testing.T) {
	if FormatCSV != "CSV" {
		t.Errorf("FormatCSV = %q, want %q", FormatCSV, "CSV")
	}
	if FormatJSON != "JSON" {
		t.Errorf("FormatJSON = %q, want %q", FormatJSON, "JSON")
	}
	if FormatText != "Text" {
		t.Errorf("FormatText = %q, want %q", FormatText, "Text")
	}
}

// Benchmark tests
func BenchmarkExportTLECSV(b *testing.B) {
	tle := TLE{
		CommonName:              "Test Satellite",
		SatelliteCatalogNumber:  12345,
		ElementSetEpoch:        24001.5,
		OrbitInclination:       51.6442,
		MeanMotion:             15.49,
	}

	tempDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tempDir, fmt.Sprintf("test_%d.csv", i))
		exportTLECSV(tle, file)
	}
}

func BenchmarkExportTLEJSON(b *testing.B) {
	tle := TLE{
		CommonName:              "Test Satellite",
		SatelliteCatalogNumber:  12345,
		ElementSetEpoch:        24001.5,
		OrbitInclination:       51.6442,
		MeanMotion:             15.49,
	}

	tempDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tempDir, fmt.Sprintf("test_%d.json", i))
		exportTLEJSON(tle, file)
	}
}

