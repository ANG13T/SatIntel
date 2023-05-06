package osint

import (
	"fmt"
	"bufio"
	"os"
	"github.com/TwiN/go-color"
)

type struct TLE {
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
	fmt.Println("Parsing TLE")

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

func TLEPlainString(string input){

}