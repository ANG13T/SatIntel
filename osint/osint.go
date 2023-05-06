package osint

import (
	"fmt"
	"bufio"
	"os"
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

}

func TLEPlainString(string input){

}