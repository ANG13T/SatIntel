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

	PrintTLE(output)
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
	scanner := bufio.NewScanner(os.Stdin)
	var lineOne string
	var lineTwo string
	var lineThree string
	fmt.Print("\n ENTER LINE ONE (leave blank for unspecified name)  >  ")
	scanner.Scan()
    lineOne = scanner.Text()

	fmt.Print("\n ENTER LINE TWO  >  ")
	scanner.Scan()
    lineTwo = scanner.Text()

	fmt.Print("\n ENTER LINE THREE  >  ")
	scanner.Scan()
    lineThree = scanner.Text()

	if (lineOne == "") {
		lineOne = "UNSPECIFIED"
	}
	
	output := TLE{}

	output = ConstructTLE(lineOne, lineTwo, lineThree)

	PrintTLE(output)
}

// TODO: Right Ascension of Ascending Node (degrees)

func PrintTLE (tle TLE) {
	fmt.Println(color.Ize(color.Purple, "╔═════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Purple, GenRowString("Name", tle.CommonName)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Satellite Catalog Number", fmt.Sprintf("%d", tle.SatelliteCatalogNumber))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Elset Classification", tle.ElsetClassificiation)))
	fmt.Println(color.Ize(color.Purple, GenRowString("International Designator", tle.InternationalDesignator)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Element Set Epoch (UTC)", fmt.Sprintf("%f", tle.ElementSetEpoch))))
	fmt.Println(color.Ize(color.Purple, GenRowString("1st Derivative of the Mean Motion", fmt.Sprintf("%f", tle.FirstDerivativeMeanMotion))))
	fmt.Println(color.Ize(color.Purple, GenRowString("2nd Derivative of the Mean Motion", tle.SecondDerivativeMeanMotion)))
	fmt.Println(color.Ize(color.Purple, GenRowString("B* Drag Term", tle.BDragTerm)))
	fmt.Println(color.Ize(color.Purple, GenRowString("Element Set Type", fmt.Sprintf("%d", tle.ElementSetType))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Element Number", fmt.Sprintf("%d", tle.ElementNumber))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Checksum Line One", fmt.Sprintf("%d", tle.ChecksumOne))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Orbit Inclination (degrees)", fmt.Sprintf("%f", tle.OrbitInclination))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Right Ascension of Ascending Node (degrees)", fmt.Sprintf("%f", tle.RightAscension))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Eccentricity", fmt.Sprintf("%f", tle.Eccentrcity))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Argument of Perigee (degrees)", fmt.Sprintf("%f", tle.Perigee))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Mean Anomaly (degrees)", fmt.Sprintf("%f", tle.MeanAnamoly))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Mean Motion (revolutions/day)", fmt.Sprintf("%f", tle.MeanMotion))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Revolution Number at Epoch", fmt.Sprintf("%d", tle.RevolutionNumber))))
	fmt.Println(color.Ize(color.Purple, GenRowString("Checksum Line Two", fmt.Sprintf("%d", tle.ChecksumTwo))))
	
	fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝ \n"))
}

func GenRowString(intro string, input string) string{
	var totalCount int = 4 + len(intro) + len(input) + 2
	var useCount = 63 - totalCount
	return "║ " + intro + ": " + input + strings.Repeat(" ", useCount) + " ║"
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