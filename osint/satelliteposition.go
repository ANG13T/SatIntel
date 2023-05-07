package osint

import (
	"io/ioutil"
	"fmt"
	"github.com/iskaa02/qalam/gradient"	
	"encoding/json"
	"github.com/TwiN/go-color"
    "net/http"
	"strconv"
	"os"
)

func SatellitePositionVisualization() {
	options, _ := ioutil.ReadFile("txt/orbital_element.txt")
	opt,_:=gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 3)

	if (selection == 1) {
		result := SelectSatellite()

		if (result == "") {
			return
		}

		GetLocation(extractNorad(result))

	} else if (selection == 2) {
		fmt.Print("\n ENTER NORAD ID > ")
		var norad string
		fmt.Scanln(&norad)
		GetLocation(norad)
	} 

	return
}

// Show visualization and info in box
func GetLocation(norad string) {
	fmt.Print("\n ENTER LATITUDE > ")
	var latitude string
	fmt.Scanln(&latitude)
	fmt.Print("\n ENTER LONGITUDE > ")
	var longitude string
	fmt.Scanln(&longitude)
	fmt.Print("\n ENTER ALTITUDE > ")
	var altitude string
	fmt.Scanln(&altitude)

	_, err := strconv.ParseFloat(latitude, 64)
	_, err2 := strconv.ParseFloat(longitude, 64)
	_, err3 := strconv.Atoi(altitude)

	if err != nil || err2 != nil || err3 != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: INVALID INPUT"))
		return
	}

	url := "https://api.n2yo.com/rest/v1/satellite/positions/" + norad + "/" + latitude + "/" + longitude + "/" + altitude + "/2/&apiKey=" + os.Getenv("N2YO_API_KEY")
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()

    var data Response
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        fmt.Println(err)
    }

	fmt.Println(color.Ize(color.Purple, "\n╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Purple, "║                    Satellite Information                    ║"))
	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))

	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Name", data.SatelliteInfo.Satname)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite ID",  fmt.Sprintf("%d", data.SatelliteInfo.Satid))))
	// fmt.Println(color.Ize(color.Purple, "║                                                             ║"))

	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
	fmt.Println(color.Ize(color.Purple, "║                     Satellite Positions                     ║"))
	fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))

    for in, pos := range data.Positions {
		PrintSatellitePosition(pos, in == len(data.Positions) - 1)
    }

}

func DisplayMap() {
	// TODO
}

func PrintSatellitePosition (pos Position, last bool) {
	fmt.Println(color.Ize(color.Purple, GenRowString("Latitude", fmt.Sprintf("%f", pos.Satlatitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Longitude", fmt.Sprintf("%f", pos.Satlongitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Altitude", fmt.Sprintf("%f", pos.Sataltitude))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Right Ascension", fmt.Sprintf("%f", pos.Azimuth))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Declination", fmt.Sprintf("%f", pos.Dec))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Timestamp", fmt.Sprintf("%d", pos.Timestamp))))
	if (last) {
		fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝\n\n"))
	} else {
		fmt.Println(color.Ize(color.Purple, "╠═════════════════════════════════════════════════════════════╣"))
	}
}