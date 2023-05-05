package osint

import (
	"fmt"
	"io/ioutil"
	"github.com/iskaa02/qalam"
)

func Banner() {
	banner, _ := ioutil.ReadFile("txt/banner.txt")
	fmt.Println("\033[37m", string(banner))

	s:=qalam.NewStyler().
		Bold().
		Italic().
		Red()
	// This will print bold italic red text to the terminal
	s.Print("Hello world")
}