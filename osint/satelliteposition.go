package osint

import (
	"io/ioutil"
	"fmt"
	"github.com/iskaa02/qalam/gradient"	
	"encoding/json"
    "net/http"
)

type Marker struct {
	X int
	Y int
}


// Satellite Position Visualization Code
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

		// PrintNORADInfo(extractNorad(result), result)
		GetLocation(extractNorad(result))

	} else if (selection == 2) {
		fmt.Print("\n ENTER NORAD ID > ")
		var norad string
		fmt.Scanln(&norad)
		PrintNORADInfo(norad, "UNSPECIFIED")
	} 

	return
}

// Show visualization and info in box
func GetLocation(norad string) {
	// fmt.Print("\n ENTER LATITUDE > ")
	// var latitude string
	// fmt.Print("\n ENTER LONGITUDE > ")
	// var longitude string
	// fmt.Print("\n ENTER ALTITUDE > ")
	// var altitude string

	// marker := getXYfromLonLat(strconv.Atoi(latitude, longitude))

	url := "https://api.n2yo.com/rest/v1/satellite/positions/25544/41.702/-76.014/0/2/&apiKey=46REG9-PS2V7M-H3B76Q-5103"
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

    fmt.Printf("Satellite Name: %s\n", data.SatelliteInfo.Satname)
    fmt.Printf("Satellite ID: %d\n", data.SatelliteInfo.Satid)

    for _, pos := range data.Positions {
        fmt.Printf("Latitude: %f\n", pos.Satlatitude)
        fmt.Printf("Longitude: %f\n", pos.Satlongitude)
        fmt.Printf("Altitude: %f\n", pos.Sataltitude)
        fmt.Printf("Azimuth: %f\n", pos.Azimuth)
        fmt.Printf("Elevation: %f\n", pos.Elevation)
        fmt.Printf("RA: %f\n", pos.Ra)
        fmt.Printf("Dec: %f\n", pos.Dec)
        fmt.Printf("Timestamp: %d\n", pos.Timestamp)
        fmt.Println("------------------------")
    }

}

func DisplayMap() {

}

func getXYfromLonLat(lat int, lon int) Marker {
	marker := Marker{}
    // Normalise the X, Y in their min -> max space
    normalX := (lat + 219) / (293 + 219);
    normalY := (lon + 244) / (266 + 244);

    // Stretch them to match the ASCII map
	const width = 70
	const height = 42
    realX := normalX * width - 3;
    realY := normalY * height - 9;

	marker.X = realX
	marker.Y = realY

    return marker
}
