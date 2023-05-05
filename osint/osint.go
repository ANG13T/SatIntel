package osint

import (
	"io/ioutil"
	"github.com/iskaa02/qalam/gradient"
)

func Banner() {
	banner, _ := ioutil.ReadFile("txt/banner.txt")
	info, _ := ioutil.ReadFile("txt/info.txt")
	options, _ := ioutil.ReadFile("txt/options.txt")
	g,_:=gradient.NewGradient("cyan", "blue")
	g.Print(string(banner))
	g.Print(string(info))
	g.Print("\n" + string(options))
}