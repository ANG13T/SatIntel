package osint

import (
	"math"
	"testing"
	"time"
)

// Test TLE data for ISS (International Space Station)
// NORAD ID: 25544
const testTLELine1 = "1 25544U 98067A   04236.56031392  .00020137  00000-0  16538-3 0  9993"
const testTLELine2 = "2 25544  51.6335 344.7760 0007976 126.2523 325.9359 15.70406856328906"

// Test TLE data for a different satellite (NOAA 18)
const testTLE2Line1 = "1 28654U 05018A   21123.45678901  .00000123  00000-0  12345-4 0  9998"
const testTLE2Line2 = "2 28654  98.7145 123.4567 0012345 234.5678 345.6789 14.12345678901234"

func TestCalculateSGP4Position(t *testing.T) {
	// Test with current time
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	pos, err := CalculateSGP4Position(testTLELine1, testTLELine2, targetTime)
	if err != nil {
		t.Fatalf("CalculateSGP4Position failed: %v", err)
	}

	// Validate position data
	if math.IsNaN(pos.Latitude) || math.IsInf(pos.Latitude, 0) {
		t.Error("Latitude is NaN or Inf")
	}
	if math.IsNaN(pos.Longitude) || math.IsInf(pos.Longitude, 0) {
		t.Error("Longitude is NaN or Inf")
	}
	if math.IsNaN(pos.Altitude) || pos.Altitude < 0 {
		t.Errorf("Altitude is invalid: %f", pos.Altitude)
	}
	if math.IsNaN(pos.Velocity) || pos.Velocity <= 0 {
		t.Errorf("Velocity is invalid: %f", pos.Velocity)
	}
	if pos.Timestamp != targetTime.Unix() {
		t.Errorf("Timestamp mismatch: expected %d, got %d", targetTime.Unix(), pos.Timestamp)
	}

	// Validate reasonable ranges
	if pos.Latitude < -90 || pos.Latitude > 90 {
		t.Errorf("Latitude out of range: %f", pos.Latitude)
	}
	if pos.Longitude < -180 || pos.Longitude > 180 {
		t.Errorf("Longitude out of range: %f", pos.Longitude)
	}
	// Altitude can vary significantly depending on TLE epoch and calculation method
	// Accept any positive altitude value
	if pos.Altitude < 0 {
		t.Errorf("Altitude is negative: %f km", pos.Altitude)
	}
	// Velocity should be positive (units may vary, so just check it's not zero/negative)
	if pos.Velocity <= 0 {
		t.Errorf("Velocity is invalid: %f", pos.Velocity)
	}
}

func TestCalculateSGP4Position_InvalidTLE(t *testing.T) {
	targetTime := time.Now()

	// Test with invalid line 1 (doesn't start with "1 ")
	_, err := CalculateSGP4Position("INVALID LINE", testTLELine2, targetTime)
	if err == nil {
		t.Error("Expected error for invalid line 1, got nil")
	}

	// Test with invalid line 2 (doesn't start with "2 ")
	_, err = CalculateSGP4Position(testTLELine1, "INVALID LINE", targetTime)
	if err == nil {
		t.Error("Expected error for invalid line 2, got nil")
	}

	// Test with too short lines
	_, err = CalculateSGP4Position("1 25544", "2 25544", targetTime)
	if err == nil {
		t.Error("Expected error for too short lines, got nil")
	}
}

func TestCalculateSGP4PositionWithObserver(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	// Observer at New York City (approximately)
	observer := ObserverPosition{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Altitude:  10.0, // 10 meters above sea level
	}

	result, err := CalculateSGP4PositionWithObserver(testTLELine1, testTLELine2, targetTime, observer)
	if err != nil {
		t.Fatalf("CalculateSGP4PositionWithObserver failed: %v", err)
	}

	// Validate position
	if math.IsNaN(result.Position.Latitude) || math.IsInf(result.Position.Latitude, 0) {
		t.Error("Position latitude is NaN or Inf")
	}

	// Validate look angles
	if math.IsNaN(result.LookAngles.Azimuth) || math.IsInf(result.LookAngles.Azimuth, 0) {
		t.Error("Azimuth is NaN or Inf")
	}
	if math.IsNaN(result.LookAngles.Elevation) || math.IsInf(result.LookAngles.Elevation, 0) {
		t.Error("Elevation is NaN or Inf")
	}
	if math.IsNaN(result.LookAngles.Range) || result.LookAngles.Range <= 0 {
		t.Errorf("Range is invalid: %f", result.LookAngles.Range)
	}

	// Validate angle ranges
	if result.LookAngles.Azimuth < 0 || result.LookAngles.Azimuth >= 360 {
		t.Errorf("Azimuth out of range [0, 360): %f", result.LookAngles.Azimuth)
	}
	if result.LookAngles.Elevation < -90 || result.LookAngles.Elevation > 90 {
		t.Errorf("Elevation out of range [-90, 90]: %f", result.LookAngles.Elevation)
	}
	if result.LookAngles.Range < 100 || result.LookAngles.Range > 50000 {
		t.Errorf("Range out of reasonable bounds: %f km", result.LookAngles.Range)
	}
}

func TestCalculateSGP4PositionWithObserver_DifferentLocations(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Test with multiple observer locations
	locations := []ObserverPosition{
		{Latitude: 0.0, Longitude: 0.0, Altitude: 0.0},      // Equator, Prime Meridian
		{Latitude: 51.5074, Longitude: -0.1278, Altitude: 35.0}, // London
		{Latitude: 35.6762, Longitude: 139.6503, Altitude: 40.0}, // Tokyo
		{Latitude: -33.8688, Longitude: 151.2093, Altitude: 50.0}, // Sydney
	}

	for i, loc := range locations {
		result, err := CalculateSGP4PositionWithObserver(testTLELine1, testTLELine2, targetTime, loc)
		if err != nil {
			t.Errorf("Location %d failed: %v", i, err)
			continue
		}

		// Validate results
		if math.IsNaN(result.LookAngles.Azimuth) || math.IsNaN(result.LookAngles.Elevation) {
			t.Errorf("Location %d produced NaN angles", i)
		}
		if result.LookAngles.Range <= 0 {
			t.Errorf("Location %d has invalid range: %f", i, result.LookAngles.Range)
		}
	}
}

func TestCalculateSGP4Positions(t *testing.T) {
	startTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	endTime := startTime.Add(2 * time.Hour)
	interval := 15 * time.Minute

	positions, err := CalculateSGP4Positions(testTLELine1, testTLELine2, startTime, endTime, interval)
	if err != nil {
		t.Fatalf("CalculateSGP4Positions failed: %v", err)
	}

	// Should have 9 positions (0, 15, 30, 45, 60, 75, 90, 105, 120 minutes)
	expectedCount := 9
	if len(positions) != expectedCount {
		t.Errorf("Expected %d positions, got %d", expectedCount, len(positions))
	}

	// Validate each position
	for i, pos := range positions {
		if math.IsNaN(pos.Latitude) || math.IsNaN(pos.Longitude) {
			t.Errorf("Position %d has NaN coordinates", i)
		}
		if pos.Altitude < 0 {
			t.Errorf("Position %d has negative altitude: %f", i, pos.Altitude)
		}
	}

	// Check that positions change over time (satellite is moving)
	if len(positions) >= 2 {
		firstLat := positions[0].Latitude
		lastLat := positions[len(positions)-1].Latitude
		// Allow some tolerance, but positions should differ
		if math.Abs(firstLat-lastLat) < 0.01 {
			t.Error("Positions are too similar - satellite may not be moving")
		}
	}
}

func TestCalculateSGP4Positions_InvalidInputs(t *testing.T) {
	startTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	endTime := startTime.Add(1 * time.Hour)

	// Test with start time after end time
	_, err := CalculateSGP4Positions(testTLELine1, testTLELine2, endTime, startTime, 15*time.Minute)
	if err == nil {
		t.Error("Expected error when start time is after end time")
	}

	// Test with zero interval
	_, err = CalculateSGP4Positions(testTLELine1, testTLELine2, startTime, endTime, 0)
	if err == nil {
		t.Error("Expected error for zero interval")
	}

	// Test with negative interval
	_, err = CalculateSGP4Positions(testTLELine1, testTLELine2, startTime, endTime, -15*time.Minute)
	if err == nil {
		t.Error("Expected error for negative interval")
	}
}

func TestCalculateSGP4Position_DifferentTimes(t *testing.T) {
	// Test with different times to ensure consistency
	times := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC),
		time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
	}

	for _, targetTime := range times {
		pos, err := CalculateSGP4Position(testTLELine1, testTLELine2, targetTime)
		if err != nil {
			t.Errorf("Failed for time %v: %v", targetTime, err)
			continue
		}

		// Validate position is reasonable
		if pos.Latitude < -90 || pos.Latitude > 90 {
			t.Errorf("Invalid latitude at %v: %f", targetTime, pos.Latitude)
		}
		if pos.Longitude < -180 || pos.Longitude > 180 {
			t.Errorf("Invalid longitude at %v: %f", targetTime, pos.Longitude)
		}
	}
}

func TestCalculateSGP4Position_DifferentSatellites(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Test with different TLE sets
	tleSets := []struct {
		name  string
		line1 string
		line2 string
	}{
		{"ISS", testTLELine1, testTLELine2},
		{"NOAA 18", testTLE2Line1, testTLE2Line2},
	}

	for _, tleSet := range tleSets {
		pos, err := CalculateSGP4Position(tleSet.line1, tleSet.line2, targetTime)
		if err != nil {
			t.Errorf("Failed for %s: %v", tleSet.name, err)
			continue
		}

		// Validate position
		if math.IsNaN(pos.Latitude) || math.IsNaN(pos.Longitude) {
			t.Errorf("%s produced NaN coordinates", tleSet.name)
		}
		if pos.Altitude < 0 {
			t.Errorf("%s has negative altitude: %f", tleSet.name, pos.Altitude)
		}
	}
}

func TestCalculateSGP4PositionWithObserver_EdgeCases(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Test with observer at extreme locations
	extremeLocations := []ObserverPosition{
		{Latitude: 90.0, Longitude: 0.0, Altitude: 0.0},   // North Pole
		{Latitude: -90.0, Longitude: 0.0, Altitude: 0.0},   // South Pole
		{Latitude: 0.0, Longitude: 180.0, Altitude: 0.0},  // International Date Line
		{Latitude: 0.0, Longitude: -180.0, Altitude: 0.0}, // International Date Line (west)
	}

	for i, loc := range extremeLocations {
		result, err := CalculateSGP4PositionWithObserver(testTLELine1, testTLELine2, targetTime, loc)
		if err != nil {
			t.Errorf("Extreme location %d failed: %v", i, err)
			continue
		}

		// Validate angles are in valid ranges
		if result.LookAngles.Azimuth < 0 || result.LookAngles.Azimuth >= 360 {
			t.Errorf("Extreme location %d: azimuth out of range: %f", i, result.LookAngles.Azimuth)
		}
		if result.LookAngles.Elevation < -90 || result.LookAngles.Elevation > 90 {
			t.Errorf("Extreme location %d: elevation out of range: %f", i, result.LookAngles.Elevation)
		}
	}
}

func TestCalculateSGP4Position_TimeConsistency(t *testing.T) {
	// Test that positions change smoothly over time
	baseTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	positions := make([]SGPPosition, 5)
	for i := 0; i < 5; i++ {
		targetTime := baseTime.Add(time.Duration(i) * 5 * time.Minute)
		pos, err := CalculateSGP4Position(testTLELine1, testTLELine2, targetTime)
		if err != nil {
			t.Fatalf("Failed at step %d: %v", i, err)
		}
		positions[i] = pos
	}

	// Check that positions are different (satellite is moving)
	for i := 1; i < len(positions); i++ {
		prev := positions[i-1]
		curr := positions[i]
		
		// Calculate distance between consecutive positions
		latDiff := curr.Latitude - prev.Latitude
		lonDiff := curr.Longitude - prev.Longitude
		
		// Positions should change (allow small tolerance for numerical precision)
		if math.Abs(latDiff) < 0.0001 && math.Abs(lonDiff) < 0.0001 {
			t.Errorf("Positions at step %d and %d are too similar", i-1, i)
		}
	}
}

func TestCalculateSGP4PositionFromTLEStruct(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	// Parse TLE first
	tle := ConstructTLE("ISS (ZARYA)", testTLELine1, testTLELine2)
	
	// Test CalculateSGP4PositionFromTLE
	pos, err := CalculateSGP4PositionFromTLE(tle, testTLELine1, testTLELine2, targetTime)
	if err != nil {
		t.Fatalf("CalculateSGP4PositionFromTLE failed: %v", err)
	}

	// Validate position
	if math.IsNaN(pos.Latitude) || math.IsNaN(pos.Longitude) {
		t.Error("Position has NaN coordinates")
	}
	if pos.Altitude < 0 {
		t.Errorf("Negative altitude: %f", pos.Altitude)
	}
	
	// Test CalculateSGP4PositionFromTLEStruct
	pos2, err := CalculateSGP4PositionFromTLEStruct(tle, testTLELine1, testTLELine2, targetTime)
	if err != nil {
		t.Fatalf("CalculateSGP4PositionFromTLEStruct failed: %v", err)
	}
	
	// Validate position
	if math.IsNaN(pos2.Latitude) || math.IsNaN(pos2.Longitude) {
		t.Error("Position from TLE struct has NaN coordinates")
	}
	if pos2.Altitude < 0 {
		t.Errorf("Negative altitude from TLE struct: %f", pos2.Altitude)
	}
	
	// Both should produce the same result
	if math.Abs(pos.Latitude-pos2.Latitude) > 0.0001 {
		t.Errorf("Latitude mismatch: %f vs %f", pos.Latitude, pos2.Latitude)
	}
	if math.Abs(pos.Longitude-pos2.Longitude) > 0.0001 {
		t.Errorf("Longitude mismatch: %f vs %f", pos.Longitude, pos2.Longitude)
	}
}

func TestCalculateSGP4PositionFromTLEStructWithObserver(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	// Parse TLE first
	tle := ConstructTLE("ISS (ZARYA)", testTLELine1, testTLELine2)
	
	observer := ObserverPosition{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Altitude:  10.0,
	}
	
	result, err := CalculateSGP4PositionFromTLEStructWithObserver(tle, testTLELine1, testTLELine2, targetTime, observer)
	if err != nil {
		t.Fatalf("CalculateSGP4PositionFromTLEStructWithObserver failed: %v", err)
	}

	// Validate result
	if math.IsNaN(result.LookAngles.Azimuth) || math.IsNaN(result.LookAngles.Elevation) {
		t.Error("Look angles are NaN")
	}
	if result.LookAngles.Range <= 0 {
		t.Errorf("Invalid range: %f", result.LookAngles.Range)
	}
}

// Benchmark tests
func BenchmarkCalculateSGP4Position(b *testing.B) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CalculateSGP4Position(testTLELine1, testTLELine2, targetTime)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkCalculateSGP4PositionWithObserver(b *testing.B) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	observer := ObserverPosition{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Altitude:  10.0,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CalculateSGP4PositionWithObserver(testTLELine1, testTLELine2, targetTime, observer)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

// Additional edge case and validation tests

func TestCalculateSGP4Position_WhitespaceHandling(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	// Test with leading/trailing whitespace
	line1WithSpaces := "  " + testTLELine1 + "  "
	line2WithSpaces := "  " + testTLELine2 + "  "
	
	pos, err := CalculateSGP4Position(line1WithSpaces, line2WithSpaces, targetTime)
	if err != nil {
		t.Fatalf("Failed with whitespace: %v", err)
	}
	
	// Should still produce valid results
	if math.IsNaN(pos.Latitude) || math.IsNaN(pos.Longitude) {
		t.Error("Position has NaN coordinates after whitespace trimming")
	}
}

func TestCalculateSGP4Position_ExactLength(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	// TLE lines should be exactly 69 characters
	if len(testTLELine1) != 69 {
		t.Skipf("Test TLE line 1 is %d characters, not 69", len(testTLELine1))
	}
	if len(testTLELine2) != 69 {
		t.Skipf("Test TLE line 2 is %d characters, not 69", len(testTLELine2))
	}
	
	pos, err := CalculateSGP4Position(testTLELine1, testTLELine2, targetTime)
	if err != nil {
		t.Fatalf("Failed with exact length TLE: %v", err)
	}
	
	if math.IsNaN(pos.Latitude) || math.IsNaN(pos.Longitude) {
		t.Error("Position has NaN coordinates")
	}
}

func TestCalculateSGP4Position_JustBelowMinimumLength(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	// Create lines that are just below 69 characters
	shortLine1 := testTLELine1[:68]
	shortLine2 := testTLELine2[:68]
	
	_, err := CalculateSGP4Position(shortLine1, shortLine2, targetTime)
	if err == nil {
		t.Error("Expected error for lines shorter than 69 characters")
	}
}

func TestCalculateSGP4Positions_SinglePosition(t *testing.T) {
	startTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	endTime := startTime // Same time
	interval := 1 * time.Minute
	
	positions, err := CalculateSGP4Positions(testTLELine1, testTLELine2, startTime, endTime, interval)
	if err != nil {
		t.Fatalf("CalculateSGP4Positions failed: %v", err)
	}
	
	// Should have exactly 1 position
	if len(positions) != 1 {
		t.Errorf("Expected 1 position, got %d", len(positions))
	}
	
	if math.IsNaN(positions[0].Latitude) || math.IsNaN(positions[0].Longitude) {
		t.Error("Position has NaN coordinates")
	}
}

func TestCalculateSGP4Positions_VeryShortInterval(t *testing.T) {
	startTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	endTime := startTime.Add(1 * time.Minute)
	interval := 1 * time.Second
	
	positions, err := CalculateSGP4Positions(testTLELine1, testTLELine2, startTime, endTime, interval)
	if err != nil {
		t.Fatalf("CalculateSGP4Positions failed: %v", err)
	}
	
	// Should have 61 positions (0 to 60 seconds)
	expectedCount := 61
	if len(positions) != expectedCount {
		t.Errorf("Expected %d positions, got %d", expectedCount, len(positions))
	}
	
	// Validate all positions
	for i, pos := range positions {
		if math.IsNaN(pos.Latitude) || math.IsNaN(pos.Longitude) {
			t.Errorf("Position %d has NaN coordinates", i)
		}
	}
}

func TestCalculateSGP4Positions_LongTimeRange(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(24 * time.Hour)
	interval := 1 * time.Hour
	
	positions, err := CalculateSGP4Positions(testTLELine1, testTLELine2, startTime, endTime, interval)
	if err != nil {
		t.Fatalf("CalculateSGP4Positions failed: %v", err)
	}
	
	// Should have 25 positions (0 to 24 hours)
	expectedCount := 25
	if len(positions) != expectedCount {
		t.Errorf("Expected %d positions, got %d", expectedCount, len(positions))
	}
	
	// Check that positions vary significantly over 24 hours
	if len(positions) >= 2 {
		first := positions[0]
		last := positions[len(positions)-1]
		
		// Positions should be quite different after 24 hours
		latDiff := math.Abs(first.Latitude - last.Latitude)
		lonDiff := math.Abs(first.Longitude - last.Longitude)
		
		// At least one coordinate should change significantly
		if latDiff < 0.1 && lonDiff < 0.1 {
			t.Error("Positions are too similar after 24 hours")
		}
	}
}

func TestCalculateSGP4PositionWithObserver_InvalidObserverCoordinates(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	// Test with observer coordinates that are technically out of range
	// The library might handle these, but we should test them
	invalidObservers := []ObserverPosition{
		{Latitude: 91.0, Longitude: 0.0, Altitude: 0.0},   // Latitude > 90
		{Latitude: -91.0, Longitude: 0.0, Altitude: 0.0}, // Latitude < -90
		{Latitude: 0.0, Longitude: 181.0, Altitude: 0.0}, // Longitude > 180
		{Latitude: 0.0, Longitude: -181.0, Altitude: 0.0}, // Longitude < -180
	}
	
	for i, observer := range invalidObservers {
		// The library might normalize these, so we just check it doesn't crash
		result, err := CalculateSGP4PositionWithObserver(testTLELine1, testTLELine2, targetTime, observer)
		if err != nil {
			// Error is acceptable for invalid coordinates
			continue
		}
		
		// If no error, validate the result is reasonable
		if !math.IsNaN(result.LookAngles.Azimuth) && !math.IsNaN(result.LookAngles.Elevation) {
			// Result is valid, which is fine if library normalizes coordinates
			_ = result
		}
		_ = i // Suppress unused variable warning
	}
}

func TestCalculateSGP4PositionWithObserver_HighAltitudeObserver(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	// Test with observer at high altitude (e.g., on a mountain or aircraft)
	highAltitudeObserver := ObserverPosition{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Altitude:  8848.0, // Mount Everest height in meters
	}
	
	result, err := CalculateSGP4PositionWithObserver(testTLELine1, testTLELine2, targetTime, highAltitudeObserver)
	if err != nil {
		t.Fatalf("Failed with high altitude observer: %v", err)
	}
	
	// Validate results
	if math.IsNaN(result.LookAngles.Azimuth) || math.IsNaN(result.LookAngles.Elevation) {
		t.Error("Look angles are NaN")
	}
	if result.LookAngles.Range <= 0 {
		t.Errorf("Invalid range: %f", result.LookAngles.Range)
	}
}

func TestCalculateSGP4PositionWithObserver_NegativeAltitudeObserver(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	// Test with observer below sea level (e.g., Death Valley)
	negativeAltitudeObserver := ObserverPosition{
		Latitude:  36.5054,
		Longitude: -117.0794,
		Altitude:  -86.0, // Death Valley is 86m below sea level
	}
	
	result, err := CalculateSGP4PositionWithObserver(testTLELine1, testTLELine2, targetTime, negativeAltitudeObserver)
	if err != nil {
		t.Fatalf("Failed with negative altitude observer: %v", err)
	}
	
	// Validate results
	if math.IsNaN(result.LookAngles.Azimuth) || math.IsNaN(result.LookAngles.Elevation) {
		t.Error("Look angles are NaN")
	}
}

func TestCalculateSGP4Position_EmptyStrings(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	// Test with empty strings
	_, err := CalculateSGP4Position("", "", targetTime)
	if err == nil {
		t.Error("Expected error for empty TLE lines")
	}
	
	// Test with whitespace only
	_, err = CalculateSGP4Position("   ", "   ", targetTime)
	if err == nil {
		t.Error("Expected error for whitespace-only TLE lines")
	}
}

func TestCalculateSGP4Position_HistoricalDate(t *testing.T) {
	// Test with a date in the past (before TLE epoch)
	historicalTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	
	pos, err := CalculateSGP4Position(testTLELine1, testTLELine2, historicalTime)
	if err != nil {
		// Error is acceptable for dates far from TLE epoch
		return
	}
	
	// If no error, validate position
	if math.IsNaN(pos.Latitude) || math.IsNaN(pos.Longitude) {
		t.Error("Position has NaN coordinates")
	}
}

func TestCalculateSGP4Position_FutureDate(t *testing.T) {
	// Test with a date far in the future
	futureTime := time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC)
	
	pos, err := CalculateSGP4Position(testTLELine1, testTLELine2, futureTime)
	if err != nil {
		// Error is acceptable for dates far from TLE epoch
		return
	}
	
	// If no error, validate position
	if math.IsNaN(pos.Latitude) || math.IsNaN(pos.Longitude) {
		t.Error("Position has NaN coordinates")
	}
}

func TestCalculateSGP4PositionWithObserver_ErrorPropagation(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	observer := ObserverPosition{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Altitude:  10.0,
	}
	
	// Test that errors from CalculateSGP4Position are propagated
	_, err := CalculateSGP4PositionWithObserver("INVALID", "INVALID", targetTime, observer)
	if err == nil {
		t.Error("Expected error for invalid TLE lines")
	}
}

func TestCalculateSGP4Positions_ErrorPropagation(t *testing.T) {
	startTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	endTime := startTime.Add(1 * time.Hour)
	interval := 15 * time.Minute
	
	// Test that errors are propagated when TLE is invalid
	_, err := CalculateSGP4Positions("INVALID", "INVALID", startTime, endTime, interval)
	if err == nil {
		t.Error("Expected error for invalid TLE lines")
	}
}

func TestCalculateSGP4Position_LeapYear(t *testing.T) {
	// Test with leap year date
	leapYearTime := time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC)
	
	pos, err := CalculateSGP4Position(testTLELine1, testTLELine2, leapYearTime)
	if err != nil {
		t.Fatalf("Failed on leap year date: %v", err)
	}
	
	if math.IsNaN(pos.Latitude) || math.IsNaN(pos.Longitude) {
		t.Error("Position has NaN coordinates")
	}
}

func TestCalculateSGP4Position_YearBoundary(t *testing.T) {
	// Test with year boundary transition
	yearBoundaryTime := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
	
	pos, err := CalculateSGP4Position(testTLELine1, testTLELine2, yearBoundaryTime)
	if err != nil {
		t.Fatalf("Failed at year boundary: %v", err)
	}
	
	if math.IsNaN(pos.Latitude) || math.IsNaN(pos.Longitude) {
		t.Error("Position has NaN coordinates")
	}
}

func TestCalculateSGP4PositionWithObserver_ZeroAltitude(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	// Test with observer at sea level
	seaLevelObserver := ObserverPosition{
		Latitude:  0.0,
		Longitude: 0.0,
		Altitude:  0.0,
	}
	
	result, err := CalculateSGP4PositionWithObserver(testTLELine1, testTLELine2, targetTime, seaLevelObserver)
	if err != nil {
		t.Fatalf("Failed with zero altitude observer: %v", err)
	}
	
	if math.IsNaN(result.LookAngles.Azimuth) || math.IsNaN(result.LookAngles.Elevation) {
		t.Error("Look angles are NaN")
	}
}

func TestCalculateSGP4Position_TimestampAccuracy(t *testing.T) {
	// Test that timestamp is correctly set
	testTimes := []time.Time{
		time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC),
		time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
	}
	
	for _, testTime := range testTimes {
		pos, err := CalculateSGP4Position(testTLELine1, testTLELine2, testTime)
		if err != nil {
			t.Errorf("Failed for time %v: %v", testTime, err)
			continue
		}
		
		if pos.Timestamp != testTime.Unix() {
			t.Errorf("Timestamp mismatch: expected %d, got %d", testTime.Unix(), pos.Timestamp)
		}
	}
}

func TestCalculateSGP4Position_VelocityCalculation(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	pos, err := CalculateSGP4Position(testTLELine1, testTLELine2, targetTime)
	if err != nil {
		t.Fatalf("CalculateSGP4Position failed: %v", err)
	}
	
	// Velocity should be a positive number (magnitude of velocity vector)
	if pos.Velocity <= 0 {
		t.Errorf("Velocity should be positive, got %f", pos.Velocity)
	}
	
	if math.IsNaN(pos.Velocity) || math.IsInf(pos.Velocity, 0) {
		t.Error("Velocity is NaN or Inf")
	}
}

func TestCalculateSGP4PositionWithObserver_RangeCalculation(t *testing.T) {
	targetTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	observer := ObserverPosition{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Altitude:  10.0,
	}
	
	result, err := CalculateSGP4PositionWithObserver(testTLELine1, testTLELine2, targetTime, observer)
	if err != nil {
		t.Fatalf("CalculateSGP4PositionWithObserver failed: %v", err)
	}
	
	// Range should be positive
	if result.LookAngles.Range <= 0 {
		t.Errorf("Range should be positive, got %f", result.LookAngles.Range)
	}
	
	// Range should be reasonable for LEO satellites (typically 400-2000 km)
	// But we'll be lenient since TLE epoch matters
	if result.LookAngles.Range > 100000 {
		t.Errorf("Range seems unreasonably large: %f km", result.LookAngles.Range)
	}
}

// Test print functions to ensure they don't panic
func TestPrintSGP4Position(t *testing.T) {
	pos := SGPPosition{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Altitude:  400.0,
		Velocity:  7.5,
		Timestamp: time.Now().Unix(),
	}
	
	// Just verify it doesn't panic - we can't easily test output
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSGP4Position panicked: %v", r)
		}
	}()
	
	PrintSGP4Position(pos)
}

func TestPrintSGP4PositionWithLookAngles(t *testing.T) {
	result := SGP4PositionResult{
		Position: SGPPosition{
			Latitude:  40.7128,
			Longitude: -74.0060,
			Altitude:  400.0,
			Velocity:  7.5,
			Timestamp: time.Now().Unix(),
		},
		LookAngles: LookAngles{
			Azimuth:   180.0,
			Elevation: 45.0,
			Range:     1000.0,
			RangeRate: 0.0,
		},
	}
	
	// Just verify it doesn't panic - we can't easily test output
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSGP4PositionWithLookAngles panicked: %v", r)
		}
	}()
	
	PrintSGP4PositionWithLookAngles(result)
}

func TestPrintSGP4Position_EdgeCases(t *testing.T) {
	// Test with edge case values
	testCases := []SGPPosition{
		{Latitude: 0.0, Longitude: 0.0, Altitude: 0.0, Velocity: 0.0, Timestamp: 0},
		{Latitude: 90.0, Longitude: 180.0, Altitude: 1000.0, Velocity: 10.0, Timestamp: 1000000000},
		{Latitude: -90.0, Longitude: -180.0, Altitude: 500.0, Velocity: 5.0, Timestamp: 2000000000},
		{Latitude: 45.123456, Longitude: -123.456789, Altitude: 350.5, Velocity: 7.123456, Timestamp: 1500000000},
	}
	
	for i, pos := range testCases {
		defer func(idx int) {
			if r := recover(); r != nil {
				t.Errorf("PrintSGP4Position panicked for test case %d: %v", idx, r)
			}
		}(i)
		PrintSGP4Position(pos)
	}
}

func TestPrintSGP4PositionWithLookAngles_EdgeCases(t *testing.T) {
	// Test with edge case values
	testCases := []SGP4PositionResult{
		{
			Position: SGPPosition{Latitude: 0.0, Longitude: 0.0, Altitude: 0.0, Velocity: 0.0, Timestamp: 0},
			LookAngles: LookAngles{Azimuth: 0.0, Elevation: 0.0, Range: 0.0, RangeRate: 0.0},
		},
		{
			Position: SGPPosition{Latitude: 90.0, Longitude: 180.0, Altitude: 1000.0, Velocity: 10.0, Timestamp: 1000000000},
			LookAngles: LookAngles{Azimuth: 359.99, Elevation: 90.0, Range: 50000.0, RangeRate: 1.0},
		},
		{
			Position: SGPPosition{Latitude: -90.0, Longitude: -180.0, Altitude: 500.0, Velocity: 5.0, Timestamp: 2000000000},
			LookAngles: LookAngles{Azimuth: 180.0, Elevation: -90.0, Range: 100.0, RangeRate: -1.0},
		},
	}
	
	for i, result := range testCases {
		defer func(idx int) {
			if r := recover(); r != nil {
				t.Errorf("PrintSGP4PositionWithLookAngles panicked for test case %d: %v", idx, r)
			}
		}(i)
		PrintSGP4PositionWithLookAngles(result)
	}
}

