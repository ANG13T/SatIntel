package osint

import (
	"fmt"
	"os"
	"github.com/TwiN/go-color"
	"io/ioutil"
	"strconv"
	"strings"
	"encoding/json"
	"net/http"	
	"github.com/manifoldco/promptui"
	"net/url"
)

type field string
type orderBy string
const authURL = "https://www.space-track.org/ajaxauth/login"
const baseurl = "https://www.space-track.org/basicspacedata/query/class"


func extractNorad(str string) string {
    start := strings.Index(str, "(")
    end := strings.Index(str, ")")
    if start == -1 || end == -1 || start >= end {
        return ""
    }
    return str[start+1:end]
}

func PrintNORADInfo(norad string, name string) {
	vals := url.Values{}
	vals.Add("identity", os.Getenv("SPACE_TRACK_USERNAME"))
	vals.Add("password", os.Getenv("SPACE_TRACK_PASSWORD"))
	vals.Add("query", "https://www.space-track.org/basicspacedata/query/class/gp_history/format/tle/NORAD_CAT_ID/" + norad + "/orderby/EPOCH%20desc/limit/1")

	client := &http.Client{}

	resp, err := client.PostForm(authURL, vals)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
	}
	respData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
	}

	tleLines := strings.Fields(string(respData))
	mid := (len(tleLines)/2) + 1
	lineOne := strings.Join(tleLines[:mid], " ")
	lineTwo := strings.Join(tleLines[mid:], " ")
	tle := ConstructTLE(name, lineOne, lineTwo)
	PrintTLE(tle)
}

func SelectSatellite() string {
	vals := url.Values{}
	vals.Add("identity", os.Getenv("SPACE_TRACK_USERNAME"))
	vals.Add("password", os.Getenv("SPACE_TRACK_PASSWORD"))
	vals.Add("query", "https://www.space-track.org/basicspacedata/query/class/satcat/orderby/SATNAME asc/limit/10/emptyresult/show")

	client := &http.Client{}

	resp, err := client.PostForm(authURL, vals)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
		return ""
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
		return ""
	}
	respData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
		return ""
	}

	var sats []Satellite
	if err := json.Unmarshal(respData, &sats); err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
		return ""
	}

	var satStrings []string
	for _, sat := range sats {
		satStrings = append(satStrings, sat.SATNAME + " (" + sat.NORAD_CAT_ID + ")")
	}
	prompt := promptui.Select{
		Label: "Select a Satellite 🛰",
		Items: satStrings,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] PROMPT FAILED"))
		return ""
	}
	return result
}

func GenRowString(intro string, input string) string{
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
		if (num == min) {
			fmt.Println(color.Ize(color.Blue, " Escaping Orbit..."))
			os.Exit(1)
			return 0
		} else if (num > min  && num < max + 1) {
			return num
		} else {
			fmt.Println(color.Ize(color.Red, "  [!] INVALID INPUT"))
			return Option(min, max)
		}
    }
}