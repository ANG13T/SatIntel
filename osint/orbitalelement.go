package osint

import (
	"fmt"
	"os"

	"github.com/iskaa02/qalam/gradient"
)

// OrbitalElement displays orbital element data for a selected satellite.
func OrbitalElement() {
	options, _ := os.ReadFile("txt/orbital_element.txt")
	opt, _ := gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 3)

	if selection == 1 {
		result := SelectSatellite()

		if result == "" {
			return
		}

		PrintNORADInfo(extractNorad(result), result)

	} else if selection == 2 {
		fmt.Print("\n ENTER NORAD ID > ")
		var norad string
		fmt.Scanln(&norad)
		PrintNORADInfo(norad, "UNSPECIFIED")
	}
}
