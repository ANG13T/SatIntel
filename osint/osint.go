package osint

import (
	"io/ioutil"
	"github.com/iskaa02/qalam/gradient"
)

func Banner() {
	banner, _ := ioutil.ReadFile("txt/banner.txt")
	g,err:=gradient.NewGradient("purple", "cyan", "blue")

	if err != nil{
	}
	g.Print(string(banner))
}