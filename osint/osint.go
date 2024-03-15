package osint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/manifoldco/promptui"
)

const authURL = "https://www.space-track.org/ajaxauth/login"

func extractNorad(str string) string {
	start := strings.Index(str, "(")
	end := strings.Index(str, ")")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return str[start+1 : end]
}

func PrintNORADInfo(norad string, name string) {
	vals := url.Values{}
	vals.Add("identity", os.Getenv("SPACE_TRACK_USERNAME"))
	vals.Add("password", os.Getenv("SPACE_TRACK_PASSWORD"))
	vals.Add("query", "https://www.space-track.org/basicspacedata/query/class/gp_history/format/tle/NORAD_CAT_ID/"+norad+"/orderby/EPOCH%20desc/limit/1")

	client := &http.Client{}

	resp, err := client.PostForm(authURL, vals)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
	}
	respData, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
	}

	tleLines := strings.Fields(string(respData))
	mid := (len(tleLines) / 2) + 1
	lineOne := strings.Join(tleLines[:mid], " ")
	lineTwo := strings.Join(tleLines[mid:], " ")
	tle := ConstructTLE(name, lineOne, lineTwo)
	PrintTLE(tle)
}

func fetchData(page int, pageSize int) ([]Satellite, error) {
	vals := url.Values{}
	vals.Add("identity", os.Getenv("SPACE_TRACK_USERNAME"))
	vals.Add("password", os.Getenv("SPACE_TRACK_PASSWORD"))
	vals.Add("query",
		fmt.Sprintf("https://www.space-track.org/basicspacedata/query/class/satcat/orderby/SATNAME asc/limit/%d,%d/emptyresult/show",
			pageSize,
			page*pageSize,
		),
	)
	client := &http.Client{}

	resp, err := client.PostForm(authURL, vals)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, err
	}

	var sats []Satellite
	if err := json.NewDecoder(resp.Body).Decode(&sats); err != nil {
		return nil, err
	}

	return sats, nil
}

func SelectSatellite() string {
	const pageSize = 10
	page := 0
	var satStrings []string

	for {
		items, err := fetchData(page, pageSize)
		if err != nil {
			fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
			break
		}
		for _, sat := range items {
			satStrings = append(satStrings, fmt.Sprintf("%s (%s)", sat.SATNAME, sat.NORAD_CAT_ID))
		}
		prompt := promptui.Select{
			Label:        "Select a Satellite 🛰",
			Items:        append(satStrings, fmt.Sprintf("Load next %d results", pageSize)),
			Size:         5,
			HideSelected: true,
		}
		cursorPos := len(satStrings) - pageSize
		index, result, err := prompt.RunCursorAt(cursorPos, cursorPos-prompt.Size+1)
		if err != nil {
			fmt.Println(color.Ize(color.Red, "  [!] PROMPT FAILED"))
			break
		}
		if index < len(satStrings) {
			return result
		}
		page++
	}
	return ""
}

func GenRowString(intro string, input string) string {
	var totalCount int = 4 + len(intro) + len(input) + 2
	var useCount = 63 - totalCount
	return "║ " + intro + ": " + input + strings.Repeat(" ", useCount) + " ║"
}

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
