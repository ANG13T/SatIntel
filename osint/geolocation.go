package osint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TwiN/go-color"
)

// LocationData represents the user's geographic location
type LocationData struct {
	Latitude  float64
	Longitude float64
	City      string
	Country   string
	Region    string
	Timezone  string
}

// GetUserLocation automatically detects the user's location using IP geolocation.
// Returns latitude, longitude, and location info, or an error if detection fails.
func GetUserLocation() (*LocationData, error) {
	fmt.Println(color.Ize(color.Cyan, "  [*] Detecting your location..."))

	// Try multiple free geolocation APIs for reliability
	apis := []struct {
		name string
		url  string
		parse func([]byte) (*LocationData, error)
	}{
		{
			name: "ip-api.com",
			url:  "http://ip-api.com/json/?fields=status,lat,lon,city,country,regionName,timezone",
			parse: parseIPAPIResponse,
		},
		{
			name: "ipapi.co",
			url:  "https://ipapi.co/json/",
			parse: parseIPAPICoResponse,
		},
		{
			name: "ip-api.com (backup)",
			url:  "https://ip-api.com/json/?fields=status,lat,lon,city,country,regionName,timezone",
			parse: parseIPAPIResponse,
		},
	}

	var lastErr error
	for _, api := range apis {
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get(api.url)
		if err != nil {
			lastErr = fmt.Errorf("failed to connect to %s: %w", api.name, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("%s returned status %d", api.name, resp.StatusCode)
			continue
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response from %s: %w", api.name, err)
			continue
		}

		location, err := api.parse(body)
		if err != nil {
			lastErr = fmt.Errorf("failed to parse %s response: %w", api.name, err)
			continue
		}

		if location != nil && location.Latitude != 0 && location.Longitude != 0 {
			fmt.Println(color.Ize(color.Green, fmt.Sprintf("  [+] Location detected: %s, %s", location.City, location.Country)))
			return location, nil
		}
	}

	return nil, fmt.Errorf("failed to detect location from all APIs: %w", lastErr)
}

// parseIPAPIResponse parses the response from ip-api.com
func parseIPAPIResponse(body []byte) (*LocationData, error) {
	var response struct {
		Status    string  `json:"status"`
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
		City      string  `json:"city"`
		Country   string  `json:"country"`
		Region    string  `json:"regionName"`
		Timezone  string  `json:"timezone"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, fmt.Errorf("API returned status: %s", response.Status)
	}

	if response.Latitude == 0 && response.Longitude == 0 {
		return nil, fmt.Errorf("invalid coordinates")
	}

	return &LocationData{
		Latitude:  response.Latitude,
		Longitude: response.Longitude,
		City:      response.City,
		Country:   response.Country,
		Region:    response.Region,
		Timezone:  response.Timezone,
	}, nil
}

// parseIPAPICoResponse parses the response from ipapi.co
func parseIPAPICoResponse(body []byte) (*LocationData, error) {
	var response struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		City      string  `json:"city"`
		Country   string  `json:"country_name"`
		Region    string  `json:"region"`
		Timezone  string  `json:"timezone"`
		Error     bool    `json:"error"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Error {
		return nil, fmt.Errorf("API returned an error")
	}

	if response.Latitude == 0 && response.Longitude == 0 {
		return nil, fmt.Errorf("invalid coordinates")
	}

	return &LocationData{
		Latitude:  response.Latitude,
		Longitude: response.Longitude,
		City:      response.City,
		Country:   response.Country,
		Region:    response.Region,
		Timezone:  response.Timezone,
	}, nil
}

// GetLocationWithPrompt attempts to get location automatically, with manual fallback option.
// Returns latitude, longitude as strings, and a boolean indicating if location was auto-detected.
func GetLocationWithPrompt() (string, string, bool) {
	// Try to auto-detect location
	location, err := GetUserLocation()
	if err != nil {
		fmt.Println(color.Ize(color.Yellow, fmt.Sprintf("  [!] Auto-detection failed: %s", err.Error())))
		fmt.Println(color.Ize(color.Cyan, "  [*] Please enter your location manually:"))
		return getManualLocation()
	}

	// Show detected location and ask for confirmation
	fmt.Println(color.Ize(color.Cyan, fmt.Sprintf("\n  Detected Location:")))
	fmt.Println(color.Ize(color.White, fmt.Sprintf("    City: %s", location.City)))
	fmt.Println(color.Ize(color.White, fmt.Sprintf("    Region: %s", location.Region)))
	fmt.Println(color.Ize(color.White, fmt.Sprintf("    Country: %s", location.Country)))
	fmt.Println(color.Ize(color.White, fmt.Sprintf("    Coordinates: %.6f, %.6f", location.Latitude, location.Longitude)))
	
	fmt.Print(color.Ize(color.Cyan, "\n  Use this location? (y/n, default: y) > "))
	var confirm string
	fmt.Scanln(&confirm)
	confirm = strings.ToLower(strings.TrimSpace(confirm))

	if confirm == "" || confirm == "y" || confirm == "yes" {
		return fmt.Sprintf("%.6f", location.Latitude), fmt.Sprintf("%.6f", location.Longitude), true
	}

	// User wants to enter manually
	fmt.Println(color.Ize(color.Cyan, "  [*] Please enter your location manually:"))
	return getManualLocation()
}

// getManualLocation prompts the user to manually enter their location.
func getManualLocation() (string, string, bool) {
	fmt.Print("\n ENTER LATITUDE > ")
	var latitude string
	fmt.Scanln(&latitude)
	if strings.TrimSpace(latitude) == "" {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Latitude cannot be empty"))
		return "", "", false
	}

	// Validate latitude
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Invalid latitude format"))
		return "", "", false
	}
	if lat < -90 || lat > 90 {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Latitude must be between -90 and 90"))
		return "", "", false
	}

	fmt.Print("\n ENTER LONGITUDE > ")
	var longitude string
	fmt.Scanln(&longitude)
	if strings.TrimSpace(longitude) == "" {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Longitude cannot be empty"))
		return "", "", false
	}

	// Validate longitude
	lon, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Invalid longitude format"))
		return "", "", false
	}
	if lon < -180 || lon > 180 {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Longitude must be between -180 and 180"))
		return "", "", false
	}

	return latitude, longitude, false
}

