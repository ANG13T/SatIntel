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

func TLEParser() {
	options, _ := ioutil.ReadFile("txt/tle_parser.txt")
	opt,_:=gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 3)

	if (selection == 1){
		TLETextFile()
	} else if (selection == 2) {
		TLEPlainString()
	} 

	return
}

func TLETextFile() {

	fmt.Print("\n ENTER TEXT FILE PATH > ")
	var path string
	fmt.Scanln(&path)
	file, err := os.Open(path)
 
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] INVALID TEXT FILE"))
	}
 
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	var count int = 0
 
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
		count += 1
	}
 
	file.Close()

	if (count < 2 || count > 3) {
		fmt.Println(color.Ize(color.Red, "  [!] INVALID TLE FORMAT"))
		return
	}

	output := TLE{}

	if (count == 3) {
		var satelliteName string = txtlines[0]
		output = ConstructTLE(satelliteName, txtlines[1], txtlines[2])
	} else {
		output = ConstructTLE("UNSPECIFIED", txtlines[0], txtlines[1])
	}

	fmt.Println(count, output)
}

func ConstructTLE(one string, two string, three string) TLE {
	output := TLE{}
	output.CommonName = one
	return output
}

func TLEPlainString(){
	fmt.Print("\n ENTER TLE > ")
	var tleString string
	fmt.Scanln(&tleString)
}

func PrintTLE (tle TLE) {

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
		if (num == min) {
			fmt.Println(color.Ize(color.Blue, " Escaping Orbit..."))
			os.Exit(1)
			return 0
		} else if (num > min  && num < max + 1) {
			return num
		} else {
			fmt.Println(color.Ize(color.Red, "  [!] INVALID INPUT"))
			return Option(min, max)
		}
    }
}

