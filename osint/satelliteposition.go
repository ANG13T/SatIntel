package osint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	fmt.Print("\n ENTER LATITUDE > ")
	var latitude string
	fmt.Scanln(&latitude)
	if strings.TrimSpace(latitude) == "" {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Latitude cannot be empty"))
		return
	}
	fmt.Print("\n ENTER LONGITUDE > ")
	var longitude string
	fmt.Scanln(&longitude)
	if strings.TrimSpace(longitude) == "" {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Longitude cannot be empty"))
		return
	}
	fmt.Print("\n ENTER ALTITUDE > ")
	var altitude string
	fmt.Scanln(&altitude)
	if strings.TrimSpace(altitude) == "" {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Altitude cannot be empty"))
		return
	}

	_, err := strconv.ParseFloat(latitude, 64)
	_, err2 := strconv.ParseFloat(longitude, 64)
	_, err3 := strconv.Atoi(altitude)

	if err != nil || err2 != nil || err3 != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: INVALID INPUT - Please enter valid numbers"))
		return
	}

	url := "https://api.n2yo.com/rest/v1/satellite/positions/" + norad + "/" + latitude + "/" + longitude + "/" + altitude + "/2/&apiKey=" + os.Getenv("N2YO_API_KEY")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to fetch satellite position data: "+err.Error()))
		return
	}
	defer resp.Body.Close()

	var data Response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to parse response: "+err.Error()))
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

// DisplayMap is a placeholder for future map visualization functionality.
func DisplayMap() {
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
