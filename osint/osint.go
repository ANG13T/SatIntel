package osint

import (
	"fmt"
	"bufio"
	"os"
	"github.com/TwiN/go-color"
	"github.com/iskaa02/qalam/gradient"	
	"io/ioutil"
	"strconv"
	"strings"
)

type TLE struct {
	CommonName string
	SatelliteCatalogNumber int
	ElsetClassificiation string
	InternationalDesignator string
	ElementSetEpoch float64
	FirstDerivativeMeanMotion float64
	SecondDerivativeMeanMotion string
	BDragTerm string
	ElementSetType int
	ElementNumber int
	ChecksumOne int
	OrbitInclination float64
	RightAscension float64
	Eccentrcity float64
	Perigee float64
	MeanAnamoly float64
	MeanMotion float64
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
	tle := TLE{}
	tle.CommonName = one
	firstArr := strings.Split(two, "")
	secondArr := strings.Split(three, "")
	tle.SatelliteCatalogNumber, _ = strconv.Atoi(firstArr[1])
	tle.ElsetClassificiation = firstArr[2]
	tle.InternationalDesignator = firstArr[3]
	tle.ElementSetEpoch, _ = strconv.ParseFloat(firstArr[4], 64)
	tle.FirstDerivativeMeanMotion, _ = strconv.ParseFloat(firstArr[5], 64)
	tle.SecondDerivativeMeanMotion = firstArr[6]
	tle.BDragTerm = firstArr[7]
	tle.ElementSetType, _ = strconv.Atoi(firstArr[8])
	tle.ElementNumber, _ = strconv.Atoi(firstArr[9])
	tle.ChecksumOne, _ = strconv.Atoi(firstArr[10])
	tle.OrbitInclination, _ = strconv.ParseFloat(secondArr[1], 64)
	tle.RightAscension, _ = strconv.ParseFloat(secondArr[2], 64)
	tle.Eccentrcity, _ = strconv.ParseFloat(secondArr[3], 64)
	tle.Perigee, _ = strconv.ParseFloat(secondArr[4], 64)
	tle.MeanAnamoly, _ = strconv.ParseFloat(secondArr[5], 64)
	tle.MeanMotion, _ = strconv.ParseFloat(secondArr[6], 64)
	tle.RevolutionNumber, _ = strconv.Atoi(secondArr[7])
	tle.ChecksumTwo, _ = strconv.Atoi(secondArr[8])

	return tle
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

