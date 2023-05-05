package osint

import (
	"io/ioutil"
	"github.com/iskaa02/qalam/gradient"
)

func Banner() {
	banner, _ := ioutil.ReadFile("txt/banner.txt")
	info, _ := ioutil.ReadFile("txt/info.txt")
	g,_:=gradient.NewGradient("cyan", "blue")
	g.Print(string(banner))
	g.Print(string(info))
}