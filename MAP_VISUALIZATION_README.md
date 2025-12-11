# Map Visualization Module

## Overview

The Map Visualization module provides three different ways to visualize satellite positions on a map:

1. **Terminal ASCII Map** - A text-based world map displayed directly in the terminal
2. **KML Export** - Export satellite positions to Google Earth format
3. **Web-based Interactive Map** - Generate an HTML file with an interactive map using Leaflet.js

## Features

### 1. Terminal ASCII Map

A terminal-based visualization that displays satellite positions on a simplified ASCII world map.

**Features:**
- 80x24 character grid representing the world
- Simplified continent outlines
- Color-coded position markers:
  - `●` (cyan) - First position
  - `·` (cyan) - Intermediate positions
  - `○` (cyan) - Last position
- Latitude and longitude axis labels
- Position details listing

**Usage:**
When viewing satellite positions, select option `1` from the map visualization menu.

**Example Output:**
```
╔═════════════════════════════════════════════════════════════╗
║              ASCII Map Visualization                      ║
╠═════════════════════════════════════════════════════════════╣
║  Satellite: ISS (ZARYA)                                     ║
╚═════════════════════════════════════════════════════════════╝

    Longitude: -180°                                   0°                                   180°
Latitude
 90° │                                                                    │
     │                                                                    │
     │                                                                    │
  0° │                                                                    │
     │                                                                    │
-90° │                                                                    │
      └────────────────────────────────────────────────────────────────────────┘
```

### 2. KML Export for Google Earth

Export satellite positions to a KML (Keyhole Markup Language) file that can be opened in Google Earth or other KML-compatible applications.

**Features:**
- Placemarks for each satellite position with detailed information
- Paths connecting consecutive positions
- Custom styling for satellite markers
- Altitude information (converted from km to meters for KML)
- Timestamp data for each position

**Usage:**
1. When viewing satellite positions, select option `2` from the map visualization menu
2. Enter a file path (or press Enter for default)
3. Open the generated `.kml` file in Google Earth

**File Format:**
- Extension: `.kml`
- Format: XML-based KML 2.2
- Compatible with: Google Earth, Google Maps, QGIS, and other KML viewers

**Example KML Structure:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>ISS (ZARYA) (NORAD ID: 25544)</name>
    <Placemark>
      <name>Position 1</name>
      <Point>
        <coordinates>-0.127800,51.507400,408000.00</coordinates>
      </Point>
    </Placemark>
    ...
  </Document>
</kml>
```

### 3. Web-based Interactive Map

Generate a standalone HTML file with an interactive map using Leaflet.js and OpenStreetMap.

**Features:**
- Interactive map with zoom and pan controls
- Clickable markers for each satellite position
- Color-coded markers:
  - Red - First position
  - Cyan - Intermediate positions
  - Green - Last position
- Polyline showing the satellite's path
- Info panel with satellite details
- Popup windows with position information
- Automatic map bounds fitting to show all positions
- Legend explaining marker colors

**Usage:**
1. When viewing satellite positions, select option `3` from the map visualization menu
2. Enter a file path (or press Enter for default)
3. Open the generated `.html` file in any web browser

**Requirements:**
- Modern web browser (Chrome, Firefox, Safari, Edge)
- Internet connection (for loading Leaflet.js and OpenStreetMap tiles)

**File Format:**
- Extension: `.html`
- Standalone file (no server required)
- Uses CDN for Leaflet.js library

**Example Features:**
- Click on any marker to see position details
- Zoom in/out using mouse wheel or controls
- Pan by dragging the map
- View satellite path as a connected line

## Integration

The map visualization is automatically integrated into the satellite position viewing workflow:

1. Run `GetLocation()` to fetch satellite positions
2. After displaying position data, you'll be prompted: "View map visualization? (y/n)"
3. If you answer 'y', you'll see the map visualization menu
4. Select your preferred visualization method (1, 2, or 3)

## Code Structure

### Main Functions

- **`DisplayMap(data Response)`** - Main entry point for map visualization
  - Validates position data
  - Displays menu
  - Routes to selected visualization method

- **`displayASCIIMap(data Response)`** - Terminal ASCII map visualization
  - Creates 80x24 character grid
  - Plots positions on map
  - Displays legend and position details

- **`drawWorldMapOutline(grid [][]rune)`** - Draws continent outlines on ASCII grid

- **`exportToKML(data Response)`** - Exports to KML file
  - Prompts for file path
  - Generates KML content
  - Writes to file

- **`generateKMLContent(data Response) string`** - Generates KML XML content

- **`generateWebMap(data Response)`** - Exports to HTML file
  - Prompts for file path
  - Generates HTML content
  - Writes to file

- **`generateHTMLMapContent(data Response) string`** - Generates HTML with Leaflet.js

## Data Structures

The map visualization uses the `Response` type from `position.go`:

```go
type Response struct {
    SatelliteInfo SatelliteInfo `json:"info"`
    Positions    []Position    `json:"positions"`
}

type Position struct {
    Satlatitude  float64 `json:"satlatitude"`
    Satlongitude float64 `json:"satlongitude"`
    Sataltitude  float64 `json:"sataltitude"`
    Azimuth      float64 `json:"azimuth"`
    Elevation    float64 `json:"elevation"`
    Ra           float64 `json:"ra"`
    Dec          float64 `json:"dec"`
    Timestamp    int64   `json:"timestamp"`
}
```

## Testing

Comprehensive test coverage is provided in `osint/map_test.go`:

### Test Coverage

- **KML Generation Tests:**
  - `TestGenerateKMLContent` - Full KML structure validation
  - `TestGenerateKMLContentEmptyPositions` - Empty data handling
  - `TestGenerateKMLContentSinglePosition` - Single position handling
  - `TestGenerateKMLContentEdgeCases` - Extreme coordinates
  - `TestKMLFileExport` - File writing validation

- **HTML Generation Tests:**
  - `TestGenerateHTMLMapContent` - Full HTML structure validation
  - `TestGenerateHTMLMapContentEmptyPositions` - Empty data handling
  - `TestGenerateHTMLMapContentEdgeCases` - Extreme coordinates
  - `TestHTMLFileExport` - File writing validation

- **ASCII Map Tests:**
  - `TestDrawWorldMapOutline` - Map outline drawing
  - `TestDrawWorldMapOutlineSmallGrid` - Small grid handling
  - `TestDisplayASCIIMapPositionMapping` - Coordinate conversion
  - `TestDisplayMapEmptyPositions` - Empty data handling

### Running Tests

```bash
# Run all map visualization tests
go test ./osint -run "Test.*Map|Test.*KML|Test.*HTML|TestDrawWorld" -v

# Run specific test
go test ./osint -run TestGenerateKMLContent -v

# Run with coverage
go test ./osint -cover -run "Test.*Map|Test.*KML|Test.*HTML|TestDrawWorld"
```

## Examples

### Example 1: View ASCII Map

```
1. Get satellite position data
2. When prompted "View map visualization? (y/n)", enter 'y'
3. Select option 1 for Terminal ASCII Map
4. View the ASCII map in your terminal
```

### Example 2: Export to Google Earth

```
1. Get satellite position data
2. When prompted "View map visualization? (y/n)", enter 'y'
3. Select option 2 for KML Export
4. Enter file path (or press Enter for default)
5. Open the .kml file in Google Earth
```

### Example 3: Generate Interactive Web Map

```
1. Get satellite position data
2. When prompted "View map visualization? (y/n)", enter 'y'
3. Select option 3 for Web-based Interactive Map
4. Enter file path (or press Enter for default)
5. Open the .html file in your web browser
```

## Limitations

1. **ASCII Map:**
   - Simplified world map (not geographically accurate)
   - Fixed 80x24 character resolution
   - Limited to terminal display

2. **KML Export:**
   - Requires external application (Google Earth) to view
   - Paths are straight lines between positions (not orbital paths)

3. **Web Map:**
   - Requires internet connection for map tiles
   - Uses CDN for Leaflet.js (requires internet)
   - Large position datasets may slow down rendering

## Future Enhancements

Potential improvements for the map visualization module:

- [ ] Real-time position updates on web map
- [ ] Multiple satellite tracking on single map
- [ ] Ground track visualization
- [ ] 3D visualization option
- [ ] Custom map styles and themes
- [ ] Export to other formats (GeoJSON, GPX)
- [ ] Offline map tile support
- [ ] Animation of satellite movement
- [ ] Integration with orbital prediction data

## Troubleshooting

### ASCII Map Not Displaying Correctly

- Ensure your terminal supports Unicode characters
- Check terminal width is at least 80 characters
- Verify terminal supports color output

### KML File Won't Open

- Verify file has `.kml` extension
- Check file is valid XML (not corrupted)
- Ensure Google Earth or compatible viewer is installed
- Try opening in a text editor to verify content

### Web Map Not Loading

- Check internet connection (required for map tiles)
- Verify browser console for JavaScript errors
- Ensure Leaflet.js CDN is accessible
- Try a different browser
- Check browser allows loading local files

### No Positions Displayed

- Verify satellite position data was successfully fetched
- Check that `data.Positions` is not empty
- Ensure coordinates are valid (latitude: -90 to 90, longitude: -180 to 180)

## Contributing

When adding new features or fixing bugs:

1. Write tests for new functionality
2. Update this README with new features
3. Follow existing code style
4. Test all three visualization methods
5. Verify edge cases (empty data, extreme coordinates, etc.)

## License

This module is part of the SatIntel project. See the main project LICENSE file for details.

