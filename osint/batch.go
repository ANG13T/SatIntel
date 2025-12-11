package osint

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/TwiN/go-color"
	"github.com/manifoldco/promptui"
)

// BatchSatellite represents a satellite selected for batch processing.
type BatchSatellite struct {
	Name      string
	NORADID   string
	Country   string
	ObjectType string
}

// BatchTLEResult contains TLE data for a satellite in batch processing.
type BatchTLEResult struct {
	Satellite BatchSatellite
	TLE       TLE
	Error     error
	Success   bool
}

// BatchComparisonResult contains comparison data for multiple satellites.
type BatchComparisonResult struct {
	Satellites []BatchSatellite
	Results    []BatchTLEResult
	Summary   BatchSummary
}

// BatchSummary provides aggregate statistics for batch operations.
type BatchSummary struct {
	TotalProcessed int
	Successful     int
	Failed         int
	AverageInclination float64
	AverageMeanMotion  float64
	LowestAltitude     float64
	HighestAltitude    float64
}

// showBatchMenu displays the batch operations menu.
func showBatchMenu() string {
	menuItems := []string{
		"Download TLE for Multiple Satellites",
		"Compare Multiple Satellites",
		"Batch Visual Predictions",
		"Batch Radio Predictions",
		"Batch Position Data",
		"Cancel",
	}

	prompt := promptui.Select{
		Label: "Batch Operations Menu",
		Items: menuItems,
		Size:  10,
	}

	idx, _, err := prompt.Run()
	if err != nil || idx == 5 {
		return ""
	}

	options := []string{"tle", "compare", "visual", "radio", "position"}
	if idx < len(options) {
		return options[idx]
	}
	return ""
}

// selectMultipleSatellites allows users to select multiple satellites for batch processing.
func selectMultipleSatellites() []BatchSatellite {
	var selected []BatchSatellite
	selectedMap := make(map[string]bool) // Track selected NORAD IDs to prevent duplicates

	for {
		fmt.Println(color.Ize(color.Cyan, "\n  [*] Current selection: "+strconv.Itoa(len(selected))+" satellite(s)"))
		
		menuItems := []string{
			"Add Satellite from Catalog",
			"Add Satellite by NORAD ID",
			"Add from Favorites",
			"Remove Satellite",
			"View Selected Satellites",
			"Clear All",
			"Done - Process Batch",
			"Cancel",
		}

		prompt := promptui.Select{
			Label: "Batch Selection Menu",
			Items: menuItems,
			Size:  10,
		}

		idx, _, err := prompt.Run()
		if err != nil || idx == 7 {
			return nil
		}

		switch idx {
		case 0: // Add from Catalog
			result := SelectSatellite()
			if result != "" {
				norad := extractNorad(result)
				if !selectedMap[norad] {
					// Fetch full satellite info
					client, err := Login()
					if err == nil {
						endpoint := fmt.Sprintf("/class/satcat/NORAD_CAT_ID/%s/format/json", norad)
						data, err := QuerySpaceTrack(client, endpoint)
						if err == nil {
							var sats []Satellite
							if json.Unmarshal([]byte(data), &sats) == nil && len(sats) > 0 {
								sat := sats[0]
								selected = append(selected, BatchSatellite{
									Name:       sat.SATNAME,
									NORADID:    sat.NORAD_CAT_ID,
									Country:    sat.COUNTRY,
									ObjectType: sat.OBJECT_TYPE,
								})
								selectedMap[norad] = true
								fmt.Println(color.Ize(color.Green, "  [+] Added: "+sat.SATNAME))
							} else {
								// Fallback with just name and NORAD
								name := strings.Split(result, " (")[0]
								selected = append(selected, BatchSatellite{
									Name:     name,
									NORADID:  norad,
									Country:  "Unknown",
									ObjectType: "Unknown",
								})
								selectedMap[norad] = true
								fmt.Println(color.Ize(color.Green, "  [+] Added: "+name))
							}
						} else {
							// Fallback
							name := strings.Split(result, " (")[0]
							selected = append(selected, BatchSatellite{
								Name:     name,
								NORADID:  norad,
								Country:  "Unknown",
								ObjectType: "Unknown",
							})
							selectedMap[norad] = true
							fmt.Println(color.Ize(color.Green, "  [+] Added: "+name))
						}
					} else {
						// Fallback
						name := strings.Split(result, " (")[0]
						selected = append(selected, BatchSatellite{
							Name:     name,
							NORADID:  norad,
							Country:  "Unknown",
							ObjectType: "Unknown",
						})
						selectedMap[norad] = true
						fmt.Println(color.Ize(color.Green, "  [+] Added: "+name))
					}
				} else {
					fmt.Println(color.Ize(color.Yellow, "  [!] Satellite already in batch"))
				}
			}

		case 1: // Add by NORAD ID
			noradPrompt := promptui.Prompt{
				Label:    "Enter NORAD ID",
				Validate: func(input string) error {
					if strings.TrimSpace(input) == "" {
						return fmt.Errorf("NORAD ID cannot be empty")
					}
					return nil
				},
			}
			norad, err := noradPrompt.Run()
			if err == nil && norad != "" {
				norad = strings.TrimSpace(norad)
				if !selectedMap[norad] {
					// Try to fetch satellite name
					client, err := Login()
					if err == nil {
						endpoint := fmt.Sprintf("/class/satcat/NORAD_CAT_ID/%s/format/json", norad)
						data, err := QuerySpaceTrack(client, endpoint)
						if err == nil {
							var sats []Satellite
							if json.Unmarshal([]byte(data), &sats) == nil && len(sats) > 0 {
								sat := sats[0]
								selected = append(selected, BatchSatellite{
									Name:       sat.SATNAME,
									NORADID:    sat.NORAD_CAT_ID,
									Country:    sat.COUNTRY,
									ObjectType: sat.OBJECT_TYPE,
								})
								selectedMap[norad] = true
								fmt.Println(color.Ize(color.Green, "  [+] Added: "+sat.SATNAME))
							} else {
								selected = append(selected, BatchSatellite{
									Name:     "NORAD " + norad,
									NORADID:  norad,
									Country:  "Unknown",
									ObjectType: "Unknown",
								})
								selectedMap[norad] = true
								fmt.Println(color.Ize(color.Green, "  [+] Added: NORAD "+norad))
							}
						} else {
							selected = append(selected, BatchSatellite{
								Name:     "NORAD " + norad,
								NORADID:  norad,
								Country:  "Unknown",
								ObjectType: "Unknown",
							})
							selectedMap[norad] = true
							fmt.Println(color.Ize(color.Green, "  [+] Added: NORAD "+norad))
						}
					} else {
						selected = append(selected, BatchSatellite{
							Name:     "NORAD " + norad,
							NORADID:  norad,
							Country:  "Unknown",
							ObjectType: "Unknown",
						})
						selectedMap[norad] = true
						fmt.Println(color.Ize(color.Green, "  [+] Added: NORAD "+norad))
					}
				} else {
					fmt.Println(color.Ize(color.Yellow, "  [!] Satellite already in batch"))
				}
			}

		case 2: // Add from Favorites
			favResult := SelectFromFavorites()
			if favResult != "" {
				norad := extractNorad(favResult)
				if !selectedMap[norad] {
					name := strings.Split(favResult, " (")[0]
					selected = append(selected, BatchSatellite{
						Name:     name,
						NORADID:  norad,
						Country:  "Unknown",
						ObjectType: "Unknown",
					})
					selectedMap[norad] = true
					fmt.Println(color.Ize(color.Green, "  [+] Added: "+name))
				} else {
					fmt.Println(color.Ize(color.Yellow, "  [!] Satellite already in batch"))
				}
			}

		case 3: // Remove Satellite
			if len(selected) == 0 {
				fmt.Println(color.Ize(color.Yellow, "  [!] No satellites to remove"))
				continue
			}
			var items []string
			for i, sat := range selected {
				items = append(items, fmt.Sprintf("%d. %s (%s)", i+1, sat.Name, sat.NORADID))
			}
			items = append(items, "Cancel")

			removePrompt := promptui.Select{
				Label: "Select satellite to remove",
				Items: items,
			}
			removeIdx, _, err := removePrompt.Run()
			if err == nil && removeIdx < len(selected) {
				removed := selected[removeIdx]
				selected = append(selected[:removeIdx], selected[removeIdx+1:]...)
				delete(selectedMap, removed.NORADID)
				fmt.Println(color.Ize(color.Green, "  [+] Removed: "+removed.Name))
			}

		case 4: // View Selected
			if len(selected) == 0 {
				fmt.Println(color.Ize(color.Yellow, "  [!] No satellites selected"))
			} else {
				fmt.Println(color.Ize(color.Cyan, "\n  Selected Satellites:"))
				for i, sat := range selected {
					fmt.Printf("  %d. %s (%s)\n", i+1, sat.Name, sat.NORADID)
					if sat.Country != "Unknown" {
						fmt.Printf("     Country: %s\n", sat.Country)
					}
					if sat.ObjectType != "Unknown" {
						fmt.Printf("     Type: %s\n", sat.ObjectType)
					}
				}
			}

		case 5: // Clear All
			if len(selected) > 0 {
				confirmPrompt := promptui.Prompt{
					Label:     "Clear all satellites? (y/n)",
					Default:   "n",
					AllowEdit: true,
				}
				confirm, _ := confirmPrompt.Run()
				if strings.ToLower(strings.TrimSpace(confirm)) == "y" {
					selected = []BatchSatellite{}
					selectedMap = make(map[string]bool)
					fmt.Println(color.Ize(color.Green, "  [+] Cleared all satellites"))
				}
			}

		case 6: // Done
			if len(selected) == 0 {
				fmt.Println(color.Ize(color.Red, "  [!] Please select at least one satellite"))
				continue
			}
			return selected
		}
	}
}

// BatchDownloadTLE downloads TLE data for multiple satellites concurrently.
func BatchDownloadTLE(satellites []BatchSatellite) []BatchTLEResult {
	if len(satellites) == 0 {
		return nil
	}

	fmt.Println(color.Ize(color.Cyan, fmt.Sprintf("\n  [*] Downloading TLE data for %d satellite(s)...", len(satellites))))

	client, err := Login()
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to login: "+err.Error()))
		return nil
	}

	results := make([]BatchTLEResult, len(satellites))
	var wg sync.WaitGroup
	var mu sync.Mutex
	completed := 0

	for i, sat := range satellites {
		wg.Add(1)
		go func(idx int, satellite BatchSatellite) {
			defer wg.Done()

			result := BatchTLEResult{
				Satellite: satellite,
				Success:   false,
			}

			endpoint := fmt.Sprintf("/class/gp_history/format/tle/NORAD_CAT_ID/%s/orderby/EPOCH%%20desc/limit/1", satellite.NORADID)
			data, err := QuerySpaceTrack(client, endpoint)
			if err != nil {
				result.Error = err
				mu.Lock()
				results[idx] = result
				completed++
				mu.Unlock()
				return
			}

			lines := strings.Split(strings.TrimSpace(data), "\n")
			var lineOne, lineTwo string

			if len(lines) >= 2 {
				lineOne = strings.TrimSpace(lines[0])
				lineTwo = strings.TrimSpace(lines[1])
			} else {
				tleLines := strings.Fields(data)
				if len(tleLines) >= 2 {
					mid := len(tleLines) / 2
					if mid < 1 {
						mid = 1
					}
					if mid >= len(tleLines) {
						mid = len(tleLines) - 1
					}
					lineOne = strings.Join(tleLines[:mid], " ")
					lineTwo = strings.Join(tleLines[mid:], " ")
				} else {
					result.Error = fmt.Errorf("insufficient TLE data")
					mu.Lock()
					results[idx] = result
					completed++
					mu.Unlock()
					return
				}
			}

			if !strings.HasPrefix(lineOne, "1 ") || !strings.HasPrefix(lineTwo, "2 ") {
				result.Error = fmt.Errorf("invalid TLE format")
				mu.Lock()
				results[idx] = result
				completed++
				mu.Unlock()
				return
			}

			tle := ConstructTLE(satellite.Name, lineOne, lineTwo)
			
			// Validate parsing
			line1Fields := strings.Fields(lineOne)
			line2Fields := strings.Fields(lineTwo)
			if len(line1Fields) < 4 || len(line2Fields) < 3 {
				result.Error = fmt.Errorf("insufficient fields in TLE")
				mu.Lock()
				results[idx] = result
				completed++
				mu.Unlock()
				return
			}

			if tle.SatelliteCatalogNumber == 0 && tle.InternationalDesignator == "" && tle.ElementSetEpoch == 0.0 {
				result.Error = fmt.Errorf("failed to parse TLE data")
				mu.Lock()
				results[idx] = result
				completed++
				mu.Unlock()
				return
			}

			result.TLE = tle
			result.Success = true

			mu.Lock()
			results[idx] = result
			completed++
			mu.Unlock()

			fmt.Printf(color.Ize(color.Green, "  [+] [%d/%d] Downloaded: %s\n"), completed, len(satellites), satellite.Name)
		}(i, sat)
	}

	wg.Wait()

	// Display summary
	successful := 0
	for _, r := range results {
		if r.Success {
			successful++
		}
	}

	fmt.Println(color.Ize(color.Cyan, fmt.Sprintf("\n  [*] Batch download complete: %d/%d successful", successful, len(satellites))))

	return results
}

// CompareSatellites compares TLE data for multiple satellites and displays a summary.
func CompareSatellites(results []BatchTLEResult) BatchComparisonResult {
	if len(results) == 0 {
		return BatchComparisonResult{}
	}

	comparison := BatchComparisonResult{
		Results: results,
		Summary: BatchSummary{
			TotalProcessed: len(results),
		},
	}

	var inclinations []float64
	var meanMotions []float64
	var altitudes []float64

	for _, result := range results {
		if result.Success {
			comparison.Summary.Successful++
			comparison.Satellites = append(comparison.Satellites, result.Satellite)

			if result.TLE.OrbitInclination > 0 {
				inclinations = append(inclinations, result.TLE.OrbitInclination)
			}
			if result.TLE.MeanMotion > 0 {
				meanMotions = append(meanMotions, result.TLE.MeanMotion)
			}
			// Estimate altitude from mean motion (rough calculation)
			if result.TLE.MeanMotion > 0 {
				// Simplified: altitude ≈ (GM / (4π² * meanMotion²))^(1/3) - Earth radius
				// For simplicity, using approximate formula
				period := 86400.0 / result.TLE.MeanMotion // period in seconds
				// Rough altitude estimate (km)
				altitude := (398600.4418 * period * period / (4 * 3.14159 * 3.14159)) / 1000 - 6371
				if altitude > 0 {
					altitudes = append(altitudes, altitude)
				}
			}
		} else {
			comparison.Summary.Failed++
		}
	}

	// Calculate averages
	if len(inclinations) > 0 {
		sum := 0.0
		for _, inc := range inclinations {
			sum += inc
		}
		comparison.Summary.AverageInclination = sum / float64(len(inclinations))
	}

	if len(meanMotions) > 0 {
		sum := 0.0
		for _, mm := range meanMotions {
			sum += mm
		}
		comparison.Summary.AverageMeanMotion = sum / float64(len(meanMotions))
	}

	if len(altitudes) > 0 {
		min := altitudes[0]
		max := altitudes[0]
		for _, alt := range altitudes {
			if alt < min {
				min = alt
			}
			if alt > max {
				max = alt
			}
		}
		comparison.Summary.LowestAltitude = min
		comparison.Summary.HighestAltitude = max
	}

	return comparison
}

// DisplayComparison displays the comparison results in a formatted table.
func DisplayComparison(comparison BatchComparisonResult) {
	if len(comparison.Results) == 0 {
		fmt.Println(color.Ize(color.Yellow, "  [!] No data to compare"))
		return
	}

	fmt.Println(color.Ize(color.Purple, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Purple, "║              Satellite Comparison Summary                  ║"))
	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))

	fmt.Println(color.Ize(color.Purple, GenRowString("Total Processed", strconv.Itoa(comparison.Summary.TotalProcessed))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Successful", strconv.Itoa(comparison.Summary.Successful))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Failed", strconv.Itoa(comparison.Summary.Failed))))

	if comparison.Summary.AverageInclination > 0 {
		fmt.Println(color.Ize(color.Purple, GenRowString("Average Inclination", fmt.Sprintf("%.2f°", comparison.Summary.AverageInclination))))
	}
	if comparison.Summary.AverageMeanMotion > 0 {
		fmt.Println(color.Ize(color.Purple, GenRowString("Average Mean Motion", fmt.Sprintf("%.4f rev/day", comparison.Summary.AverageMeanMotion))))
	}
	if comparison.Summary.LowestAltitude > 0 {
		fmt.Println(color.Ize(color.Purple, GenRowString("Lowest Altitude (est.)", fmt.Sprintf("%.2f km", comparison.Summary.LowestAltitude))))
	}
	if comparison.Summary.HighestAltitude > 0 {
		fmt.Println(color.Ize(color.Purple, GenRowString("Highest Altitude (est.)", fmt.Sprintf("%.2f km", comparison.Summary.HighestAltitude))))
	}

	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
	fmt.Println(color.Ize(color.Purple, "║                    Individual Results                     ║"))
	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))

	for i, result := range comparison.Results {
		status := "✅ Success"
		if !result.Success {
			status = "❌ Failed"
			if result.Error != nil {
				status += ": " + result.Error.Error()
			}
		}

		fmt.Println(color.Ize(color.Purple, GenRowString(fmt.Sprintf("Satellite %d", i+1), result.Satellite.Name)))
		fmt.Println(color.Ize(color.Purple, GenRowString("  NORAD ID", result.Satellite.NORADID)))
		fmt.Println(color.Ize(color.Purple, GenRowString("  Status", status)))

		if result.Success {
			fmt.Println(color.Ize(color.Purple, GenRowString("  Inclination", fmt.Sprintf("%.2f°", result.TLE.OrbitInclination))))
			fmt.Println(color.Ize(color.Purple, GenRowString("  Mean Motion", fmt.Sprintf("%.4f rev/day", result.TLE.MeanMotion))))
			fmt.Println(color.Ize(color.Purple, GenRowString("  Eccentricity", fmt.Sprintf("%.6f", result.TLE.Eccentrcity))))
		}

		if i < len(comparison.Results)-1 {
			fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
		}
	}

	fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝\n\n"))
}

// BatchOperations provides the main entry point for batch operations.
func BatchOperations() {
	operation := showBatchMenu()
	if operation == "" {
		return
	}

	satellites := selectMultipleSatellites()
	if len(satellites) == 0 {
		return
	}

	switch operation {
	case "tle":
		results := BatchDownloadTLE(satellites)
		if len(results) > 0 {
			// Display results
			fmt.Println(color.Ize(color.Cyan, "\n  [*] Batch TLE Download Results:"))
			for i, result := range results {
				if result.Success {
					fmt.Printf("\n  %d. %s (%s) - ✅ Success\n", i+1, result.Satellite.Name, result.Satellite.NORADID)
					PrintTLE(result.TLE)
				} else {
					fmt.Printf("\n  %d. %s (%s) - ❌ Failed", i+1, result.Satellite.Name, result.Satellite.NORADID)
					if result.Error != nil {
						fmt.Printf(": %s\n", result.Error.Error())
					} else {
						fmt.Println()
					}
				}
			}

			// Offer export
			exportPrompt := promptui.Prompt{
				Label:     "Export batch results? (y/n)",
				Default:   "n",
				AllowEdit: true,
			}
			exportAnswer, _ := exportPrompt.Run()
			if strings.ToLower(strings.TrimSpace(exportAnswer)) == "y" {
				exportBatchTLE(results)
			}
		}

	case "compare":
		results := BatchDownloadTLE(satellites)
		if len(results) > 0 {
			comparison := CompareSatellites(results)
			DisplayComparison(comparison)

			// Offer export
			exportPrompt := promptui.Prompt{
				Label:     "Export comparison results? (y/n)",
				Default:   "n",
				AllowEdit: true,
			}
			exportAnswer, _ := exportPrompt.Run()
			if strings.ToLower(strings.TrimSpace(exportAnswer)) == "y" {
				exportBatchComparison(comparison)
			}
		}

	case "visual", "radio", "position":
		fmt.Println(color.Ize(color.Yellow, "  [!] Batch predictions and positions coming soon"))
		// TODO: Implement batch predictions and positions
	}
}

// exportBatchTLE exports batch TLE results to a file.
func exportBatchTLE(results []BatchTLEResult) {
	formatPrompt := promptui.Select{
		Label: "Select Export Format",
		Items: []string{"CSV", "JSON", "Text", "Cancel"},
	}
	formatIdx, formatChoice, err := formatPrompt.Run()
	if err != nil || formatIdx == 3 {
		return
	}

	pathPrompt := promptui.Prompt{
		Label:    "Enter file path",
		Default: fmt.Sprintf("batch_tle_%s", time.Now().Format("20060102_150405")),
		AllowEdit: true,
	}
	filePath, err := pathPrompt.Run()
	if err != nil {
		return
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		filePath = fmt.Sprintf("batch_tle_%s", time.Now().Format("20060102_150405"))
	}

	ext := ""
	switch formatChoice {
	case "CSV":
		ext = ".csv"
	case "JSON":
		ext = ".json"
	case "Text":
		ext = ".txt"
	}

	if !strings.HasSuffix(filePath, ext) {
		filePath += ext
	}

	switch formatChoice {
	case "CSV":
		exportBatchTLECSV(results, filePath)
	case "JSON":
		exportBatchTLEJSON(results, filePath)
	case "Text":
		exportBatchTLEText(results, filePath)
	}
}

// exportBatchComparison exports comparison results to a file.
func exportBatchComparison(comparison BatchComparisonResult) {
	formatPrompt := promptui.Select{
		Label: "Select Export Format",
		Items: []string{"CSV", "JSON", "Text", "Cancel"},
	}
	formatIdx, formatChoice, err := formatPrompt.Run()
	if err != nil || formatIdx == 3 {
		return
	}

	pathPrompt := promptui.Prompt{
		Label:    "Enter file path",
		Default: fmt.Sprintf("batch_comparison_%s", time.Now().Format("20060102_150405")),
		AllowEdit: true,
	}
	filePath, err := pathPrompt.Run()
	if err != nil {
		return
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		filePath = fmt.Sprintf("batch_comparison_%s", time.Now().Format("20060102_150405"))
	}

	ext := ""
	switch formatChoice {
	case "CSV":
		ext = ".csv"
	case "JSON":
		ext = ".json"
	case "Text":
		ext = ".txt"
	}

	if !strings.HasSuffix(filePath, ext) {
		filePath += ext
	}

	switch formatChoice {
	case "CSV":
		exportBatchComparisonCSV(comparison, filePath)
	case "JSON":
		exportBatchComparisonJSON(comparison, filePath)
	case "Text":
		exportBatchComparisonText(comparison, filePath)
	}
}

// exportBatchTLECSV exports batch TLE results to CSV format.
func exportBatchTLECSV(results []BatchTLEResult, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	headers := []string{
		"Satellite Name", "NORAD ID", "Status", "Common Name", "Catalog Number",
		"Inclination", "Mean Motion", "Eccentricity", "Error",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for _, result := range results {
		status := "Success"
		errorMsg := ""
		if !result.Success {
			status = "Failed"
			if result.Error != nil {
				errorMsg = result.Error.Error()
			}
		}

		row := []string{
			result.Satellite.Name,
			result.Satellite.NORADID,
			status,
		}

		if result.Success {
			row = append(row,
				result.TLE.CommonName,
				strconv.Itoa(result.TLE.SatelliteCatalogNumber),
				fmt.Sprintf("%.2f", result.TLE.OrbitInclination),
				fmt.Sprintf("%.4f", result.TLE.MeanMotion),
				fmt.Sprintf("%.6f", result.TLE.Eccentrcity),
				errorMsg,
			)
		} else {
			row = append(row, "", "", "", "", "", errorMsg)
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	fmt.Println(color.Ize(color.Green, "  [+] Exported to: "+filePath))
	return nil
}

// exportBatchTLEJSON exports batch TLE results to JSON format.
func exportBatchTLEJSON(results []BatchTLEResult, filePath string) error {
	type ExportResult struct {
		Satellite BatchSatellite `json:"satellite"`
		TLE       *TLE           `json:"tle,omitempty"`
		Success   bool           `json:"success"`
		Error     string         `json:"error,omitempty"`
	}

	var exportResults []ExportResult
	for _, result := range results {
		exportResult := ExportResult{
			Satellite: result.Satellite,
			Success:   result.Success,
		}
		if result.Success {
			exportResult.TLE = &result.TLE
		}
		if result.Error != nil {
			exportResult.Error = result.Error.Error()
		}
		exportResults = append(exportResults, exportResult)
	}

	data := map[string]interface{}{
		"batch_results": exportResults,
		"export_timestamp": time.Now().Format(time.RFC3339),
		"total_count": len(results),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println(color.Ize(color.Green, "  [+] Exported to: "+filePath))
	return nil
}

// exportBatchTLEText exports batch TLE results to text format.
func exportBatchTLEText(results []BatchTLEResult, filePath string) error {
	var builder strings.Builder

	builder.WriteString("Batch TLE Download Results\n")
	builder.WriteString(strings.Repeat("=", 60) + "\n\n")
	builder.WriteString(fmt.Sprintf("Total Satellites: %d\n", len(results)))
	builder.WriteString(fmt.Sprintf("Export Date: %s\n\n", time.Now().Format(time.RFC3339)))

	for i, result := range results {
		builder.WriteString(fmt.Sprintf("Satellite %d: %s (%s)\n", i+1, result.Satellite.Name, result.Satellite.NORADID))
		if result.Success {
			builder.WriteString("Status: ✅ Success\n")
			builder.WriteString(fmt.Sprintf("  Common Name: %s\n", result.TLE.CommonName))
			builder.WriteString(fmt.Sprintf("  Catalog Number: %d\n", result.TLE.SatelliteCatalogNumber))
			builder.WriteString(fmt.Sprintf("  Inclination: %.2f°\n", result.TLE.OrbitInclination))
			builder.WriteString(fmt.Sprintf("  Mean Motion: %.4f rev/day\n", result.TLE.MeanMotion))
			builder.WriteString(fmt.Sprintf("  Eccentricity: %.6f\n", result.TLE.Eccentrcity))
		} else {
			builder.WriteString("Status: ❌ Failed\n")
			if result.Error != nil {
				builder.WriteString(fmt.Sprintf("  Error: %s\n", result.Error.Error()))
			}
		}
		builder.WriteString("\n")
	}

	if err := os.WriteFile(filePath, []byte(builder.String()), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println(color.Ize(color.Green, "  [+] Exported to: "+filePath))
	return nil
}

// exportBatchComparisonCSV exports comparison results to CSV format.
func exportBatchComparisonCSV(comparison BatchComparisonResult, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write summary
	summaryHeaders := []string{"Metric", "Value"}
	if err := writer.Write(summaryHeaders); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	summaryRows := [][]string{
		{"Total Processed", strconv.Itoa(comparison.Summary.TotalProcessed)},
		{"Successful", strconv.Itoa(comparison.Summary.Successful)},
		{"Failed", strconv.Itoa(comparison.Summary.Failed)},
	}
	if comparison.Summary.AverageInclination > 0 {
		summaryRows = append(summaryRows, []string{"Average Inclination", fmt.Sprintf("%.2f", comparison.Summary.AverageInclination)})
	}
	if comparison.Summary.AverageMeanMotion > 0 {
		summaryRows = append(summaryRows, []string{"Average Mean Motion", fmt.Sprintf("%.4f", comparison.Summary.AverageMeanMotion)})
	}

	for _, row := range summaryRows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write summary row: %w", err)
		}
	}

	// Empty row separator
	writer.Write([]string{})

	// Write individual results
	resultHeaders := []string{"Name", "NORAD ID", "Status", "Inclination", "Mean Motion", "Eccentricity", "Error"}
	if err := writer.Write(resultHeaders); err != nil {
		return fmt.Errorf("failed to write result header: %w", err)
	}

	for _, result := range comparison.Results {
		status := "Success"
		errorMsg := ""
		if !result.Success {
			status = "Failed"
			if result.Error != nil {
				errorMsg = result.Error.Error()
			}
		}

		row := []string{
			result.Satellite.Name,
			result.Satellite.NORADID,
			status,
		}

		if result.Success {
			row = append(row,
				fmt.Sprintf("%.2f", result.TLE.OrbitInclination),
				fmt.Sprintf("%.4f", result.TLE.MeanMotion),
				fmt.Sprintf("%.6f", result.TLE.Eccentrcity),
				errorMsg,
			)
		} else {
			row = append(row, "", "", errorMsg)
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write result row: %w", err)
		}
	}

	fmt.Println(color.Ize(color.Green, "  [+] Exported to: "+filePath))
	return nil
}

// exportBatchComparisonJSON exports comparison results to JSON format.
func exportBatchComparisonJSON(comparison BatchComparisonResult, filePath string) error {
	data := map[string]interface{}{
		"summary": comparison.Summary,
		"results": comparison.Results,
		"export_timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println(color.Ize(color.Green, "  [+] Exported to: "+filePath))
	return nil
}

// exportBatchComparisonText exports comparison results to text format.
func exportBatchComparisonText(comparison BatchComparisonResult, filePath string) error {
	var builder strings.Builder

	builder.WriteString("Satellite Comparison Results\n")
	builder.WriteString(strings.Repeat("=", 60) + "\n\n")

	builder.WriteString("Summary:\n")
	builder.WriteString(fmt.Sprintf("  Total Processed: %d\n", comparison.Summary.TotalProcessed))
	builder.WriteString(fmt.Sprintf("  Successful: %d\n", comparison.Summary.Successful))
	builder.WriteString(fmt.Sprintf("  Failed: %d\n", comparison.Summary.Failed))
	if comparison.Summary.AverageInclination > 0 {
		builder.WriteString(fmt.Sprintf("  Average Inclination: %.2f°\n", comparison.Summary.AverageInclination))
	}
	if comparison.Summary.AverageMeanMotion > 0 {
		builder.WriteString(fmt.Sprintf("  Average Mean Motion: %.4f rev/day\n", comparison.Summary.AverageMeanMotion))
	}

	builder.WriteString("\nIndividual Results:\n")
	builder.WriteString(strings.Repeat("-", 60) + "\n")

	for i, result := range comparison.Results {
		builder.WriteString(fmt.Sprintf("\n%d. %s (%s)\n", i+1, result.Satellite.Name, result.Satellite.NORADID))
		if result.Success {
			builder.WriteString("  Status: ✅ Success\n")
			builder.WriteString(fmt.Sprintf("  Inclination: %.2f°\n", result.TLE.OrbitInclination))
			builder.WriteString(fmt.Sprintf("  Mean Motion: %.4f rev/day\n", result.TLE.MeanMotion))
			builder.WriteString(fmt.Sprintf("  Eccentricity: %.6f\n", result.TLE.Eccentrcity))
		} else {
			builder.WriteString("  Status: ❌ Failed\n")
			if result.Error != nil {
				builder.WriteString(fmt.Sprintf("  Error: %s\n", result.Error.Error()))
			}
		}
	}

	builder.WriteString(fmt.Sprintf("\nExported: %s\n", time.Now().Format(time.RFC3339)))

	if err := os.WriteFile(filePath, []byte(builder.String()), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println(color.Ize(color.Green, "  [+] Exported to: "+filePath))
	return nil
}

