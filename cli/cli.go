package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/ANG13T/SatIntel/osint"
	"github.com/TwiN/go-color"
	"github.com/iskaa02/qalam/gradient"
)

// Option prompts the user to select a menu option and validates the input.
// It recursively prompts until a valid option between 0 and 4 is entered.
func Option() {
	fmt.Print("\n ENTER INPUT > ")
	var selection string
	fmt.Scanln(&selection)
	num, err := strconv.Atoi(selection)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] INVALID INPUT"))
		Option()
	} else {
		if num >= 0 && num < 5 {
			DisplayFunctions(num)
		} else {
			fmt.Println(color.Ize(color.Red, "  [!] INVALID INPUT"))
			Option()
		}
	}
}

// DisplayFunctions executes the selected function based on the menu choice.
// After execution, it waits for user input, clears the screen, and shows the menu again.
func DisplayFunctions(x int) {
	if x == 0 {
		fmt.Println(color.Ize(color.Blue, " Escaping Orbit..."))
		os.Exit(1)
	} else if x == 1 {
		osint.OrbitalElement()
		waitForEnter()
		clearScreen()
		Banner()
		Option()
	} else if x == 2 {
		osint.SatellitePositionVisualization()
		waitForEnter()
		clearScreen()
		Banner()
		Option()
	} else if x == 3 {
		osint.OrbitalPrediction()
		waitForEnter()
		clearScreen()
		Banner()
		Option()
	} else if x == 4 {
		osint.TLEParser()
		waitForEnter()
		clearScreen()
		Banner()
		Option()
	}
}

// Banner displays the application banner, info, and menu options with gradient colors.
func Banner() {
	banner, _ := os.ReadFile("txt/banner.txt")
	info, _ := os.ReadFile("txt/info.txt")
	options, _ := os.ReadFile("txt/options.txt")
	g, _ := gradient.NewGradient("cyan", "blue")
	solid, _ := gradient.NewGradient("blue", "#1179ef")
	opt, _ := gradient.NewGradient("#1179ef", "cyan")
	g.Print(string(banner))
	solid.Print(string(info))
	opt.Print("\n" + string(options))
}

// waitForEnter pauses execution and waits for the user to press Enter.
func waitForEnter() {
	fmt.Print("\n\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// clearScreen clears the terminal screen using the appropriate command for the operating system.
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// SatIntel initializes the SatIntel CLI application by displaying the banner and starting the menu loop.
func SatIntel() {
	Banner()
	Option()
}
