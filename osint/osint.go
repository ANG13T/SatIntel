package osint

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Banner() {
	banner, err := ioutil.ReadFile("../txt/banner.txt")
	res(err)
	fmt.Println("\033[37m", string(banner))
}