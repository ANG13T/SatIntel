package osint

import (
	"io/ioutil"
	"fmt"
	"github.com/iskaa02/qalam/gradient"	
)

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

		PrintNORADInfo(extractNorad(result), result)

	} else if (selection == 2) {
		fmt.Print("\n ENTER NORAD ID > ")
		var norad string
		fmt.Scanln(&norad)
		PrintNORADInfo(norad, "UNSPECIFIED")
	} 

	return
}