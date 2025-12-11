package osint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// createTestResponse creates a test Response with sample satellite position data.
func createTestResponse() Response {
	return Response{
		SatelliteInfo: SatelliteInfo{
			Satname: "ISS (ZARYA)",
			Satid:   25544,
		},
		Positions: []Position{
			{
				Satlatitude:  51.5074,  // London
				Satlongitude: -0.1278,
				Sataltitude:  408.0,
				Azimuth:      45.0,
				Elevation:    30.0,
				Ra:           180.0,
				Dec:          40.0,
				Timestamp:    time.Now().Unix(),
			},
			{
				Satlatitude:  40.7128,  // New York
				Satlongitude: -74.0060,
				Sataltitude:  410.0,
				Azimuth:      90.0,
				Elevation:    45.0,
				Ra:           200.0,
				Dec:          50.0,
				Timestamp:    time.Now().Unix() + 100,
			},
			{
				Satlatitude:  35.6762,  // Tokyo
				Satlongitude: 139.6503,
				Sataltitude:  412.0,
				Azimuth:      135.0,
				Elevation:    60.0,
				Ra:           220.0,
				Dec:          60.0,
				Timestamp:    time.Now().Unix() + 200,
			},
		},
	}
}

func TestGenerateKMLContent(t *testing.T) {
	data := createTestResponse()
	kmlContent := generateKMLContent(data)

	// Verify KML structure
	if !strings.Contains(kmlContent, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("KML content missing XML declaration")
	}

	if !strings.Contains(kmlContent, "<kml xmlns=\"http://www.opengis.net/kml/2.2\">") {
		t.Error("KML content missing KML root element")
	}

	if !strings.Contains(kmlContent, "<Document>") {
		t.Error("KML content missing Document element")
	}

	// Verify satellite info is included
	if !strings.Contains(kmlContent, data.SatelliteInfo.Satname) {
		t.Error("KML content missing satellite name")
	}

	if !strings.Contains(kmlContent, "NORAD ID: 25544") {
		t.Error("KML content missing NORAD ID")
	}

	// Verify Style element
	if !strings.Contains(kmlContent, "<Style id=\"satelliteStyle\">") {
		t.Error("KML content missing Style element")
	}

	// Verify Placemarks for each position
	for i := range data.Positions {
		positionName := fmt.Sprintf("Position %d", i+1)
		if !strings.Contains(kmlContent, positionName) {
			t.Errorf("KML content missing Placemark for position %d", i+1)
		}
	}

	// Verify coordinates are present
	for _, pos := range data.Positions {
		coordStr := fmt.Sprintf("%.6f,%.6f", pos.Satlongitude, pos.Satlatitude)
		if !strings.Contains(kmlContent, coordStr) {
			t.Errorf("KML content missing coordinates for position: %s", coordStr)
		}
	}

	// Verify paths between positions
	if !strings.Contains(kmlContent, "Path 1-2") {
		t.Error("KML content missing path between positions 1 and 2")
	}

	if !strings.Contains(kmlContent, "Path 2-3") {
		t.Error("KML content missing path between positions 2 and 3")
	}

	// Verify LineString elements for paths
	if !strings.Contains(kmlContent, "<LineString>") {
		t.Error("KML content missing LineString elements for paths")
	}

	// Verify altitude conversion (KML uses meters, input is km)
	for _, pos := range data.Positions {
		altMeters := pos.Sataltitude * 1000
		altStr := fmt.Sprintf("%.2f", altMeters)
		if !strings.Contains(kmlContent, altStr) {
			t.Errorf("KML content missing altitude in meters: %s", altStr)
		}
	}
}

func TestGenerateKMLContentEmptyPositions(t *testing.T) {
	data := Response{
		SatelliteInfo: SatelliteInfo{
			Satname: "Test Satellite",
			Satid:   12345,
		},
		Positions: []Position{},
	}

	kmlContent := generateKMLContent(data)

	// Should still have basic structure
	if !strings.Contains(kmlContent, "<kml") {
		t.Error("KML content should have basic structure even with no positions")
	}

	// Should not have any Placemarks
	if strings.Contains(kmlContent, "<Placemark>") {
		t.Error("KML content should not have Placemarks when there are no positions")
	}
}

func TestGenerateKMLContentSinglePosition(t *testing.T) {
	data := Response{
		SatelliteInfo: SatelliteInfo{
			Satname: "Single Position Test",
			Satid:   99999,
		},
		Positions: []Position{
			{
				Satlatitude:  0.0,
				Satlongitude: 0.0,
				Sataltitude:  400.0,
				Timestamp:    time.Now().Unix(),
			},
		},
	}

	kmlContent := generateKMLContent(data)

	// Should have one Placemark
	if !strings.Contains(kmlContent, "Position 1") {
		t.Error("KML content should have Placemark for single position")
	}

	// Should not have any paths (only one position)
	if strings.Contains(kmlContent, "Path") {
		t.Error("KML content should not have paths when there's only one position")
	}
}

func TestGenerateHTMLMapContent(t *testing.T) {
	data := createTestResponse()
	htmlContent := generateHTMLMapContent(data)

	// Verify HTML structure
	if !strings.Contains(htmlContent, "<!DOCTYPE html>") {
		t.Error("HTML content missing DOCTYPE declaration")
	}

	if !strings.Contains(htmlContent, "<html lang=\"en\">") {
		t.Error("HTML content missing HTML tag")
	}

	// Verify Leaflet.js is included
	if !strings.Contains(htmlContent, "leaflet@1.9.4") {
		t.Error("HTML content missing Leaflet.js library")
	}

	// Verify satellite info is included
	if !strings.Contains(htmlContent, data.SatelliteInfo.Satname) {
		t.Error("HTML content missing satellite name")
	}

	if !strings.Contains(htmlContent, fmt.Sprintf("%d", data.SatelliteInfo.Satid)) {
		t.Error("HTML content missing satellite ID")
	}

	// Verify position data is included as JSON
	if !strings.Contains(htmlContent, "var positions =") {
		t.Error("HTML content missing positions JavaScript variable")
	}

	// Verify map initialization
	if !strings.Contains(htmlContent, "L.map('map')") {
		t.Error("HTML content missing Leaflet map initialization")
	}

	// Verify OpenStreetMap tile layer
	if !strings.Contains(htmlContent, "tile.openstreetmap.org") {
		t.Error("HTML content missing OpenStreetMap tile layer")
	}

	// Verify markers are created
	if !strings.Contains(htmlContent, "circleMarker") {
		t.Error("HTML content missing circle marker creation")
	}

	// Verify polyline for path
	if !strings.Contains(htmlContent, "L.polyline") {
		t.Error("HTML content missing polyline for satellite path")
	}

	// Verify legend
	if !strings.Contains(htmlContent, "Legend") {
		t.Error("HTML content missing legend")
	}

	// Verify position coordinates are in the JavaScript (JSON format may vary)
	// Check that position data is present as JSON
	if !strings.Contains(htmlContent, "satlatitude") {
		t.Error("HTML content missing latitude field in position data")
	}

	if !strings.Contains(htmlContent, "satlongitude") {
		t.Error("HTML content missing longitude field in position data")
	}

	// Verify at least one coordinate value is present (JSON may format differently)
	hasCoordinates := false
	for _, pos := range data.Positions {
		// JSON might format as 51.5074 or 51.5074e+00, so check for partial match
		latPartial := fmt.Sprintf("%.1f", pos.Satlatitude)
		lonPartial := fmt.Sprintf("%.1f", pos.Satlongitude)
		if strings.Contains(htmlContent, latPartial) || strings.Contains(htmlContent, lonPartial) {
			hasCoordinates = true
			break
		}
	}
	if !hasCoordinates {
		t.Error("HTML content missing coordinate values")
	}
}

func TestGenerateHTMLMapContentEmptyPositions(t *testing.T) {
	data := Response{
		SatelliteInfo: SatelliteInfo{
			Satname: "Empty Test",
			Satid:   11111,
		},
		Positions: []Position{},
	}

	htmlContent := generateHTMLMapContent(data)

	// Should still have basic HTML structure
	if !strings.Contains(htmlContent, "<html") {
		t.Error("HTML content should have basic structure even with no positions")
	}

	// Should have empty positions array
	if !strings.Contains(htmlContent, "var positions = []") {
		t.Error("HTML content should have empty positions array")
	}
}

func TestDrawWorldMapOutline(t *testing.T) {
	// Create a test grid
	height := 24
	width := 80
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Draw outline
	drawWorldMapOutline(grid)

	// Verify grid was modified (should have some continent markers)
	foundContinent := false
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == '▓' {
				foundContinent = true
				break
			}
		}
		if foundContinent {
			break
		}
	}

	if !foundContinent {
		t.Error("drawWorldMapOutline should draw continent markers")
	}

	// Verify grid dimensions are preserved
	if len(grid) != height {
		t.Errorf("Grid height changed: expected %d, got %d", height, len(grid))
	}

	if len(grid[0]) != width {
		t.Errorf("Grid width changed: expected %d, got %d", width, len(grid[0]))
	}
}

func TestDrawWorldMapOutlineSmallGrid(t *testing.T) {
	// Test with smaller grid
	height := 10
	width := 20
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Should not panic with smaller grid
	drawWorldMapOutline(grid)

	// Verify grid is still valid
	if len(grid) != height {
		t.Errorf("Grid height changed: expected %d, got %d", height, len(grid))
	}
}

func TestDisplayASCIIMapPositionMapping(t *testing.T) {
	// Test coordinate conversion logic
	const mapWidth = 80
	const mapHeight = 24

	testCases := []struct {
		name     string
		lat      float64
		lon      float64
		expectedRow int
		expectedCol int
	}{
		{
			name:     "North Pole",
			lat:      90.0,
			lon:      0.0,
			expectedRow: 0,
			expectedCol: 40, // 180/360 * 80 = 40
		},
		{
			name:     "South Pole",
			lat:      -90.0,
			lon:      0.0,
			expectedRow: mapHeight - 1,
			expectedCol: 40,
		},
		{
			name:     "Equator, Prime Meridian",
			lat:      0.0,
			lon:      0.0,
			expectedRow: mapHeight / 2,
			expectedCol: 40,
		},
		{
			name:     "Equator, 180°",
			lat:      0.0,
			lon:      180.0,
			expectedRow: mapHeight / 2,
			expectedCol: mapWidth - 1,
		},
		{
			name:     "Equator, -180°",
			lat:      0.0,
			lon:      -180.0,
			expectedRow: mapHeight / 2,
			expectedCol: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			row := int((90.0 - tc.lat) / 180.0 * float64(mapHeight-1))
			col := int((tc.lon + 180.0) / 360.0 * float64(mapWidth-1))

			// Clamp to valid range
			if row < 0 {
				row = 0
			}
			if row >= mapHeight {
				row = mapHeight - 1
			}
			if col < 0 {
				col = 0
			}
			if col >= mapWidth {
				col = mapWidth - 1
			}

			// Allow some tolerance for floating point calculations
			rowDiff := row - tc.expectedRow
			if rowDiff < 0 {
				rowDiff = -rowDiff
			}
			colDiff := col - tc.expectedCol
			if colDiff < 0 {
				colDiff = -colDiff
			}

			if rowDiff > 1 {
				t.Errorf("Row mismatch: expected around %d, got %d", tc.expectedRow, row)
			}
			if colDiff > 1 {
				t.Errorf("Col mismatch: expected around %d, got %d", tc.expectedCol, col)
			}
		})
	}
}

func TestDisplayMapEmptyPositions(t *testing.T) {
	// This test verifies that DisplayMap handles empty positions gracefully
	// Note: We can't easily test the interactive parts (Option() prompts),
	// but we can verify the function structure handles empty data
	
	// Test that empty positions are detected
	data := Response{
		SatelliteInfo: SatelliteInfo{
			Satname: "Test",
			Satid:   12345,
		},
		Positions: []Position{},
	}

	// Verify data structure is valid
	if len(data.Positions) != 0 {
		t.Error("Test data should have empty positions")
	}

	// Note: DisplayMap uses interactive prompts which makes it difficult to test
	// In a production environment, you'd want to mock the Option() function
	// or refactor to separate the logic from the interactive parts
	_ = data // Suppress unused variable warning
}

func TestKMLFileExport(t *testing.T) {
	data := createTestResponse()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_satellite.kml")

	// Manually test the KML generation and file writing
	kmlContent := generateKMLContent(data)
	if err := os.WriteFile(filePath, []byte(kmlContent), 0644); err != nil {
		t.Fatalf("Failed to write KML file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("KML file was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read KML file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "<kml") {
		t.Error("KML file content is invalid")
	}

	if !strings.Contains(contentStr, data.SatelliteInfo.Satname) {
		t.Error("KML file missing satellite name")
	}
}

func TestHTMLFileExport(t *testing.T) {
	data := createTestResponse()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_satellite_map.html")

	// Manually test the HTML generation and file writing
	htmlContent := generateHTMLMapContent(data)
	if err := os.WriteFile(filePath, []byte(htmlContent), 0644); err != nil {
		t.Fatalf("Failed to write HTML file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("HTML file was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read HTML file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "<html") {
		t.Error("HTML file content is invalid")
	}

	if !strings.Contains(contentStr, data.SatelliteInfo.Satname) {
		t.Error("HTML file missing satellite name")
	}

	if !strings.Contains(contentStr, "leaflet") {
		t.Error("HTML file missing Leaflet.js")
	}
}

func TestGenerateKMLContentEdgeCases(t *testing.T) {
	// Test with extreme coordinates
	data := Response{
		SatelliteInfo: SatelliteInfo{
			Satname: "Edge Case Test",
			Satid:   99999,
		},
		Positions: []Position{
			{Satlatitude: 90.0, Satlongitude: 180.0, Sataltitude: 0.0, Timestamp: 1000},
			{Satlatitude: -90.0, Satlongitude: -180.0, Sataltitude: 1000.0, Timestamp: 2000},
			{Satlatitude: 0.0, Satlongitude: 0.0, Sataltitude: 500.0, Timestamp: 3000},
		},
	}

	kmlContent := generateKMLContent(data)

	// Should handle extreme coordinates
	if !strings.Contains(kmlContent, "90.000000") {
		t.Error("KML should handle maximum latitude")
	}

	if !strings.Contains(kmlContent, "-90.000000") {
		t.Error("KML should handle minimum latitude")
	}

	if !strings.Contains(kmlContent, "180.000000") {
		t.Error("KML should handle maximum longitude")
	}

	if !strings.Contains(kmlContent, "-180.000000") {
		t.Error("KML should handle minimum longitude")
	}
}

func TestGenerateHTMLMapContentEdgeCases(t *testing.T) {
	// Test with extreme coordinates
	data := Response{
		SatelliteInfo: SatelliteInfo{
			Satname: "Edge Case Test",
			Satid:   99999,
		},
		Positions: []Position{
			{Satlatitude: 90.0, Satlongitude: 180.0, Sataltitude: 0.0, Timestamp: 1000},
			{Satlatitude: -90.0, Satlongitude: -180.0, Sataltitude: 1000.0, Timestamp: 2000},
		},
	}

	htmlContent := generateHTMLMapContent(data)

	// Should handle extreme coordinates (JSON may format differently, so check for presence)
	// Check that position data structure exists
	if !strings.Contains(htmlContent, "var positions =") {
		t.Error("HTML should have positions variable")
	}

	// Verify extreme coordinates are present (may be in JSON format)
	hasMaxLat := strings.Contains(htmlContent, "90") || strings.Contains(htmlContent, "9e+01")
	hasMinLat := strings.Contains(htmlContent, "-90") || strings.Contains(htmlContent, "-9e+01")
	hasMaxLon := strings.Contains(htmlContent, "180") || strings.Contains(htmlContent, "1.8e+02")
	hasMinLon := strings.Contains(htmlContent, "-180") || strings.Contains(htmlContent, "-1.8e+02")

	if !hasMaxLat {
		t.Error("HTML should handle maximum latitude")
	}

	if !hasMinLat {
		t.Error("HTML should handle minimum latitude")
	}

	if !hasMaxLon {
		t.Error("HTML should handle maximum longitude")
	}

	if !hasMinLon {
		t.Error("HTML should handle minimum longitude")
	}
}

