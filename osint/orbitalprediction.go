package osint

import (
	"fmt"
	"github.com/iskaa02/qalam/gradient"
	"io/ioutil"
)

func OrbitalPrediction() {
	options, _ := ioutil.ReadFile("txt/orbital_prediction.txt")
	opt,_:=gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 4)

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

	} else if (selection == 3) {

	}

	return
}