package osint

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

// ExportFormat represents the available export formats.
type ExportFormat string

const (
	FormatCSV  ExportFormat = "CSV"
	FormatJSON ExportFormat = "JSON"
	FormatText ExportFormat = "Text"
)

// showExportMenu displays a menu for selecting export format and file path.
func showExportMenu(defaultFilename string) (ExportFormat, string, error) {
	formatItems := []string{"CSV", "JSON", "Text", "Cancel"}
	
	formatPrompt := promptui.Select{
		Label: "Select Export Format",
		Items: formatItems,
	}

	formatIdx, formatChoice, err := formatPrompt.Run()
	if err != nil || formatIdx == 3 {
		return "", "", fmt.Errorf("export cancelled")
	}

	format := ExportFormat(formatChoice)

	pathPrompt := promptui.Prompt{
		Label:    "Enter file path (or press Enter for default)",
		Default:  defaultFilename,
		AllowEdit: true,
	}

	filePath, err := pathPrompt.Run()
	if err != nil {
		return "", "", fmt.Errorf("export cancelled")
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		filePath = defaultFilename
	}

	// Add appropriate extension if not present
	ext := filepath.Ext(filePath)
	expectedExt := ""
	switch format {
	case FormatCSV:
		expectedExt = ".csv"
	case FormatJSON:
		expectedExt = ".json"
	case FormatText:
		expectedExt = ".txt"
	}

	if ext != expectedExt {
		filePath += expectedExt
	}

	return format, filePath, nil
}

// ExportTLE exports TLE data to the specified format and file.
func ExportTLE(tle TLE, format ExportFormat, filePath string) error {
	switch format {
	case FormatCSV:
		return exportTLECSV(tle, filePath)
	case FormatJSON:
		return exportTLEJSON(tle, filePath)
	case FormatText:
		return exportTLEText(tle, filePath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// exportTLECSV exports TLE data to CSV format.
func exportTLECSV(tle TLE, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	headers := []string{
		"Field", "Value",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	rows := [][]string{
		{"Common Name", tle.CommonName},
		{"Satellite Catalog Number", strconv.Itoa(tle.SatelliteCatalogNumber)},
		{"Elset Classification", tle.ElsetClassificiation},
		{"International Designator", tle.InternationalDesignator},
		{"Element Set Epoch (UTC)", fmt.Sprintf("%f", tle.ElementSetEpoch)},
		{"1st Derivative of Mean Motion", fmt.Sprintf("%f", tle.FirstDerivativeMeanMotion)},
		{"2nd Derivative of Mean Motion", tle.SecondDerivativeMeanMotion},
		{"B* Drag Term", tle.BDragTerm},
		{"Element Set Type", strconv.Itoa(tle.ElementSetType)},
		{"Element Number", strconv.Itoa(tle.ElementNumber)},
		{"Checksum Line One", strconv.Itoa(tle.ChecksumOne)},
		{"Orbit Inclination (degrees)", fmt.Sprintf("%f", tle.OrbitInclination)},
		{"Right Ascension (degrees)", fmt.Sprintf("%f", tle.RightAscension)},
		{"Eccentricity", fmt.Sprintf("%f", tle.Eccentrcity)},
		{"Argument of Perigee (degrees)", fmt.Sprintf("%f", tle.Perigee)},
		{"Mean Anomaly (degrees)", fmt.Sprintf("%f", tle.MeanAnamoly)},
		{"Mean Motion (revolutions/day)", fmt.Sprintf("%f", tle.MeanMotion)},
		{"Revolution Number at Epoch", strconv.Itoa(tle.RevolutionNumber)},
		{"Checksum Line Two", strconv.Itoa(tle.ChecksumTwo)},
	}

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// exportTLEJSON exports TLE data to JSON format.
func exportTLEJSON(tle TLE, filePath string) error {
	data := map[string]interface{}{
		"common_name":                      tle.CommonName,
		"satellite_catalog_number":         tle.SatelliteCatalogNumber,
		"elset_classification":             tle.ElsetClassificiation,
		"international_designator":          tle.InternationalDesignator,
		"element_set_epoch_utc":            tle.ElementSetEpoch,
		"first_derivative_mean_motion":      tle.FirstDerivativeMeanMotion,
		"second_derivative_mean_motion":     tle.SecondDerivativeMeanMotion,
		"b_drag_term":                       tle.BDragTerm,
		"element_set_type":                 tle.ElementSetType,
		"element_number":                    tle.ElementNumber,
		"checksum_line_one":                 tle.ChecksumOne,
		"orbit_inclination_degrees":         tle.OrbitInclination,
		"right_ascension_degrees":            tle.RightAscension,
		"eccentricity":                      tle.Eccentrcity,
		"argument_of_perigee_degrees":       tle.Perigee,
		"mean_anomaly_degrees":              tle.MeanAnamoly,
		"mean_motion_revolutions_per_day":   tle.MeanMotion,
		"revolution_number_at_epoch":        tle.RevolutionNumber,
		"checksum_line_two":                 tle.ChecksumTwo,
		"export_timestamp":                  time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// exportTLEText exports TLE data to text format.
func exportTLEText(tle TLE, filePath string) error {
	var builder strings.Builder

	builder.WriteString("Two-Line Element (TLE) Data\n")
	builder.WriteString(strings.Repeat("=", 60) + "\n\n")

	builder.WriteString(fmt.Sprintf("Common Name: %s\n", tle.CommonName))
	builder.WriteString(fmt.Sprintf("Satellite Catalog Number: %d\n", tle.SatelliteCatalogNumber))
	builder.WriteString(fmt.Sprintf("Elset Classification: %s\n", tle.ElsetClassificiation))
	builder.WriteString(fmt.Sprintf("International Designator: %s\n", tle.InternationalDesignator))
	builder.WriteString(fmt.Sprintf("Element Set Epoch (UTC): %f\n", tle.ElementSetEpoch))
	builder.WriteString(fmt.Sprintf("1st Derivative of Mean Motion: %f\n", tle.FirstDerivativeMeanMotion))
	builder.WriteString(fmt.Sprintf("2nd Derivative of Mean Motion: %s\n", tle.SecondDerivativeMeanMotion))
	builder.WriteString(fmt.Sprintf("B* Drag Term: %s\n", tle.BDragTerm))
	builder.WriteString(fmt.Sprintf("Element Set Type: %d\n", tle.ElementSetType))
	builder.WriteString(fmt.Sprintf("Element Number: %d\n", tle.ElementNumber))
	builder.WriteString(fmt.Sprintf("Checksum Line One: %d\n", tle.ChecksumOne))
	builder.WriteString(fmt.Sprintf("Orbit Inclination (degrees): %f\n", tle.OrbitInclination))
	builder.WriteString(fmt.Sprintf("Right Ascension (degrees): %f\n", tle.RightAscension))
	builder.WriteString(fmt.Sprintf("Eccentricity: %f\n", tle.Eccentrcity))
	builder.WriteString(fmt.Sprintf("Argument of Perigee (degrees): %f\n", tle.Perigee))
	builder.WriteString(fmt.Sprintf("Mean Anomaly (degrees): %f\n", tle.MeanAnamoly))
	builder.WriteString(fmt.Sprintf("Mean Motion (revolutions/day): %f\n", tle.MeanMotion))
	builder.WriteString(fmt.Sprintf("Revolution Number at Epoch: %d\n", tle.RevolutionNumber))
	builder.WriteString(fmt.Sprintf("Checksum Line Two: %d\n", tle.ChecksumTwo))
	builder.WriteString(fmt.Sprintf("\nExported: %s\n", time.Now().Format(time.RFC3339)))

	if err := os.WriteFile(filePath, []byte(builder.String()), 0644); err != nil {
		return fmt.Errorf("failed to write text file: %w", err)
	}

	return nil
}

// ExportVisualPrediction exports visual pass predictions to the specified format.
func ExportVisualPrediction(data VisualPassesResponse, format ExportFormat, filePath string) error {
	switch format {
	case FormatCSV:
		return exportVisualPredictionCSV(data, filePath)
	case FormatJSON:
		return exportVisualPredictionJSON(data, filePath)
	case FormatText:
		return exportVisualPredictionText(data, filePath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// exportVisualPredictionCSV exports visual pass predictions to CSV format.
func exportVisualPredictionCSV(data VisualPassesResponse, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write satellite info
	infoHeaders := []string{"Satellite Name", "Satellite ID", "Transactions Count", "Passes Count"}
	infoRow := []string{
		data.Info.SatName,
		strconv.Itoa(data.Info.SatID),
		strconv.Itoa(data.Info.TransactionsCount),
		strconv.Itoa(data.Info.PassesCount),
	}
	if err := writer.Write(infoHeaders); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}
	if err := writer.Write(infoRow); err != nil {
		return fmt.Errorf("failed to write CSV info row: %w", err)
	}

	// Write empty row separator
	if err := writer.Write([]string{}); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	// Write passes header
	passHeaders := []string{
		"Pass #", "Start Azimuth", "Start Azimuth Compass", "Start Elevation",
		"Start UTC", "Max Azimuth", "Max Azimuth Compass", "Max Elevation",
		"Max UTC", "End Azimuth", "End Azimuth Compass", "End Elevation",
		"End UTC", "Magnitude", "Duration (seconds)",
	}
	if err := writer.Write(passHeaders); err != nil {
		return fmt.Errorf("failed to write pass headers: %w", err)
	}

	// Write passes data
	for i, pass := range data.Passes {
		row := []string{
			strconv.Itoa(i + 1),
			fmt.Sprintf("%f", pass.StartAz),
			pass.StartAzCompass,
			fmt.Sprintf("%f", pass.StartEl),
			strconv.Itoa(pass.StartUTC),
			fmt.Sprintf("%f", pass.MaxAz),
			pass.MaxAzCompass,
			fmt.Sprintf("%f", pass.MaxEl),
			strconv.Itoa(pass.MaxUTC),
			fmt.Sprintf("%f", pass.EndAz),
			pass.EndAzCompass,
			fmt.Sprintf("%f", pass.EndEl),
			strconv.Itoa(pass.EndUTC),
			fmt.Sprintf("%f", pass.Mag),
			strconv.Itoa(pass.Duration),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write pass row: %w", err)
		}
	}

	return nil
}

// exportVisualPredictionJSON exports visual pass predictions to JSON format.
func exportVisualPredictionJSON(data VisualPassesResponse, filePath string) error {
	exportData := map[string]interface{}{
		"satellite_info": map[string]interface{}{
			"satellite_name":      data.Info.SatName,
			"satellite_id":        data.Info.SatID,
			"transactions_count":  data.Info.TransactionsCount,
			"passes_count":        data.Info.PassesCount,
		},
		"passes": data.Passes,
		"export_timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// exportVisualPredictionText exports visual pass predictions to text format.
func exportVisualPredictionText(data VisualPassesResponse, filePath string) error {
	var builder strings.Builder

	builder.WriteString("Visual Pass Predictions\n")
	builder.WriteString(strings.Repeat("=", 60) + "\n\n")

	builder.WriteString("Satellite Information:\n")
	builder.WriteString(fmt.Sprintf("  Name: %s\n", data.Info.SatName))
	builder.WriteString(fmt.Sprintf("  ID: %d\n", data.Info.SatID))
	builder.WriteString(fmt.Sprintf("  Transactions Count: %d\n", data.Info.TransactionsCount))
	builder.WriteString(fmt.Sprintf("  Passes Count: %d\n\n", data.Info.PassesCount))

	builder.WriteString("Passes:\n")
	builder.WriteString(strings.Repeat("-", 60) + "\n")

	for i, pass := range data.Passes {
		builder.WriteString(fmt.Sprintf("\nPass #%d:\n", i+1))
		builder.WriteString(fmt.Sprintf("  Start: Azimuth %.2f° (%s), Elevation %.2f°, UTC %d\n",
			pass.StartAz, pass.StartAzCompass, pass.StartEl, pass.StartUTC))
		builder.WriteString(fmt.Sprintf("  Max:   Azimuth %.2f° (%s), Elevation %.2f°, UTC %d\n",
			pass.MaxAz, pass.MaxAzCompass, pass.MaxEl, pass.MaxUTC))
		builder.WriteString(fmt.Sprintf("  End:   Azimuth %.2f° (%s), Elevation %.2f°, UTC %d\n",
			pass.EndAz, pass.EndAzCompass, pass.EndEl, pass.EndUTC))
		builder.WriteString(fmt.Sprintf("  Magnitude: %.2f, Duration: %d seconds\n",
			pass.Mag, pass.Duration))
	}

	builder.WriteString(fmt.Sprintf("\nExported: %s\n", time.Now().Format(time.RFC3339)))

	if err := os.WriteFile(filePath, []byte(builder.String()), 0644); err != nil {
		return fmt.Errorf("failed to write text file: %w", err)
	}

	return nil
}

// ExportRadioPrediction exports radio pass predictions to the specified format.
func ExportRadioPrediction(data RadioPassResponse, format ExportFormat, filePath string) error {
	switch format {
	case FormatCSV:
		return exportRadioPredictionCSV(data, filePath)
	case FormatJSON:
		return exportRadioPredictionJSON(data, filePath)
	case FormatText:
		return exportRadioPredictionText(data, filePath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// exportRadioPredictionCSV exports radio pass predictions to CSV format.
func exportRadioPredictionCSV(data RadioPassResponse, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write satellite info
	infoHeaders := []string{"Satellite Name", "Satellite ID", "Transactions Count", "Passes Count"}
	infoRow := []string{
		data.Info.SatName,
		strconv.Itoa(data.Info.SatID),
		strconv.Itoa(data.Info.TransactionsCount),
		strconv.Itoa(data.Info.PassesCount),
	}
	if err := writer.Write(infoHeaders); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}
	if err := writer.Write(infoRow); err != nil {
		return fmt.Errorf("failed to write CSV info row: %w", err)
	}

	// Write empty row separator
	if err := writer.Write([]string{}); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	// Write passes header
	passHeaders := []string{
		"Pass #", "Start Azimuth", "Start Azimuth Compass", "Start UTC",
		"Max Azimuth", "Max Azimuth Compass", "Max Elevation", "Max UTC",
		"End Azimuth", "End Azimuth Compass", "End UTC",
	}
	if err := writer.Write(passHeaders); err != nil {
		return fmt.Errorf("failed to write pass headers: %w", err)
	}

	// Write passes data
	for i, pass := range data.Passes {
		row := []string{
			strconv.Itoa(i + 1),
			fmt.Sprintf("%f", pass.StartAz),
			pass.StartAzCompass,
			strconv.FormatInt(pass.StartUTC, 10),
			fmt.Sprintf("%f", pass.MaxAz),
			pass.MaxAzCompass,
			fmt.Sprintf("%f", pass.MaxEl),
			strconv.FormatInt(pass.MaxUTC, 10),
			fmt.Sprintf("%f", pass.EndAz),
			pass.EndAzCompass,
			strconv.FormatInt(pass.EndUTC, 10),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write pass row: %w", err)
		}
	}

	return nil
}

// exportRadioPredictionJSON exports radio pass predictions to JSON format.
func exportRadioPredictionJSON(data RadioPassResponse, filePath string) error {
	exportData := map[string]interface{}{
		"satellite_info": map[string]interface{}{
			"satellite_name":      data.Info.SatName,
			"satellite_id":        data.Info.SatID,
			"transactions_count":  data.Info.TransactionsCount,
			"passes_count":        data.Info.PassesCount,
		},
		"passes": data.Passes,
		"export_timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// exportRadioPredictionText exports radio pass predictions to text format.
func exportRadioPredictionText(data RadioPassResponse, filePath string) error {
	var builder strings.Builder

	builder.WriteString("Radio Pass Predictions\n")
	builder.WriteString(strings.Repeat("=", 60) + "\n\n")

	builder.WriteString("Satellite Information:\n")
	builder.WriteString(fmt.Sprintf("  Name: %s\n", data.Info.SatName))
	builder.WriteString(fmt.Sprintf("  ID: %d\n", data.Info.SatID))
	builder.WriteString(fmt.Sprintf("  Transactions Count: %d\n", data.Info.TransactionsCount))
	builder.WriteString(fmt.Sprintf("  Passes Count: %d\n\n", data.Info.PassesCount))

	builder.WriteString("Passes:\n")
	builder.WriteString(strings.Repeat("-", 60) + "\n")

	for i, pass := range data.Passes {
		builder.WriteString(fmt.Sprintf("\nPass #%d:\n", i+1))
		builder.WriteString(fmt.Sprintf("  Start: Azimuth %.2f° (%s), UTC %d\n",
			pass.StartAz, pass.StartAzCompass, pass.StartUTC))
		builder.WriteString(fmt.Sprintf("  Max:   Azimuth %.2f° (%s), Elevation %.2f°, UTC %d\n",
			pass.MaxAz, pass.MaxAzCompass, pass.MaxEl, pass.MaxUTC))
		builder.WriteString(fmt.Sprintf("  End:   Azimuth %.2f° (%s), UTC %d\n",
			pass.EndAz, pass.EndAzCompass, pass.EndUTC))
	}

	builder.WriteString(fmt.Sprintf("\nExported: %s\n", time.Now().Format(time.RFC3339)))

	if err := os.WriteFile(filePath, []byte(builder.String()), 0644); err != nil {
		return fmt.Errorf("failed to write text file: %w", err)
	}

	return nil
}

// ExportSatellitePosition exports satellite position data to the specified format.
func ExportSatellitePosition(data Response, format ExportFormat, filePath string) error {
	switch format {
	case FormatCSV:
		return exportSatellitePositionCSV(data, filePath)
	case FormatJSON:
		return exportSatellitePositionJSON(data, filePath)
	case FormatText:
		return exportSatellitePositionText(data, filePath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// exportSatellitePositionCSV exports satellite positions to CSV format.
func exportSatellitePositionCSV(data Response, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write satellite info
	infoHeaders := []string{"Satellite Name", "Satellite ID"}
	infoRow := []string{
		data.SatelliteInfo.Satname,
		strconv.Itoa(data.SatelliteInfo.Satid),
	}
	if err := writer.Write(infoHeaders); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}
	if err := writer.Write(infoRow); err != nil {
		return fmt.Errorf("failed to write CSV info row: %w", err)
	}

	// Write empty row separator
	if err := writer.Write([]string{}); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	// Write positions header
	posHeaders := []string{
		"Position #", "Latitude", "Longitude", "Altitude (km)",
		"Azimuth", "Declination", "Timestamp",
	}
	if err := writer.Write(posHeaders); err != nil {
		return fmt.Errorf("failed to write position headers: %w", err)
	}

	// Write positions data
	for i, pos := range data.Positions {
		row := []string{
			strconv.Itoa(i + 1),
			fmt.Sprintf("%f", pos.Satlatitude),
			fmt.Sprintf("%f", pos.Satlongitude),
			fmt.Sprintf("%f", pos.Sataltitude),
			fmt.Sprintf("%f", pos.Azimuth),
			fmt.Sprintf("%f", pos.Dec),
			strconv.FormatInt(pos.Timestamp, 10),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write position row: %w", err)
		}
	}

	return nil
}

// exportSatellitePositionJSON exports satellite positions to JSON format.
func exportSatellitePositionJSON(data Response, filePath string) error {
	exportData := map[string]interface{}{
		"satellite_info": map[string]interface{}{
			"satellite_name": data.SatelliteInfo.Satname,
			"satellite_id":   data.SatelliteInfo.Satid,
		},
		"positions": data.Positions,
		"export_timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// exportSatellitePositionText exports satellite positions to text format.
func exportSatellitePositionText(data Response, filePath string) error {
	var builder strings.Builder

	builder.WriteString("Satellite Position Data\n")
	builder.WriteString(strings.Repeat("=", 60) + "\n\n")

	builder.WriteString("Satellite Information:\n")
	builder.WriteString(fmt.Sprintf("  Name: %s\n", data.SatelliteInfo.Satname))
	builder.WriteString(fmt.Sprintf("  ID: %d\n\n", data.SatelliteInfo.Satid))

	builder.WriteString("Positions:\n")
	builder.WriteString(strings.Repeat("-", 60) + "\n")

	for i, pos := range data.Positions {
		builder.WriteString(fmt.Sprintf("\nPosition #%d:\n", i+1))
		builder.WriteString(fmt.Sprintf("  Latitude:  %.6f°\n", pos.Satlatitude))
		builder.WriteString(fmt.Sprintf("  Longitude: %.6f°\n", pos.Satlongitude))
		builder.WriteString(fmt.Sprintf("  Altitude:  %.2f km\n", pos.Sataltitude))
		builder.WriteString(fmt.Sprintf("  Azimuth:   %.2f°\n", pos.Azimuth))
		builder.WriteString(fmt.Sprintf("  Declination: %.2f°\n", pos.Dec))
		builder.WriteString(fmt.Sprintf("  Timestamp: %d\n", pos.Timestamp))
	}

	builder.WriteString(fmt.Sprintf("\nExported: %s\n", time.Now().Format(time.RFC3339)))

	if err := os.WriteFile(filePath, []byte(builder.String()), 0644); err != nil {
		return fmt.Errorf("failed to write text file: %w", err)
	}

	return nil
}

