package osint

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/TwiN/go-color"
	satellite "github.com/joshuaferrara/go-satellite"
)

// SGPPosition represents a satellite position calculated using SGP4.
type SGPPosition struct {
	Latitude  float64 // Satellite latitude in degrees
	Longitude float64 // Satellite longitude in degrees
	Altitude  float64 // Satellite altitude in kilometers
	Velocity  float64 // Satellite velocity in km/s
	Timestamp int64   // Unix timestamp
}

// ObserverPosition represents the position of an observer on Earth.
type ObserverPosition struct {
	Latitude  float64 // Observer latitude in degrees
	Longitude float64 // Observer longitude in degrees
	Altitude  float64 // Observer altitude in meters
}

// LookAngles represents the viewing angles from an observer to a satellite.
type LookAngles struct {
	Azimuth   float64 // Azimuth angle in degrees (0-360)
	Elevation float64 // Elevation angle in degrees (-90 to 90)
	Range     float64 // Range to satellite in kilometers
	RangeRate float64 // Range rate in km/s
}

// SGP4PositionResult contains the calculated position and look angles.
type SGP4PositionResult struct {
	Position   SGPPosition
	LookAngles LookAngles
}

// CalculateSGP4Position calculates the satellite position using SGP4 algorithm from raw TLE line strings.
// This is the recommended function to use as it works directly with TLE line strings.
func CalculateSGP4Position(line1, line2 string, targetTime time.Time) (SGPPosition, error) {
	// Validate TLE lines
	line1 = strings.TrimSpace(line1)
	line2 = strings.TrimSpace(line2)

	if len(line1) < 69 || len(line2) < 69 {
		return SGPPosition{}, fmt.Errorf("invalid TLE: lines must be at least 69 characters")
	}

	if !strings.HasPrefix(line1, "1 ") {
		return SGPPosition{}, fmt.Errorf("invalid TLE: line 1 must start with '1 '")
	}
	if !strings.HasPrefix(line2, "2 ") {
		return SGPPosition{}, fmt.Errorf("invalid TLE: line 2 must start with '2 '")
	}

	// Parse TLE using the library (using WGS72 as default gravity model)
	sat := satellite.TLEToSat(line1, line2, satellite.GravityWGS72)

	// Propagate satellite to target time
	year := targetTime.Year()
	month := int(targetTime.Month())
	day := targetTime.Day()
	hour := targetTime.Hour()
	minute := targetTime.Minute()
	second := targetTime.Second()

	position, velocity := satellite.Propagate(sat, year, month, day, hour, minute, second)

	// Calculate Julian Day for the target time
	jday := satellite.JDay(year, month, day, hour, minute, second)

	// Calculate Greenwich Mean Sidereal Time
	gmst := satellite.ThetaG_JD(jday)

	// Convert ECI to Lat/Long/Alt
	altitude, _, latLong := satellite.ECIToLLA(position, gmst)

	// Calculate velocity magnitude
	velocityMagnitude := math.Sqrt(velocity.X*velocity.X + velocity.Y*velocity.Y + velocity.Z*velocity.Z)

	return SGPPosition{
		Latitude:  latLong.Latitude * satellite.RAD2DEG,
		Longitude: latLong.Longitude * satellite.RAD2DEG,
		Altitude:  altitude / 1000.0, // Convert meters to kilometers
		Velocity:  velocityMagnitude,
		Timestamp: targetTime.Unix(),
	}, nil
}

// CalculateSGP4PositionFromTLE calculates position from a TLE struct.
// Note: This requires the original TLE lines. If you have raw TLE lines, use CalculateSGP4Position instead.
func CalculateSGP4PositionFromTLE(tle TLE, line1, line2 string, targetTime time.Time) (SGPPosition, error) {
	return CalculateSGP4Position(line1, line2, targetTime)
}

// CalculateSGP4PositionWithObserver calculates satellite position and look angles from an observer's perspective.
// This is the recommended function to use as it works directly with TLE line strings.
func CalculateSGP4PositionWithObserver(line1, line2 string, targetTime time.Time, observer ObserverPosition) (SGP4PositionResult, error) {
	// First calculate the satellite position
	satPosition, err := CalculateSGP4Position(line1, line2, targetTime)
	if err != nil {
		return SGP4PositionResult{}, err
	}

	// Validate and trim TLE lines
	line1 = strings.TrimSpace(line1)
	line2 = strings.TrimSpace(line2)

	// Parse TLE
	sat := satellite.TLEToSat(line1, line2, satellite.GravityWGS72)

	// Propagate satellite to target time
	year := targetTime.Year()
	month := int(targetTime.Month())
	day := targetTime.Day()
	hour := targetTime.Hour()
	minute := targetTime.Minute()
	second := targetTime.Second()

	position, _ := satellite.Propagate(sat, year, month, day, hour, minute, second)

	// Calculate Julian Day
	jday := satellite.JDay(year, month, day, hour, minute, second)

	// Convert observer position to ECI coordinates
	obsLatLong := satellite.LatLong{
		Latitude:  observer.Latitude * satellite.DEG2RAD,
		Longitude: observer.Longitude * satellite.DEG2RAD,
	}
	obsAlt := observer.Altitude / 1000.0 // Convert meters to kilometers
	obsECI := satellite.LLAToECI(obsLatLong, obsAlt, jday)

	// Calculate look angles
	lookAngles := satellite.ECIToLookAngles(position, obsLatLong, obsAlt, jday)

	// Calculate range (distance from observer to satellite)
	dx := position.X - obsECI.X
	dy := position.Y - obsECI.Y
	dz := position.Z - obsECI.Z
	rangeKm := math.Sqrt(dx*dx+dy*dy+dz*dz) / 1000.0 // Convert meters to kilometers

	return SGP4PositionResult{
		Position: satPosition,
		LookAngles: LookAngles{
			Azimuth:   lookAngles.Az * satellite.RAD2DEG,
			Elevation: lookAngles.El * satellite.RAD2DEG,
			Range:     rangeKm,
			RangeRate: 0.0, // Range rate calculation would require velocity comparison
		},
	}, nil
}

// CalculateSGP4Positions calculates multiple positions over a time range.
func CalculateSGP4Positions(line1, line2 string, startTime time.Time, endTime time.Time, interval time.Duration) ([]SGPPosition, error) {
	if startTime.After(endTime) {
		return nil, fmt.Errorf("start time must be before end time")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("interval must be positive")
	}

	var positions []SGPPosition
	currentTime := startTime

	for !currentTime.After(endTime) {
		pos, err := CalculateSGP4Position(line1, line2, currentTime)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate position at %v: %w", currentTime, err)
		}
		positions = append(positions, pos)
		currentTime = currentTime.Add(interval)
	}

	return positions, nil
}

// CalculateSGP4PositionFromTLEStruct calculates position from a TLE struct with original lines.
// This is a convenience wrapper that uses the provided TLE lines.
func CalculateSGP4PositionFromTLEStruct(tle TLE, originalLine1, originalLine2 string, targetTime time.Time) (SGPPosition, error) {
	return CalculateSGP4Position(originalLine1, originalLine2, targetTime)
}

// CalculateSGP4PositionFromTLEStructWithObserver calculates position and look angles from a TLE struct with original lines.
func CalculateSGP4PositionFromTLEStructWithObserver(tle TLE, originalLine1, originalLine2 string, targetTime time.Time, observer ObserverPosition) (SGP4PositionResult, error) {
	return CalculateSGP4PositionWithObserver(originalLine1, originalLine2, targetTime, observer)
}

// PrintSGP4Position displays SGP4-calculated position in a formatted table.
func PrintSGP4Position(pos SGPPosition) {
	fmt.Println(color.Ize(color.Purple, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Purple, "║              SGP4 Calculated Position                       ║"))
	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
	fmt.Println(color.Ize(color.Purple, GenRowString("Latitude (degrees)", fmt.Sprintf("%.6f", pos.Latitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Longitude (degrees)", fmt.Sprintf("%.6f", pos.Longitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Altitude (km)", fmt.Sprintf("%.2f", pos.Altitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Velocity (km/s)", fmt.Sprintf("%.4f", pos.Velocity))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Timestamp", fmt.Sprintf("%d", pos.Timestamp))))
	fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝\n\n"))
}

// PrintSGP4PositionWithLookAngles displays position and look angles in a formatted table.
func PrintSGP4PositionWithLookAngles(result SGP4PositionResult) {
	fmt.Println(color.Ize(color.Purple, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Purple, "║         SGP4 Calculated Position & Look Angles             ║"))
	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Latitude (degrees)", fmt.Sprintf("%.6f", result.Position.Latitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Longitude (degrees)", fmt.Sprintf("%.6f", result.Position.Longitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Altitude (km)", fmt.Sprintf("%.2f", result.Position.Altitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Velocity (km/s)", fmt.Sprintf("%.4f", result.Position.Velocity))))
	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
	fmt.Println(color.Ize(color.Purple, GenRowString("Azimuth (degrees)", fmt.Sprintf("%.2f", result.LookAngles.Azimuth))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Elevation (degrees)", fmt.Sprintf("%.2f", result.LookAngles.Elevation))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Range (km)", fmt.Sprintf("%.2f", result.LookAngles.Range))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Range Rate (km/s)", fmt.Sprintf("%.4f", result.LookAngles.RangeRate))))
	fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝\n\n"))
}
