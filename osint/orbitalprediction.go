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

// OrbitalPrediction provides an interactive menu for visual and radio pass predictions.
func OrbitalPrediction() {
	options, _ := os.ReadFile("txt/orbital_prediction.txt")
	opt, _ := gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 3)

	if selection == 1 {
		GetVisualPrediction()
	} else if selection == 2 {
		GetRadioPrediction()
	}
}

// GetVisualPrediction fetches and displays visual pass predictions for a satellite.
func GetVisualPrediction() {
	selection := SatelliteSelection()
	if selection.norad == "" {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: INVALID INPUT"))
		return
	}
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
	fmt.Print("\n ENTER DAYS OF PREDICTION > ")
	var days string
	fmt.Scanln(&days)
	if strings.TrimSpace(days) == "" {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Days cannot be empty"))
		return
	}
	fmt.Print("\n ENTER MIN VISIBILITY > ")
	var vis string
	fmt.Scanln(&vis)
	if strings.TrimSpace(vis) == "" {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Minimum visibility cannot be empty"))
		return
	}

	_, err := strconv.ParseFloat(latitude, 64)
	_, err2 := strconv.ParseFloat(longitude, 64)
	_, err3 := strconv.ParseFloat(altitude, 64)
	_, err4 := strconv.Atoi(days)
	_, err5 := strconv.Atoi(vis)

	if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: INVALID INPUT - Please enter valid numbers"))
		return
	}

	url := "https://api.n2yo.com/rest/v1/satellite/visualpasses/" + selection.norad + "/" + latitude + "/" + longitude + "/" + altitude + "/" + days + "/" + vis + "/&apiKey=" + os.Getenv("N2YO_API_KEY")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to fetch visual pass data: "+err.Error()))
		return
	}
	defer resp.Body.Close()

	var data VisualPassesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to parse response: "+err.Error()))
		return
	}

	fmt.Println(color.Ize(color.Purple, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Purple, "║                    Satellite Information                    ║"))
	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))

	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Name", data.Info.SatName)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite ID", fmt.Sprintf("%d", data.Info.SatID))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Transactions Count", fmt.Sprintf("%d", data.Info.TransactionsCount))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Passes Count", fmt.Sprintf("%d", data.Info.PassesCount))))

	if len(data.Passes) > 0 {
		fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
		fmt.Println(color.Ize(color.Purple, "║                       Satellite Passes                      ║"))
		fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))

		for in, pos := range data.Passes {
			PrintVisualPass(pos, in == len(data.Passes)-1)
		}
	} else {
		fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝\n\n"))
	}

	// Offer export option
	exportPrompt := promptui.Prompt{
		Label:     "Export visual pass predictions? (y/n)",
		Default:   "n",
		AllowEdit: true,
	}
	exportAnswer, _ := exportPrompt.Run()
	if strings.ToLower(strings.TrimSpace(exportAnswer)) == "y" {
		defaultFilename := fmt.Sprintf("visual_passes_%s_%d", strings.ReplaceAll(data.Info.SatName, " ", "_"), data.Info.SatID)
		format, filePath, err := showExportMenu(defaultFilename)
		if err == nil {
			if err := ExportVisualPrediction(data, format, filePath); err != nil {
				fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to export: "+err.Error()))
			} else {
				fmt.Println(color.Ize(color.Green, fmt.Sprintf("  [+] Exported to: %s", filePath)))
			}
		}
	}
}

// GetRadioPrediction fetches and displays radio pass predictions for a satellite.
func GetRadioPrediction() {
	selection := SatelliteSelection()
	if selection.norad == "" {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: INVALID INPUT"))
		return
	}
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
	fmt.Print("\n ENTER DAYS OF PREDICTION > ")
	var days string
	fmt.Scanln(&days)
	if strings.TrimSpace(days) == "" {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Days cannot be empty"))
		return
	}
	fmt.Print("\n ENTER MIN ELEVATION > ")
	var elevation string
	fmt.Scanln(&elevation)
	if strings.TrimSpace(elevation) == "" {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Minimum elevation cannot be empty"))
		return
	}

	_, err := strconv.ParseFloat(latitude, 64)
	_, err2 := strconv.ParseFloat(longitude, 64)
	_, err3 := strconv.ParseFloat(altitude, 64)
	_, err4 := strconv.Atoi(days)
	_, err5 := strconv.Atoi(elevation)

	if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: INVALID INPUT - Please enter valid numbers"))
		return
	}

	url := "https://api.n2yo.com/rest/v1/satellite/radiopasses/" + selection.norad + "/" + latitude + "/" + longitude + "/" + altitude + "/" + days + "/" + elevation + "/&apiKey=" + os.Getenv("N2YO_API_KEY")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to fetch radio pass data: "+err.Error()))
		return
	}
	defer resp.Body.Close()

	var data RadioPassResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to parse response: "+err.Error()))
		return
	}

	fmt.Println(color.Ize(color.Purple, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Purple, "║                    Satellite Information                    ║"))
	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))

	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Name", data.Info.SatName)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite ID", fmt.Sprintf("%d", data.Info.SatID))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Transactions Count", fmt.Sprintf("%d", data.Info.TransactionsCount))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Passes Count", fmt.Sprintf("%d", data.Info.PassesCount))))

	if len(data.Passes) > 0 {
		fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
		fmt.Println(color.Ize(color.Purple, "║                       Satellite Passes                      ║"))
		fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))

		for in, pos := range data.Passes {
			PrintRadioPass(pos, in == len(data.Passes)-1)
		}
	} else {
		fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝\n\n"))
	}

	// Offer export option
	exportPrompt := promptui.Prompt{
		Label:     "Export radio pass predictions? (y/n)",
		Default:   "n",
		AllowEdit: true,
	}
	exportAnswer, _ := exportPrompt.Run()
	if strings.ToLower(strings.TrimSpace(exportAnswer)) == "y" {
		defaultFilename := fmt.Sprintf("radio_passes_%s_%d", strings.ReplaceAll(data.Info.SatName, " ", "_"), data.Info.SatID)
		format, filePath, err := showExportMenu(defaultFilename)
		if err == nil {
			if err := ExportRadioPrediction(data, format, filePath); err != nil {
				fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to export: "+err.Error()))
			} else {
				fmt.Println(color.Ize(color.Green, fmt.Sprintf("  [+] Exported to: %s", filePath)))
			}
		}
	}
}

// SatelliteSelection provides an interactive menu for selecting a satellite by catalog or NORAD ID.
func SatelliteSelection() SatelliteSelectionType {
	options, _ := os.ReadFile("txt/orbital_element.txt")
	opt, _ := gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 3)
	if selection == 1 {
		result := SelectSatellite()

		if result == "" {
			return SatelliteSelectionType{}
		}

		return SatelliteSelectionType{norad: extractNorad(result), name: result}

	} else if selection == 2 {
		fmt.Print("\n ENTER NORAD ID > ")
		var norad string
		fmt.Scanln(&norad)
		return SatelliteSelectionType{norad: norad, name: "UNSPECIFIED"}
	}

	return SatelliteSelectionType{}
}

type SatelliteSelectionType struct {
	norad string
	name  string
}

type VisualPassesResponse struct {
	Info   Info   `json:"info"`
	Passes []Pass `json:"passes"`
}

type Info struct {
	SatID             int    `json:"satid"`
	SatName           string `json:"satname"`
	TransactionsCount int    `json:"transactionscount"`
	PassesCount       int    `json:"passescount"`
}

type Pass struct {
	StartAz        float64 `json:"startAz"`
	StartAzCompass string  `json:"startAzCompass"`
	StartEl        float64 `json:"startEl"`
	StartUTC       int     `json:"startUTC"`
	MaxAz          float64 `json:"maxAz"`
	MaxAzCompass   string  `json:"maxAzCompass"`
	MaxEl          float64 `json:"maxEl"`
	MaxUTC         int     `json:"maxUTC"`
	EndAz          float64 `json:"endAz"`
	EndAzCompass   string  `json:"endAzCompass"`
	EndEl          float64 `json:"endEl"`
	EndUTC         int     `json:"endUTC"`
	Mag            float64 `json:"mag"`
	Duration       int     `json:"duration"`
}

type RadioPass struct {
	StartAz        float64 `json:"startAz"`
	StartAzCompass string  `json:"startAzCompass"`
	StartUTC       int64   `json:"startUTC"`
	MaxAz          float64 `json:"maxAz"`
	MaxAzCompass   string  `json:"maxAzCompass"`
	MaxEl          float64 `json:"maxEl"`
	MaxUTC         int64   `json:"maxUTC"`
	EndAz          float64 `json:"endAz"`
	EndAzCompass   string  `json:"endAzCompass"`
	EndUTC         int64   `json:"endUTC"`
}

type RadioPassResponse struct {
	Info   Info       `json:"info"`
	Passes []RadioPass `json:"passes"`
}

// PrintVisualPass displays visual pass information in a formatted table.
func PrintVisualPass(pass Pass, last bool) {
	fmt.Println(color.Ize(color.Purple, GenRowString("Start Azimuth", fmt.Sprintf("%f", pass.StartAz))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Start Azimuth Compass", pass.StartAzCompass)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Start Elevation", fmt.Sprintf("%f", pass.StartEl))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Start UTC", fmt.Sprintf("%d", pass.StartUTC))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Azimuth for Max Elevation", fmt.Sprintf("%f", pass.MaxAz))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Azimuth Compass for Max Elevation", pass.MaxAzCompass)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Max Elevation", fmt.Sprintf("%f", pass.MaxEl))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Max UTC", fmt.Sprintf("%d", pass.MaxUTC))))
	fmt.Println(color.Ize(color.Purple, GenRowString("End Azimuth", fmt.Sprintf("%f", pass.EndAz))))
	fmt.Println(color.Ize(color.Purple, GenRowString("End Azimuth Compass", pass.EndAzCompass)))
	fmt.Println(color.Ize(color.Purple, GenRowString("End Elevation", fmt.Sprintf("%f", pass.EndEl))))
	fmt.Println(color.Ize(color.Purple, GenRowString("End UTC", fmt.Sprintf("%d", pass.EndUTC))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Max Visual Magnitude", fmt.Sprintf("%f", pass.Mag))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Visible Duration", fmt.Sprintf("%d", pass.Duration))))
	if last {
		fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝\n\n"))
	} else {
		fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
	}
}

// PrintRadioPass displays radio pass information in a formatted table.
func PrintRadioPass(pass RadioPass, last bool) {
	fmt.Println(color.Ize(color.Purple, GenRowString("Start Azimuth", fmt.Sprintf("%f", pass.StartAz))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Start Azimuth Compass", pass.StartAzCompass)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Start UTC", fmt.Sprintf("%d", pass.StartUTC))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Azimuth for Max Elevation", fmt.Sprintf("%f", pass.MaxAz))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Azimuth Compass for Max Elevation", pass.MaxAzCompass)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Max Elevation", fmt.Sprintf("%f", pass.MaxEl))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Max UTC", fmt.Sprintf("%d", pass.MaxUTC))))
	fmt.Println(color.Ize(color.Purple, GenRowString("End Azimuth", fmt.Sprintf("%f", pass.EndAz))))
	fmt.Println(color.Ize(color.Purple, GenRowString("End Azimuth Compass", pass.EndAzCompass)))
	fmt.Println(color.Ize(color.Purple, GenRowString("End UTC", fmt.Sprintf("%d", pass.EndUTC))))
	if last {
		fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝\n\n"))
	} else {
		fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
	}
}
