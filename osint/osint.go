package osint

import (
	"fmt"
	"bufio"
	"os"
	"github.com/TwiN/go-color"
	"github.com/iskaa02/qalam/gradient"	
	"io/ioutil"
	"strconv"
)

type TLE struct {
	CommonName string
	SatelliteCatalogNumber int
	ElsetClassificiation string
	InternationalDesignator string
	ElementSetEpoch float32
	FirstDerivativeMeanMotion float32
	SecondDerivativeMeanMotion string
	BDragTerm string
	ElementSetType int
	ElementNumber int
	ChecksumOne int
	OrbitInclination float32
	RightAscension float32
	Eccentrcity float32
	Perigee float32
	MeanAnamoly float32
	MeanMotion float32
	RevolutionNumber int
	ChecksumTwo int
}

func TLEParser(){
	options, _ := ioutil.ReadFile("txt/tle_parser.txt")
	opt,_:=gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 3)
	if (selection == 3) {
		return
	}
}

func TLETextFile(){

	fmt.Print("\n ENTER INPUT > ")
	var path string
	file, err := os.Open(path)
 
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] INVALID TEXT FILE"))
	}
 
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
 
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
 
	file.Close()
 
	for _, eachline := range txtlines {
		fmt.Println(eachline)
	}
}

func Option(min int, max int) int {
	fmt.Print("\n ENTER INPUT > ")
	var selection string
	fmt.Scanln(&selection)
	num, err := strconv.Atoi(selection)
    if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] INVALID INPUT"))
		return Option(min, max)
    } else {
        if (num >= min  && num < max + 1) {
			return num
		} else {
			fmt.Println(color.Ize(color.Red, "  [!] INVALID INPUT"))
			return Option(min, max)
		}
    }
}

func TLEPlainString(input string){

}