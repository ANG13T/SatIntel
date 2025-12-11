package osint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TwiN/go-color"
	"github.com/iskaa02/qalam/gradient"
	"github.com/manifoldco/promptui"
)

// SatellitePositionVisualization provides an interactive menu for viewing satellite positions.
func SatellitePositionVisualization() {
	options, _ := os.ReadFile("txt/orbital_element.txt")
	opt, _ := gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 3)

	if selection == 1 {
		result := SelectSatellite()

		if result == "" {
			return
		}

		GetLocation(extractNorad(result))

	} else if selection == 2 {
		fmt.Print("\n ENTER NORAD ID > ")
		var norad string
		fmt.Scanln(&norad)
		GetLocation(norad)
	}
}

// GetLocation fetches and displays the current position of a satellite for a given observer location.
func GetLocation(norad string) {
	// Automatically detect user location
	latitude, longitude, autoDetected := GetLocationWithPrompt()
	if latitude == "" || longitude == "" {
		return
	}

	if autoDetected {
		fmt.Println(color.Ize(color.Green, "  [+] Using auto-detected location"))
	}

	fmt.Print("\n ENTER ALTITUDE (meters, default: 0) > ")
	var altitude string
	fmt.Scanln(&altitude)
	if strings.TrimSpace(altitude) == "" {
		altitude = "0"
	}

	// Validate inputs
	_, err := strconv.ParseFloat(latitude, 64)
	_, err2 := strconv.ParseFloat(longitude, 64)
	_, err3 := strconv.Atoi(altitude)

	if err != nil || err2 != nil || err3 != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: INVALID INPUT - Please enter valid numbers"))
		return
	}

	spinner := ShowProgressWithSpinner("Fetching satellite position data")
	url := "https://api.n2yo.com/rest/v1/satellite/positions/" + norad + "/" + latitude + "/" + longitude + "/" + altitude + "/2/&apiKey=" + os.Getenv("N2YO_API_KEY")
	resp, err := http.Get(url)
	spinner.Stop()
	if err != nil {
		context := fmt.Sprintf("NORAD ID: %s, Latitude: %s, Longitude: %s", norad, latitude, longitude)
		HandleErrorWithContext(err, ErrCodeAPIRequestFailed, "Failed to fetch satellite position data from N2YO API", context)
		return
	}
	defer resp.Body.Close()

	var data Response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		context := fmt.Sprintf("NORAD ID: %s", norad)
		HandleErrorWithContext(err, ErrCodeAPIParseFailed, "Failed to parse satellite position response", context)
		return
	}

	fmt.Println(color.Ize(color.Purple, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Purple, "║                    Satellite Information                    ║"))
	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))

	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Name", data.SatelliteInfo.Satname)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite ID", fmt.Sprintf("%d", data.SatelliteInfo.Satid))))

	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
	fmt.Println(color.Ize(color.Purple, "║                     Satellite Positions                     ║"))
	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))

	for in, pos := range data.Positions {
		PrintSatellitePosition(pos, in == len(data.Positions)-1)
	}

	// Offer map visualization option
	mapPrompt := promptui.Prompt{
		Label:     "View map visualization? (y/n)",
		Default:   "n",
		AllowEdit: true,
	}
	mapAnswer, _ := mapPrompt.Run()
	if strings.ToLower(strings.TrimSpace(mapAnswer)) == "y" {
		DisplayMap(data)
	}

	// Offer export option
	exportPrompt := promptui.Prompt{
		Label:     "Export satellite positions? (y/n)",
		Default:   "n",
		AllowEdit: true,
	}
	exportAnswer, _ := exportPrompt.Run()
	if strings.ToLower(strings.TrimSpace(exportAnswer)) == "y" {
		defaultFilename := fmt.Sprintf("positions_%s_%d", strings.ReplaceAll(data.SatelliteInfo.Satname, " ", "_"), data.SatelliteInfo.Satid)
		format, filePath, err := showExportMenu(defaultFilename)
		if err == nil {
			if err := ExportSatellitePosition(data, format, filePath); err != nil {
				fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to export: "+err.Error()))
			} else {
				fmt.Println(color.Ize(color.Green, fmt.Sprintf("  [+] Exported to: %s", filePath)))
			}
		}
	}
}

// DisplayMap provides interactive map visualization options for satellite positions.
// It offers three visualization methods: ASCII terminal map, KML export, and web-based map.
func DisplayMap(data Response) {
	if len(data.Positions) == 0 {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: No position data available for visualization"))
		return
	}

	fmt.Println(color.Ize(color.Cyan, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Cyan, "║              Map Visualization Options                     ║"))
	fmt.Println(color.Ize(color.Cyan, "╠═════════════════════════════════════════════════════════════╣"))
	fmt.Println(color.Ize(color.Cyan, "║  1. Terminal ASCII Map                                     ║"))
	fmt.Println(color.Ize(color.Cyan, "║  2. Export to KML (Google Earth)                           ║"))
	fmt.Println(color.Ize(color.Cyan, "║  3. Web-based Interactive Map                               ║"))
	fmt.Println(color.Ize(color.Cyan, "║  0. Cancel                                                 ║"))
	fmt.Println(color.Ize(color.Cyan, "╚═════════════════════════════════════════════════════════════╝"))

	selection := Option(0, 3)

	switch selection {
	case 1:
		displayASCIIMap(data)
	case 2:
		exportToKML(data)
	case 3:
		generateWebMap(data)
	}
}

// displayASCIIMap creates a terminal-based ASCII visualization of satellite positions.
// It loads the world map from txt/map.txt and overlays satellite positions with telemetry data.
func displayASCIIMap(data Response) {
	fmt.Println(color.Ize(color.Green, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Green, "║              ASCII Map Visualization                      ║"))
	fmt.Println(color.Ize(color.Green, "╠═════════════════════════════════════════════════════════════╣"))
	fmt.Printf(color.Ize(color.Green, "║  Satellite: %-45s ║\n"), data.SatelliteInfo.Satname)
	fmt.Printf(color.Ize(color.Green, "║  NORAD ID: %-47d ║\n"), data.SatelliteInfo.Satid)
	fmt.Println(color.Ize(color.Green, "╚═════════════════════════════════════════════════════════════╝\n"))

	// Load world map from txt/map.txt
	mapContent, err := os.ReadFile("txt/map.txt")
	if err != nil {
		// Fallback to generated map if file not found
		fmt.Println(color.Ize(color.Yellow, "  [*] Map file not found, using generated map..."))
		displayASCIIMapGenerated(data)
		return
	}

	// Parse map into lines
	mapLines := strings.Split(string(mapContent), "\n")
	if len(mapLines) == 0 {
		fmt.Println(color.Ize(color.Yellow, "  [*] Map file is empty, using generated map..."))
		displayASCIIMapGenerated(data)
		return
	}

	// Find maximum line length to determine map width
	maxWidth := 0
	for _, line := range mapLines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	mapHeight := len(mapLines)
	mapWidth := maxWidth

	// Create map grid from loaded map
	mapGrid := make([][]rune, mapHeight)
	for i, line := range mapLines {
		mapGrid[i] = make([]rune, mapWidth)
		// Copy existing map characters
		for j, char := range line {
			if j < mapWidth {
				mapGrid[i][j] = char
			}
		}
		// Fill remaining with spaces
		for j := len(line); j < mapWidth; j++ {
			mapGrid[i][j] = ' '
		}
	}

	// Plot satellite positions on the map
	positionMarkers := make([]struct {
		row int
		col int
		pos Position
		idx int
	}, len(data.Positions))

	for i, pos := range data.Positions {
		// Convert lat/lon to map coordinates
		// Latitude: -90 to 90 -> 0 to mapHeight-1
		// Longitude: -180 to 180 -> 0 to mapWidth-1
		row := int((90.0 - pos.Satlatitude) / 180.0 * float64(mapHeight-1))
		col := int((pos.Satlongitude + 180.0) / 360.0 * float64(mapWidth-1))

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

		positionMarkers[i] = struct {
			row int
			col int
			pos Position
			idx int
		}{row, col, pos, i}

		// Use different symbols for different positions
		symbol := '*'
		if i == 0 {
			symbol = '●' // First position
		} else if i == len(data.Positions)-1 {
			symbol = '○' // Last position
		} else {
			symbol = '·' // Intermediate positions
		}

		// Only overlay if the position is on a non-land character (water/space)
		// This helps visibility
		if row >= 0 && row < mapHeight && col >= 0 && col < mapWidth {
			currentChar := mapGrid[row][col]
			// Overlay on spaces or light characters, preserve land mass
			if currentChar == ' ' || currentChar == '.' || currentChar == '-' || currentChar == '_' {
				mapGrid[row][col] = symbol
			} else {
				// Try adjacent cells if current is land
				if col+1 < mapWidth && (mapGrid[row][col+1] == ' ' || mapGrid[row][col+1] == '.') {
					mapGrid[row][col+1] = symbol
					positionMarkers[i].col = col + 1
				} else if col-1 >= 0 && (mapGrid[row][col-1] == ' ' || mapGrid[row][col-1] == '.') {
					mapGrid[row][col-1] = symbol
					positionMarkers[i].col = col - 1
				} else {
					mapGrid[row][col] = symbol
				}
			}
		}
	}

	// Display the map with positions
	fmt.Println(color.Ize(color.Cyan, "                    WORLD MAP - SATELLITE POSITIONS"))
	fmt.Println(color.Ize(color.Yellow, "    Longitude: -180°                                   0°                                   180°\n"))
	
	for i, row := range mapGrid {
		// Print latitude labels on the left
		lat := 90.0 - float64(i)*180.0/float64(mapHeight-1)
		if i%4 == 0 || i == 0 || i == mapHeight-1 {
			fmt.Printf(color.Ize(color.Yellow, "%5.0f° "), lat)
		} else {
			fmt.Print("      ")
		}
		
		// Print map row with colored positions
		for j, cell := range row {
			char := string(cell)
			// Check if this is a position marker
			isMarker := false
			markerIdx := -1
			for idx, marker := range positionMarkers {
				if marker.row == i && marker.col == j {
					isMarker = true
					markerIdx = idx
					break
				}
			}
			
			if isMarker {
				// Color code the markers
				if markerIdx == 0 {
					fmt.Print(color.Ize(color.Red, char)) // First position - red
				} else if markerIdx == len(positionMarkers)-1 {
					fmt.Print(color.Ize(color.Green, char)) // Last position - green
				} else {
					fmt.Print(color.Ize(color.Cyan, char)) // Intermediate - cyan
				}
			} else {
				// Regular map characters in dim color
				fmt.Print(color.Ize(color.White, char))
			}
		}
		fmt.Println()
	}
	fmt.Println()

	// Display telemetry data in a formatted table
	fmt.Println(color.Ize(color.Green, "╔════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Green, "║                                    SATELLITE TELEMETRY DATA                                    ║"))
	fmt.Println(color.Ize(color.Green, "╠════════════════════════════════════════════════════════════════════════════════════════════════════════════════╣"))
	
	for i, pos := range data.Positions {
		// Format timestamp
		timestamp := time.Unix(pos.Timestamp, 0)
		timeStr := timestamp.Format("2006-01-02 15:04:05 UTC")
		
		// Determine position type
		posType := "Intermediate"
		posColor := color.Cyan
		if i == 0 {
			posType = "First"
			posColor = color.Red
		} else if i == len(data.Positions)-1 {
			posType = "Last"
			posColor = color.Green
		}
		
		fmt.Println(color.Ize(color.Green, "╠════════════════════════════════════════════════════════════════════════════════════════════════════════════════╣"))
		fmt.Printf(color.Ize(posColor, "║  Position #%d (%s)                                                                                              ║\n"), i+1, posType)
		fmt.Println(color.Ize(color.Green, "╠════════════════════════════════════════════════════════════════════════════════════════════════════════════════╣"))
		fmt.Printf(color.Ize(color.White, "║  Latitude:     %10.6f°                                                                                      ║\n"), pos.Satlatitude)
		fmt.Printf(color.Ize(color.White, "║  Longitude:    %10.6f°                                                                                      ║\n"), pos.Satlongitude)
		fmt.Printf(color.Ize(color.White, "║  Altitude:     %10.2f km                                                                                    ║\n"), pos.Sataltitude)
		fmt.Printf(color.Ize(color.White, "║  Azimuth:      %10.2f°                                                                                      ║\n"), pos.Azimuth)
		fmt.Printf(color.Ize(color.White, "║  Elevation:    %10.2f°                                                                                      ║\n"), pos.Elevation)
		fmt.Printf(color.Ize(color.White, "║  Right Asc:    %10.2f°                                                                                      ║\n"), pos.Ra)
		fmt.Printf(color.Ize(color.White, "║  Declination:  %10.2f°                                                                                      ║\n"), pos.Dec)
		fmt.Printf(color.Ize(color.White, "║  Timestamp:    %-60s ║\n"), timeStr)
		
		// Show map coordinates
		row := int((90.0 - pos.Satlatitude) / 180.0 * float64(mapHeight-1))
		col := int((pos.Satlongitude + 180.0) / 360.0 * float64(mapWidth-1))
		fmt.Printf(color.Ize(color.Yellow, "║  Map Position: Row %3d, Col %3d                                                                              ║\n"), row, col)
	}
	
	fmt.Println(color.Ize(color.Green, "╚════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝"))

	// Print legend
	fmt.Println(color.Ize(color.Green, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Green, "║                         Legend                            ║"))
	fmt.Println(color.Ize(color.Green, "╠═════════════════════════════════════════════════════════════╣"))
	fmt.Println(color.Ize(color.Red, "║  ● First Position (Red)                                   ║"))
	fmt.Println(color.Ize(color.Cyan, "║  · Intermediate Positions (Cyan)                          ║"))
	fmt.Println(color.Ize(color.Green, "║  ○ Last Position (Green)                                 ║"))
	fmt.Println(color.Ize(color.Green, "╚═════════════════════════════════════════════════════════════╝\n"))
}

// displayASCIIMapGenerated is a fallback function that generates a simple map if txt/map.txt is not available.
func displayASCIIMapGenerated(data Response) {
	// Create a simple ASCII world map representation
	// Map dimensions: 80 columns (longitude) x 24 rows (latitude)
	const mapWidth = 80
	const mapHeight = 24

	// Initialize map grid
	mapGrid := make([][]rune, mapHeight)
	for i := range mapGrid {
		mapGrid[i] = make([]rune, mapWidth)
		for j := range mapGrid[i] {
			mapGrid[i][j] = ' '
		}
	}

	// Draw basic world map outline (simplified)
	drawWorldMapOutline(mapGrid)

	// Plot satellite positions
	for i, pos := range data.Positions {
		// Convert lat/lon to map coordinates
		row := int((90.0 - pos.Satlatitude) / 180.0 * float64(mapHeight-1))
		col := int((pos.Satlongitude + 180.0) / 360.0 * float64(mapWidth-1))

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

		// Use different symbols for different positions
		symbol := '*'
		if i == 0 {
			symbol = '●'
		} else if i == len(data.Positions)-1 {
			symbol = '○'
		} else {
			symbol = '·'
		}

		if row >= 0 && row < mapHeight && col >= 0 && col < mapWidth {
			mapGrid[row][col] = symbol
		}
	}

	// Print the map
	fmt.Println(color.Ize(color.Yellow, "    Longitude: -180°                                   0°                                   180°"))
	fmt.Println(color.Ize(color.Yellow, "Latitude"))
	for i, row := range mapGrid {
		lat := 90.0 - float64(i)*180.0/float64(mapHeight-1)
		if i%3 == 0 || i == 0 || i == mapHeight-1 {
			fmt.Printf(color.Ize(color.Yellow, "%5.0f° "), lat)
		} else {
			fmt.Print("      ")
		}
		fmt.Print("│")
		for _, cell := range row {
			if cell == ' ' {
				fmt.Print(" ")
			} else {
				fmt.Print(color.Ize(color.Cyan, string(cell)))
			}
		}
		fmt.Println("│")
	}
	fmt.Println(color.Ize(color.Yellow, "      └────────────────────────────────────────────────────────────────────────┘"))

	// Print legend
	fmt.Println(color.Ize(color.Green, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Green, "║                         Legend                            ║"))
	fmt.Println(color.Ize(color.Green, "╠═════════════════════════════════════════════════════════════╣"))
	fmt.Println(color.Ize(color.Green, "║  ● First Position                                        ║"))
	fmt.Println(color.Ize(color.Green, "║  · Intermediate Positions                                ║"))
	fmt.Println(color.Ize(color.Green, "║  ○ Last Position                                         ║"))
	fmt.Println(color.Ize(color.Green, "╚═════════════════════════════════════════════════════════════╝\n"))

	// Print position details
	fmt.Println(color.Ize(color.Cyan, "\nPosition Details:"))
	for i, pos := range data.Positions {
		fmt.Printf(color.Ize(color.Cyan, "  Position %d: Lat %.4f°, Lon %.4f°, Alt %.2f km\n"),
			i+1, pos.Satlatitude, pos.Satlongitude, pos.Sataltitude)
	}
	fmt.Println()
}

// drawWorldMapOutline draws a simplified ASCII world map outline.
func drawWorldMapOutline(grid [][]rune) {
	height := len(grid)
	width := len(grid[0])

	// Draw continents (very simplified)
	// North America
	for i := height/3; i < height*2/3; i++ {
		for j := width/8; j < width*3/8; j++ {
			if (i-height/3)%2 == 0 && (j-width/8)%3 != 0 {
				grid[i][j] = '▓'
			}
		}
	}

	// Europe/Africa
	for i := height/3; i < height*2/3; i++ {
		for j := width*2/5; j < width*3/5; j++ {
			if (i-height/3)%2 == 0 && (j-width*2/5)%3 != 0 {
				grid[i][j] = '▓'
			}
		}
	}

	// Asia
	for i := height/3; i < height*2/3; i++ {
		for j := width*3/5; j < width*4/5; j++ {
			if (i-height/3)%2 == 0 && (j-width*3/5)%3 != 0 {
				grid[i][j] = '▓'
			}
		}
	}

	// Australia
	for i := height*2/3; i < height*4/5; i++ {
		for j := width*3/4; j < width*7/8; j++ {
			if (i-height*2/3)%2 == 0 && (j-width*3/4)%2 != 0 {
				grid[i][j] = '▓'
			}
		}
	}

	// South America
	for i := height*2/3; i < height*4/5; i++ {
		for j := width/6; j < width/3; j++ {
			if (i-height*2/3)%2 == 0 && (j-width/6)%3 != 0 {
				grid[i][j] = '▓'
			}
		}
	}
}

// exportToKML exports satellite positions to a KML file for Google Earth.
func exportToKML(data Response) {
	defaultFilename := fmt.Sprintf("satellite_%s_%d.kml",
		strings.ReplaceAll(data.SatelliteInfo.Satname, " ", "_"), data.SatelliteInfo.Satid)

	pathPrompt := promptui.Prompt{
		Label:     "Enter KML file path (or press Enter for default)",
		Default:   defaultFilename,
		AllowEdit: true,
	}

	filePath, err := pathPrompt.Run()
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] Export cancelled"))
		return
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		filePath = defaultFilename
	}

	// Ensure .kml extension
	if !strings.HasSuffix(strings.ToLower(filePath), ".kml") {
		filePath += ".kml"
	}

	// Generate KML content
	kmlContent := generateKMLContent(data)

	// Write to file
	if err := os.WriteFile(filePath, []byte(kmlContent), 0644); err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to write KML file: "+err.Error()))
		return
	}

	fmt.Println(color.Ize(color.Green, fmt.Sprintf("  [+] KML file exported to: %s", filePath)))
	fmt.Println(color.Ize(color.Cyan, "  [*] You can open this file in Google Earth or other KML-compatible applications"))
}

// generateKMLContent creates KML XML content for satellite positions.
func generateKMLContent(data Response) string {
	var builder strings.Builder

	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString("\n<kml xmlns=\"http://www.opengis.net/kml/2.2\">\n")
	builder.WriteString("  <Document>\n")
	builder.WriteString(fmt.Sprintf("    <name>%s (NORAD ID: %d)</name>\n",
		data.SatelliteInfo.Satname, data.SatelliteInfo.Satid))
	builder.WriteString("    <description>Satellite position data exported from SatIntel</description>\n")

	// Add a style for satellite markers
	builder.WriteString("    <Style id=\"satelliteStyle\">\n")
	builder.WriteString("      <IconStyle>\n")
	builder.WriteString("        <color>ff00ffff</color>\n")
	builder.WriteString("        <scale>1.2</scale>\n")
	builder.WriteString("        <Icon>\n")
	builder.WriteString("          <href>http://maps.google.com/mapfiles/kml/shapes/arrow.png</href>\n")
	builder.WriteString("        </Icon>\n")
	builder.WriteString("      </IconStyle>\n")
	builder.WriteString("      <LabelStyle>\n")
	builder.WriteString("        <scale>0.8</scale>\n")
	builder.WriteString("      </LabelStyle>\n")
	builder.WriteString("    </Style>\n")

	// Add placemarks for each position
	for i, pos := range data.Positions {
		builder.WriteString("    <Placemark>\n")
		builder.WriteString(fmt.Sprintf("      <name>Position %d</name>\n", i+1))
		builder.WriteString(fmt.Sprintf("      <description>\n"))
		builder.WriteString(fmt.Sprintf("        Latitude: %.6f°\n", pos.Satlatitude))
		builder.WriteString(fmt.Sprintf("        Longitude: %.6f°\n", pos.Satlongitude))
		builder.WriteString(fmt.Sprintf("        Altitude: %.2f km\n", pos.Sataltitude))
		builder.WriteString(fmt.Sprintf("        Timestamp: %d\n", pos.Timestamp))
		builder.WriteString(fmt.Sprintf("      </description>\n"))
		builder.WriteString("      <styleUrl>#satelliteStyle</styleUrl>\n")
		builder.WriteString("      <Point>\n")
		builder.WriteString(fmt.Sprintf("        <coordinates>%.6f,%.6f,%.2f</coordinates>\n",
			pos.Satlongitude, pos.Satlatitude, pos.Sataltitude*1000)) // KML uses meters
		builder.WriteString("      </Point>\n")
		builder.WriteString("    </Placemark>\n")

		// Add a path between positions if not the last one
		if i < len(data.Positions)-1 {
			nextPos := data.Positions[i+1]
			builder.WriteString("    <Placemark>\n")
			builder.WriteString(fmt.Sprintf("      <name>Path %d-%d</name>\n", i+1, i+2))
			builder.WriteString("      <styleUrl>#satelliteStyle</styleUrl>\n")
			builder.WriteString("      <LineString>\n")
			builder.WriteString("        <tessellate>1</tessellate>\n")
			builder.WriteString("        <coordinates>\n")
			builder.WriteString(fmt.Sprintf("          %.6f,%.6f,%.2f\n",
				pos.Satlongitude, pos.Satlatitude, pos.Sataltitude*1000))
			builder.WriteString(fmt.Sprintf("          %.6f,%.6f,%.2f\n",
				nextPos.Satlongitude, nextPos.Satlatitude, nextPos.Sataltitude*1000))
			builder.WriteString("        </coordinates>\n")
			builder.WriteString("      </LineString>\n")
			builder.WriteString("    </Placemark>\n")
		}
	}

	builder.WriteString("  </Document>\n")
	builder.WriteString("</kml>\n")

	return builder.String()
}

// generateWebMap creates an HTML file with an interactive web-based map using Leaflet.
func generateWebMap(data Response) {
	defaultFilename := fmt.Sprintf("satellite_map_%s_%d.html",
		strings.ReplaceAll(data.SatelliteInfo.Satname, " ", "_"), data.SatelliteInfo.Satid)

	pathPrompt := promptui.Prompt{
		Label:     "Enter HTML file path (or press Enter for default)",
		Default:   defaultFilename,
		AllowEdit: true,
	}

	filePath, err := pathPrompt.Run()
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] Export cancelled"))
		return
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		filePath = defaultFilename
	}

	// Ensure .html extension
	if !strings.HasSuffix(strings.ToLower(filePath), ".html") {
		filePath += ".html"
	}

	// Generate HTML content
	htmlContent := generateHTMLMapContent(data)

	// Write to file
	if err := os.WriteFile(filePath, []byte(htmlContent), 0644); err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to write HTML file: "+err.Error()))
		return
	}

	fmt.Println(color.Ize(color.Green, fmt.Sprintf("  [+] Interactive map exported to: %s", filePath)))
	fmt.Println(color.Ize(color.Cyan, "  [*] Open this file in your web browser to view the interactive map"))
}

// generateHTMLMapContent creates HTML content with Leaflet.js for interactive map visualization.
func generateHTMLMapContent(data Response) string {
	var builder strings.Builder

	// Prepare position data as JSON for JavaScript
	positionsJSON, _ := json.Marshal(data.Positions)

	builder.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Satellite Map - `)
	builder.WriteString(data.SatelliteInfo.Satname)
	builder.WriteString(`</title>
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <style>
        body {
            margin: 0;
            padding: 0;
            font-family: Arial, sans-serif;
        }
        #map {
            height: 100vh;
            width: 100%;
        }
        .info-panel {
            position: absolute;
            top: 10px;
            right: 10px;
            background: white;
            padding: 15px;
            border-radius: 5px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.3);
            z-index: 1000;
            max-width: 300px;
        }
        .info-panel h3 {
            margin-top: 0;
            color: #333;
        }
        .info-panel p {
            margin: 5px 0;
            font-size: 14px;
        }
        .position-marker {
            background: #00ffff;
            border: 2px solid #fff;
            border-radius: 50%;
            width: 12px;
            height: 12px;
        }
    </style>
</head>
<body>
    <div id="map"></div>
    <div class="info-panel">
        <h3>`)
	builder.WriteString(data.SatelliteInfo.Satname)
	builder.WriteString(`</h3>
        <p><strong>NORAD ID:</strong> `)
	builder.WriteString(fmt.Sprintf("%d", data.SatelliteInfo.Satid))
	builder.WriteString(`</p>
        <p><strong>Positions:</strong> `)
	builder.WriteString(fmt.Sprintf("%d", len(data.Positions)))
	builder.WriteString(`</p>
        <p><small>Click on markers to see details</small></p>
    </div>

    <script>
        // Initialize map
        var map = L.map('map').setView([`)
	
	// Center map on first position
	if len(data.Positions) > 0 {
		builder.WriteString(fmt.Sprintf("%.6f, %.6f", data.Positions[0].Satlatitude, data.Positions[0].Satlongitude))
	} else {
		builder.WriteString("0, 0")
	}
	
	builder.WriteString(`], 2);

        // Add tile layer
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '© OpenStreetMap contributors',
            maxZoom: 19
        }).addTo(map);

        // Position data
        var positions = `)
	builder.WriteString(string(positionsJSON))
	builder.WriteString(`;

        // Create polyline for path
        var pathCoordinates = positions.map(function(pos) {
            return [pos.satlatitude, pos.satlongitude];
        });
        
        var polyline = L.polyline(pathCoordinates, {
            color: '#00ffff',
            weight: 3,
            opacity: 0.7
        }).addTo(map);

        // Fit map to show all positions
        map.fitBounds(polyline.getBounds());

        // Add markers for each position
        positions.forEach(function(pos, index) {
            var marker = L.circleMarker([pos.satlatitude, pos.satlongitude], {
                radius: 8,
                fillColor: index === 0 ? '#ff0000' : (index === positions.length - 1 ? '#00ff00' : '#00ffff'),
                color: '#fff',
                weight: 2,
                opacity: 1,
                fillOpacity: 0.8
            }).addTo(map);

            var popupContent = '<div style="min-width: 200px;">' +
                '<h4>Position ' + (index + 1) + '</h4>' +
                '<p><strong>Latitude:</strong> ' + pos.satlatitude.toFixed(6) + '°</p>' +
                '<p><strong>Longitude:</strong> ' + pos.satlongitude.toFixed(6) + '°</p>' +
                '<p><strong>Altitude:</strong> ' + pos.sataltitude.toFixed(2) + ' km</p>' +
                '<p><strong>Timestamp:</strong> ' + new Date(pos.timestamp * 1000).toLocaleString() + '</p>' +
                '</div>';

            marker.bindPopup(popupContent);
        });

        // Add legend
        var legend = L.control({position: 'bottomright'});
        legend.onAdd = function(map) {
            var div = L.DomUtil.create('div', 'info legend');
            div.style.backgroundColor = 'white';
            div.style.padding = '10px';
            div.style.borderRadius = '5px';
            div.innerHTML = '<h4>Legend</h4>' +
                '<p><span style="color: #ff0000;">●</span> First Position</p>' +
                '<p><span style="color: #00ffff;">●</span> Intermediate</p>' +
                '<p><span style="color: #00ff00;">●</span> Last Position</p>';
            return div;
        };
        legend.addTo(map);
    </script>
</body>
</html>`)

	return builder.String()
}

// PrintSatellitePosition displays satellite position data in a formatted table.
func PrintSatellitePosition(pos Position, last bool) {
	fmt.Println(color.Ize(color.Purple, GenRowString("Latitude", fmt.Sprintf("%f", pos.Satlatitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Longitude", fmt.Sprintf("%f", pos.Satlongitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Altitude", fmt.Sprintf("%f", pos.Sataltitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Right Ascension", fmt.Sprintf("%f", pos.Azimuth))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Declination", fmt.Sprintf("%f", pos.Dec))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Timestamp", fmt.Sprintf("%d", pos.Timestamp))))
	if last {
		fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝\n\n"))
	} else {
		fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
	}
}
