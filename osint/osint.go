package osint

import (
	"fmt"
	"bufio"
	"os"
	"github.com/TwiN/go-color"
	"github.com/iskaa02/qalam/gradient"	
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

// Orbital Element Data Display Code
func OrbitalElement() {
	options, _ := ioutil.ReadFile("txt/orbital_element.txt")
	opt,_:=gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 3)

	if (selection == 1) {
		vals := url.Values{}
		vals.Add("identity", os.Getenv("SPACE_TRACK_USERNAME"))
		vals.Add("password", os.Getenv("SPACE_TRACK_PASSWORD"))
		vals.Add("query", "https://www.space-track.org/basicspacedata/query/class/satcat/orderby/SATNAME asc/limit/10/emptyresult/show")
	
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
	
		var sats []Satellite
		if err := json.Unmarshal(respData, &sats); err != nil {
			fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
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
			return
		}
		PrintNORADInfo(extractNorad(result))

	} else if (selection == 2) {
		fmt.Print("\n ENTER NORAD ID > ")
		var norad string
		fmt.Scanln(&norad)
		PrintNORADInfo(norad)
	} 

	return
}

func extractNorad(str string) string {
    start := strings.Index(str, "(")
    end := strings.Index(str, ")")
    if start == -1 || end == -1 || start >= end {
        return ""
    }
    return str[start+1:end]
}

func PrintNORADInfo(norad string) {
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
	tle := ConstructTLE("UNSPECIFIED", lineOne, lineTwo)
	PrintTLE(tle)
}

// Satellite Position Visualization Code
func SatellitePositionVisualization() {
	options, _ := ioutil.ReadFile("txt/orbital_element.txt")
	opt,_:=gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 3)

	if (selection == 1) {
		vals := url.Values{}
		vals.Add("identity", os.Getenv("SPACE_TRACK_USERNAME"))
		vals.Add("password", os.Getenv("SPACE_TRACK_PASSWORD"))
		vals.Add("query", "https://www.space-track.org/basicspacedata/query/class/satcat/orderby/SATNAME asc/limit/10/emptyresult/show")
	
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
	
		var sats []Satellite
		if err := json.Unmarshal(respData, &sats); err != nil {
			fmt.Println(color.Ize(color.Red, "  [!] ERROR: API REQUEST TO SPACE TRACK"))
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
			return
		}
		PrintNORADInfo(extractNorad(result))

	} else if (selection == 2) {
		fmt.Print("\n ENTER NORAD ID > ")
		var norad string
		fmt.Scanln(&norad)
		PrintNORADInfo(norad)
	} 

	return
}

// TLE Parser Code

func TLEParser() {
	options, _ := ioutil.ReadFile("txt/tle_parser.txt")
	opt,_:=gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 3)

	if (selection == 1){
		TLETextFile()
	} else if (selection == 2) {
		TLEPlainString()
	} 

	return
}

func TLETextFile() {

	fmt.Print("\n ENTER TEXT FILE PATH > ")
	var path string
	fmt.Scanln(&path)
	file, err := os.Open(path)
 
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] INVALID TEXT FILE"))
	}
 
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	var count int = 0
 
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
		count += 1
	}
 
	file.Close()

	if (count < 2 || count > 3) {
		fmt.Println(color.Ize(color.Red, "  [!] INVALID TLE FORMAT"))
		return
	}

	output := TLE{}

	if (count == 3) {
		var satelliteName string = txtlines[0]
		output = ConstructTLE(satelliteName, txtlines[1], txtlines[2])
	} else {
		output = ConstructTLE("UNSPECIFIED", txtlines[0], txtlines[1])
	}

	PrintTLE(output)
}

func ConstructTLE(one string, two string, three string) TLE {
	tle := TLE{}
	tle.CommonName = one
	firstArr := strings.Fields(two)
	secondArr := strings.Fields(three)
	tle.SatelliteCatalogNumber, _ = strconv.Atoi(firstArr[1][:len(firstArr[1])-1])
	tle.ElsetClassificiation = string(firstArr[1][len(firstArr[1])-1])
	tle.InternationalDesignator = firstArr[2]
	tle.ElementSetEpoch, _ = strconv.ParseFloat(firstArr[3], 64)
	tle.FirstDerivativeMeanMotion, _ = strconv.ParseFloat(firstArr[4], 64)
	tle.SecondDerivativeMeanMotion = firstArr[5]
	tle.BDragTerm = firstArr[6]
	tle.ElementSetType, _ = strconv.Atoi(firstArr[7])
	tle.ElementNumber, _ = strconv.Atoi(firstArr[8][:len(firstArr[8])-1])
	tle.ChecksumOne, _ = strconv.Atoi(string(firstArr[8][len(firstArr[8])-1]))
	tle.SatelliteCatalogNumber, _ = strconv.Atoi(secondArr[1])
	tle.OrbitInclination, _ = strconv.ParseFloat(secondArr[2], 64)
	tle.RightAscension, _ = strconv.ParseFloat(secondArr[3], 64)
	tle.Eccentrcity, _ = strconv.ParseFloat("0." + secondArr[4], 64)
	tle.Perigee, _ = strconv.ParseFloat(secondArr[5], 64)
	tle.MeanAnamoly, _ = strconv.ParseFloat(secondArr[6], 64)
	tle.MeanMotion, _ = strconv.ParseFloat(secondArr[7][:11], 64)
	tle.RevolutionNumber, _ = strconv.Atoi(secondArr[7][11:16])
	tle.ChecksumTwo, _ = strconv.Atoi(string(secondArr[7][len(secondArr[7])-1]))

	return tle
}

func TLEPlainString(){
	scanner := bufio.NewScanner(os.Stdin)
	var lineOne string
	var lineTwo string
	var lineThree string
	fmt.Print("\n ENTER LINE ONE (leave blank for unspecified name)  >  ")
	scanner.Scan()
    lineOne = scanner.Text()

	fmt.Print("\n ENTER LINE TWO  >  ")
	scanner.Scan()
    lineTwo = scanner.Text()

	fmt.Print("\n ENTER LINE THREE  >  ")
	scanner.Scan()
    lineThree = scanner.Text()

	if (lineOne == "") {
		lineOne = "UNSPECIFIED"
	}
	
	output := TLE{}

	output = ConstructTLE(lineOne, lineTwo, lineThree)

	PrintTLE(output)
}

// TODO: Right Ascension of Ascending Node (degrees)

func PrintTLE (tle TLE) {
	fmt.Println(color.Ize(color.Purple, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Purple, GenRowString("Name", tle.CommonName)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Catalog Number", fmt.Sprintf("%d", tle.SatelliteCatalogNumber))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Elset Classification", tle.ElsetClassificiation)))
	fmt.Println(color.Ize(color.Purple, GenRowString("International Designator", tle.InternationalDesignator)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Element Set Epoch (UTC)", fmt.Sprintf("%f", tle.ElementSetEpoch))))
	fmt.Println(color.Ize(color.Purple, GenRowString("1st Derivative of the Mean Motion", fmt.Sprintf("%f", tle.FirstDerivativeMeanMotion))))
	fmt.Println(color.Ize(color.Purple, GenRowString("2nd Derivative of the Mean Motion", tle.SecondDerivativeMeanMotion)))
	fmt.Println(color.Ize(color.Purple, GenRowString("B* Drag Term", tle.BDragTerm)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Element Set Type", fmt.Sprintf("%d", tle.ElementSetType))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Element Number", fmt.Sprintf("%d", tle.ElementNumber))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Checksum Line One", fmt.Sprintf("%d", tle.ChecksumOne))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Orbit Inclination (degrees)", fmt.Sprintf("%f", tle.OrbitInclination))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Right Ascension of Ascending Node (degrees)", fmt.Sprintf("%f", tle.RightAscension))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Eccentricity", fmt.Sprintf("%f", tle.Eccentrcity))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Argument of Perigee (degrees)", fmt.Sprintf("%f", tle.Perigee))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Mean Anomaly (degrees)", fmt.Sprintf("%f", tle.MeanAnamoly))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Mean Motion (revolutions/day)", fmt.Sprintf("%f", tle.MeanMotion))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Revolution Number at Epoch", fmt.Sprintf("%d", tle.RevolutionNumber))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Checksum Line Two", fmt.Sprintf("%d", tle.ChecksumTwo))))
	
	fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝ \n\n"))
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