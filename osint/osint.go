package osint

import (
	"fmt"
	"io/ioutil"
	"github.com/iskaa02/qalam/gradient"
	"github.com/TwiN/go-color"
	"strconv"
)

func Option() {
	fmt.Print("\n ENTER INPUT > ")
	var selection string
	fmt.Scanln(&selection)
	num, err := strconv.Atoi(selection)
    if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] INVALID INPUT"))
		Option()
    } else {
        if (num >= 0  && num < 5) {
			DisplayFunctions()
		} else {
			fmt.Println(color.Ize(color.Red, "  [!] INVALID INPUT"))
			Option()
		}
    }
}

func DisplayFunctions() {
	fmt.Print("heh")
}

func Banner() {
	banner, _ := ioutil.ReadFile("txt/banner.txt")
	info, _ := ioutil.ReadFile("txt/info.txt")
	options, _ := ioutil.ReadFile("txt/options.txt")
	g,_:=gradient.NewGradient("cyan", "blue")
	solid,_:=gradient.NewGradient("blue", "#1179ef")
	opt,_:=gradient.NewGradient("#1179ef", "cyan")
	g.Print(string(banner))
	solid.Print(string(info))
	opt.Print("\n" + string(options))
}

func SatIntel() {
	Banner()
	Option()
}
