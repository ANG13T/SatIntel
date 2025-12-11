package osint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/manifoldco/promptui"
)

const (
	authURL      = "https://www.space-track.org/ajaxauth/login"
	queryBaseURL = "https://www.space-track.org/basicspacedata/query"
)

// Login authenticates with Space-Track API using credentials from environment variables.
// Returns an HTTP client with a cookie jar to maintain the session.
func Login() (*http.Client, error) {
	spinner := ShowLoginProgress()
	defer spinner.Stop()

	vals := url.Values{}
	vals.Add("identity", os.Getenv("SPACE_TRACK_USERNAME"))
	vals.Add("password", os.Getenv("SPACE_TRACK_PASSWORD"))

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, NewAppErrorWithErr(ErrCodeAuthCookieJar, "Failed to create cookie jar for authentication", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	resp, err := client.PostForm(authURL, vals)
	if err != nil {
		return nil, NewAppErrorWithErr(ErrCodeAuthConnection, "Unable to connect to Space-Track API", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		username := os.Getenv("SPACE_TRACK_USERNAME")
		context := fmt.Sprintf("Status code: %d, Username: %s", resp.StatusCode, username)
		return nil, NewAppErrorWithContext(ErrCodeAuthFailed, "Authentication failed with Space-Track API", context)
	}

	spinner.Stop()
	fmt.Println(color.Ize(color.Green, "  [+] Logged in successfully"))
	return client, nil
}

// QuerySpaceTrack sends a GET request to the Space-Track API using the authenticated client.
// Returns the response body as a string.
func QuerySpaceTrack(client *http.Client, endpoint string) (string, error) {
	spinner := ShowQueryProgress(endpoint)
	defer spinner.Stop()

	req, err := http.NewRequest("GET", queryBaseURL+endpoint, nil)
	if err != nil {
		spinner.Stop()
		return "", fmt.Errorf("failed to create query request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		spinner.Stop()
		return "", fmt.Errorf("failed to fetch data from Space-Track: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		spinner.Stop()
		return "", fmt.Errorf("query returned non-success status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		spinner.Stop()
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	spinner.Stop()
	return string(body), nil
}

// extractNorad extracts the NORAD ID from a string in the format "Name (NORAD_ID)".
func extractNorad(str string) string {
	start := strings.Index(str, "(")
	end := strings.Index(str, ")")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return str[start+1 : end]
}

// PrintNORADInfo fetches and displays TLE data for a satellite identified by its NORAD ID.
func PrintNORADInfo(norad string, name string) {
	client, err := Login()
	if err != nil {
		HandleError(err, ErrCodeAuthFailed, "Failed to authenticate with Space-Track")
		return
	}

	endpoint := fmt.Sprintf("/class/gp_history/format/tle/NORAD_CAT_ID/%s/orderby/EPOCH%%20desc/limit/1", norad)
	data, err := QuerySpaceTrack(client, endpoint)
	if err != nil {
		context := fmt.Sprintf("NORAD ID: %s, Satellite: %s", norad, name)
		HandleErrorWithContext(err, ErrCodeAPINoData, "Failed to fetch TLE data for satellite", context)
		return
	}

	lines := strings.Split(strings.TrimSpace(data), "\n")

	var lineOne, lineTwo string

	if len(lines) >= 2 {
		lineOne = strings.TrimSpace(lines[0])
		lineTwo = strings.TrimSpace(lines[1])
	} else {
		tleLines := strings.Fields(data)
		if len(tleLines) < 2 {
			err := NewAppErrorWithContext(
				ErrCodeTLEInsufficientData,
				"Invalid TLE data - insufficient fields",
				fmt.Sprintf("NORAD ID: %s, Fields found: %d", norad, len(tleLines)),
			)
			err.Display()
			return
		}

		// Calculate split point, ensuring we can create two meaningful lines
		mid := 9
		if len(tleLines) < 9 {
			mid = len(tleLines) / 2
		}
		if mid > len(tleLines) {
			mid = len(tleLines) / 2
		}
		// Ensure mid is at least 1 to prevent empty lineOne
		if mid < 1 {
			mid = 1
		}
		// Ensure mid is less than len(tleLines) to ensure lineTwo is non-empty
		if mid >= len(tleLines) {
			mid = len(tleLines) - 1
		}

		lineOne = strings.Join(tleLines[:mid], " ")
		lineTwo = strings.Join(tleLines[mid:], " ")
	}

	if !strings.HasPrefix(lineOne, "1 ") {
		prefixLen := 10
		if len(lineOne) < prefixLen {
			prefixLen = len(lineOne)
		}
		err := NewAppErrorWithContext(
			ErrCodeTLEInvalidFormat,
			"Invalid TLE format - line 1 should start with '1 '",
			fmt.Sprintf("NORAD ID: %s, Line 1 prefix: %s", norad, lineOne[:prefixLen]),
		)
		err.Display()
		return
	}
	if !strings.HasPrefix(lineTwo, "2 ") {
		prefixLen := 10
		if len(lineTwo) < prefixLen {
			prefixLen = len(lineTwo)
		}
		err := NewAppErrorWithContext(
			ErrCodeTLEInvalidFormat,
			"Invalid TLE format - line 2 should start with '2 '",
			fmt.Sprintf("NORAD ID: %s, Line 2 prefix: %s", norad, lineTwo[:prefixLen]),
		)
		err.Display()
		return
	}

	tle := ConstructTLE(name, lineOne, lineTwo)

	parsingFailed := false

	line1Fields := strings.Fields(lineOne)
	line2Fields := strings.Fields(lineTwo)

	if len(line1Fields) < 4 || len(line2Fields) < 3 {
		parsingFailed = true
	} else if tle.SatelliteCatalogNumber == 0 && tle.InternationalDesignator == "" && tle.ElementSetEpoch == 0.0 {
		parsingFailed = true
	}

	if parsingFailed {
		context := fmt.Sprintf("NORAD ID: %s, Line 1 fields: %d, Line 2 fields: %d", norad, len(line1Fields), len(line2Fields))
		err := NewAppErrorWithContext(
			ErrCodeTLEParseFailed,
			"Failed to parse TLE data",
			context,
		)
		if len(line1Fields) >= 4 && len(line2Fields) >= 3 {
			err.Suggestions = append(err.Suggestions, "Field count is sufficient, but parsing failed. Check TLE format and data integrity.")
		}
		err.Display()
		return
	}

	PrintTLE(tle)
}

// buildSatcatQuery constructs a Space-Track API query string with optional filters and pagination.
// Note: Space-Track API uses path segments for filtering. For name search, we'll filter client-side.
func buildSatcatQuery(searchName, country, objectType, launchYear string, page, pageSize int) string {
	var parts []string
	parts = append(parts, "/class/satcat")

	// Add filters (name search is handled client-side for partial matching)
	if country != "" {
		parts = append(parts, fmt.Sprintf("/COUNTRY/%s", url.QueryEscape(country)))
	}
	if objectType != "" {
		parts = append(parts, fmt.Sprintf("/OBJECT_TYPE/%s", url.QueryEscape(objectType)))
	}
	if launchYear != "" {
		parts = append(parts, fmt.Sprintf("/LAUNCH_YEAR/%s", url.QueryEscape(launchYear)))
	}

	// Add ordering
	parts = append(parts, "/orderby/SATNAME%20asc")

	// For name search, fetch more results to filter client-side
	// Otherwise use normal pagination
	if searchName != "" {
		// Fetch more results for client-side filtering
		parts = append(parts, "/limit/500")
	} else if pageSize > 0 {
		offset := (page - 1) * pageSize
		parts = append(parts, fmt.Sprintf("/limit/%d,%d", pageSize, offset))
	} else {
		parts = append(parts, "/limit/50")
	}

	parts = append(parts, "/emptyresult/show")
	return strings.Join(parts, "")
}

// filterSatellitesByName filters satellites by name (case-insensitive partial match).
func filterSatellitesByName(sats []Satellite, searchName string) []Satellite {
	if searchName == "" {
		return sats
	}
	searchLower := strings.ToLower(searchName)
	var filtered []Satellite
	for _, sat := range sats {
		if strings.Contains(strings.ToLower(sat.SATNAME), searchLower) {
			filtered = append(filtered, sat)
		}
	}
	return filtered
}

// showSearchMenu displays an interactive menu for searching satellites.
func showSearchMenu() (string, string, string, string) {
	searchName := ""
	country := ""
	objectType := ""
	launchYear := ""

	for {
		menuItems := []string{
			"Search by Name",
			"Filter by Country",
			"Filter by Object Type",
			"Filter by Launch Year",
			"Clear All Filters",
			"Search & Continue",
		}

		prompt := promptui.Select{
			Label: "Satellite Search & Filter Options",
			Items: menuItems,
			Size:  10,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return "", "", "", ""
		}

		switch idx {
		case 0: // Search by Name
			namePrompt := promptui.Prompt{
				Label:     "Enter satellite name (or part of name)",
				Default:   searchName,
				AllowEdit: true,
			}
			result, err := namePrompt.Run()
			if err == nil {
				searchName = strings.TrimSpace(result)
			}

		case 1: // Filter by Country
			countryPrompt := promptui.Prompt{
				Label:     "Enter country code (e.g., US, RU, CN)",
				Default:   country,
				AllowEdit: true,
			}
			result, err := countryPrompt.Run()
			if err == nil {
				country = strings.TrimSpace(result)
			}

		case 2: // Filter by Object Type
			typeItems := []string{
				"PAYLOAD",
				"ROCKET BODY",
				"DEBRIS",
				"UNKNOWN",
				"TBA",
				"",
			}
			typePrompt := promptui.Select{
				Label: "Select Object Type",
				Items: typeItems,
			}
			_, result, err := typePrompt.Run()
			if err == nil {
				objectType = result
			}

		case 3: // Filter by Launch Year
			yearPrompt := promptui.Prompt{
				Label:     "Enter launch year (e.g., 2020)",
				Default:   launchYear,
				AllowEdit: true,
			}
			result, err := yearPrompt.Run()
			if err == nil {
				launchYear = strings.TrimSpace(result)
			}

		case 4: // Clear All Filters
			searchName = ""
			country = ""
			objectType = ""
			launchYear = ""
			fmt.Println(color.Ize(color.Green, "  [+] All filters cleared"))

		case 5: // Search & Continue
			return searchName, country, objectType, launchYear
		}

		// Show current filters
		if searchName != "" || country != "" || objectType != "" || launchYear != "" {
			fmt.Println(color.Ize(color.Cyan, "\n  Current Filters:"))
			if searchName != "" {
				fmt.Printf("    Name: %s\n", searchName)
			}
			if country != "" {
				fmt.Printf("    Country: %s\n", country)
			}
			if objectType != "" {
				fmt.Printf("    Object Type: %s\n", objectType)
			}
			if launchYear != "" {
				fmt.Printf("    Launch Year: %s\n", launchYear)
			}
			fmt.Println()
		}
	}
}

// SelectSatellite fetches a list of satellites from Space-Track with search, filter, and pagination support.
// Returns the selected satellite name with its NORAD ID in parentheses.
func SelectSatellite() string {
	// First, show option to select from favorites or search
	initialMenu := []string{
		"⭐ Select from Favorites",
		"🔍 Search Satellites",
		"❌ Cancel",
	}

	initialPrompt := promptui.Select{
		Label: "Satellite Selection",
		Items: initialMenu,
	}

	initialIdx, _, err := initialPrompt.Run()
	if err != nil {
		return ""
	}

	if initialIdx == 0 {
		// Select from favorites
		result := SelectFromFavorites()
		if result != "" {
			// Offer to remove from favorites or continue
			return result
		}
		return ""
	} else if initialIdx == 2 {
		// Cancel
		return ""
	}

	// Continue with search
	client, err := Login()
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: "+err.Error()))
		return ""
	}

	// Show search/filter menu
	searchName, country, objectType, launchYear := showSearchMenu()

	page := 1
	pageSize := 20
	var allFilteredSats []Satellite
	var totalPages int

	for {
		var sats []Satellite

		// If we have name search, we need to fetch all and filter client-side
		// Cache the filtered results to avoid refetching
		if searchName != "" && len(allFilteredSats) == 0 {
			// Fetch a larger batch for client-side filtering
			spinner := ShowProgressWithSpinner("Searching satellite catalog")
			endpoint := buildSatcatQuery(searchName, country, objectType, launchYear, 1, 0)
			data, err := QuerySpaceTrack(client, endpoint)
			spinner.Stop()
			if err != nil {
				context := fmt.Sprintf("Search: %s, Country: %s, Object Type: %s, Launch Year: %s", searchName, country, objectType, launchYear)
				HandleErrorWithContext(err, ErrCodeAPINoData, "Failed to fetch satellite catalog", context)
				return ""
			}

			var fetchedSats []Satellite
			if err := json.Unmarshal([]byte(data), &fetchedSats); err != nil {
				context := fmt.Sprintf("Search: %s, Response length: %d bytes", searchName, len(data))
				HandleErrorWithContext(err, ErrCodeAPIParseFailed, "Failed to parse satellite catalog data", context)
				return ""
			}

			// Apply client-side name filtering
			allFilteredSats = filterSatellitesByName(fetchedSats, searchName)
			totalPages = (len(allFilteredSats) + pageSize - 1) / pageSize

			// Apply pagination
			startIdx := (page - 1) * pageSize
			endIdx := startIdx + pageSize
			if endIdx > len(allFilteredSats) {
				endIdx = len(allFilteredSats)
			}
			if startIdx < len(allFilteredSats) {
				sats = allFilteredSats[startIdx:endIdx]
			} else {
				sats = []Satellite{}
			}
		} else if searchName != "" {
			// Use cached filtered results with pagination
			startIdx := (page - 1) * pageSize
			endIdx := startIdx + pageSize
			if endIdx > len(allFilteredSats) {
				endIdx = len(allFilteredSats)
			}
			if startIdx < len(allFilteredSats) {
				sats = allFilteredSats[startIdx:endIdx]
			} else {
				sats = []Satellite{}
			}
		} else {
			// No name search - use server-side pagination
			spinner := ShowProgressWithSpinner("Loading satellite catalog")
			endpoint := buildSatcatQuery(searchName, country, objectType, launchYear, page, pageSize)
			data, err := QuerySpaceTrack(client, endpoint)
			spinner.Stop()
			if err != nil {
				context := fmt.Sprintf("Page: %d, Country: %s, Object Type: %s", page, country, objectType)
				HandleErrorWithContext(err, ErrCodeAPINoData, "Failed to fetch satellite catalog", context)
				return ""
			}

			if err := json.Unmarshal([]byte(data), &sats); err != nil {
				context := fmt.Sprintf("Page: %d, Response length: %d bytes", page, len(data))
				HandleErrorWithContext(err, ErrCodeAPIParseFailed, "Failed to parse satellite catalog data", context)
				return ""
			}
			totalPages = 0 // Unknown for server-side pagination
		}

		if len(sats) == 0 {
			err := NewAppErrorWithContext(
				ErrCodeSatNoResults,
				"No satellites found with current filters",
				fmt.Sprintf("Search: %s, Country: %s, Object Type: %s, Launch Year: %s", searchName, country, objectType, launchYear),
			)
			err.Display()
			return ""
		}

		// Build display strings with additional info
		var satStrings []string
		for _, sat := range sats {
			info := fmt.Sprintf("%s (%s)", sat.SATNAME, sat.NORAD_CAT_ID)
			if sat.COUNTRY != "" {
				info += fmt.Sprintf(" - %s", sat.COUNTRY)
			}
			if sat.OBJECT_TYPE != "" {
				info += fmt.Sprintf(" [%s]", sat.OBJECT_TYPE)
			}
			satStrings = append(satStrings, info)
		}

		// Add navigation options
		var menuItems []string
		hasNextPage := false
		if searchName != "" {
			hasNextPage = page < totalPages
		} else {
			hasNextPage = len(sats) == pageSize
		}

		if page > 1 {
			menuItems = append(menuItems, "◄ Previous Page")
		}
		menuItems = append(menuItems, satStrings...)
		if hasNextPage {
			menuItems = append(menuItems, "Next Page ►")
		}
		menuItems = append(menuItems, "⭐ View Favorites", "🔍 New Search", "❌ Cancel")

		pageInfo := fmt.Sprintf("Page %d", page)
		if searchName != "" && totalPages > 0 {
			pageInfo += fmt.Sprintf(" of %d", totalPages)
		}
		if len(sats) == pageSize && hasNextPage {
			pageInfo += " (showing 20 results)"
		} else {
			pageInfo += fmt.Sprintf(" (%d results)", len(sats))
		}

		prompt := promptui.Select{
			Label: fmt.Sprintf("Select a Satellite 🛰 - %s", pageInfo),
			Items: menuItems,
			Size:  15,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			fmt.Println(color.Ize(color.Red, "  [!] PROMPT FAILED"))
			return ""
		}

		// Handle navigation
		if page > 1 && idx == 0 {
			// Previous Page
			page--
			continue
		}

		startIdx := 0
		if page > 1 {
			startIdx = 1
		}

		if idx >= startIdx && idx < startIdx+len(satStrings) {
			// Selected a satellite - extract just the name and NORAD ID for compatibility
			selectedIdx := idx - startIdx
			selectedSat := sats[selectedIdx]
			result := fmt.Sprintf("%s (%s)", selectedSat.SATNAME, selectedSat.NORAD_CAT_ID)

			// Check if already in favorites and offer to save/remove
			isFav, _ := IsFavorite(selectedSat.NORAD_CAT_ID)
			if !isFav {
				savePrompt := promptui.Prompt{
					Label:     fmt.Sprintf("Save %s to favorites? (y/n)", selectedSat.SATNAME),
					Default:   "n",
					AllowEdit: true,
				}
				saveAnswer, _ := savePrompt.Run()
				if strings.ToLower(strings.TrimSpace(saveAnswer)) == "y" {
					if err := AddFavorite(selectedSat.SATNAME, selectedSat.NORAD_CAT_ID, selectedSat.COUNTRY, selectedSat.OBJECT_TYPE); err != nil {
						fmt.Println(color.Ize(color.Yellow, "  [!] "+err.Error()))
					} else {
						fmt.Println(color.Ize(color.Green, fmt.Sprintf("  [+] Saved %s to favorites", selectedSat.SATNAME)))
					}
				}
			}

			return result
		}

		nextPageIdx := startIdx + len(satStrings)
		if idx == nextPageIdx && hasNextPage {
			// Next Page
			page++
			continue
		}

		favoritesIdx := nextPageIdx
		if hasNextPage {
			favoritesIdx++
		}
		newSearchIdx := favoritesIdx + 1

		if idx == favoritesIdx {
			// View Favorites
			favResult := SelectFromFavorites()
			if favResult != "" {
				return favResult
			}
			// Continue showing current page
			continue
		}

		if idx == newSearchIdx || (idx == favoritesIdx && !hasNextPage) {
			// New Search - reset cache
			allFilteredSats = []Satellite{}
			searchName, country, objectType, launchYear = showSearchMenu()
			page = 1
			totalPages = 0
			continue
		}

		// Cancel
		return ""
	}
}

// GenRowString formats a key-value pair into a table row with proper spacing.
func GenRowString(intro string, input string) string {
	var totalCount int = 4 + len(intro) + len(input) + 2
	var useCount = 63 - totalCount
	return "║ " + intro + ": " + input + strings.Repeat(" ", useCount) + " ║"
}

// Option prompts the user for a numeric input within a specified range.
// Returns the selected number, or exits the program if the minimum value is chosen.
func Option(min int, max int) int {
	fmt.Print("\n ENTER INPUT > ")
	var selection string
	fmt.Scanln(&selection)
	num, err := strconv.Atoi(selection)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] INVALID INPUT"))
		return Option(min, max)
	} else {
		if num == min {
			fmt.Println(color.Ize(color.Blue, " Escaping Orbit..."))
			os.Exit(1)
			return 0
		} else if num > min && num < max+1 {
			return num
		} else {
			fmt.Println(color.Ize(color.Red, "  [!] INVALID INPUT"))
			return Option(min, max)
		}
	}
}
