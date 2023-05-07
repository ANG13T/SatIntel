package osint

import (
	"fmt"
	"github.com/TwiN/go-color"
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

func ConstructTLE(one string, two string, three string) TLE {
	tle := TLE{}
	tle.CommonName = one
	firstArr := strings.Fields(two)
	secondArr := strings.Fields(three)
	tle.SatelliteCatalogNumber, _ = strconv.Atoi(firstArr[1][:len(firstArr[1])-1])
	tle.ElsetClassificiation = string(firstArr[1][len(firstArr[1])-1])
	tle.InternationalDesignator = firstArr[2]
	tle.ElementSetEpoch, _ = strconv.ParseFloat(firstArr[3], 64)
	tle.FirstDerivativeMeanMotion, _ = strconv.ParseFloat(firstArr[4], 64)
	tle.SecondDerivativeMeanMotion = firstArr[5]
	tle.BDragTerm = firstArr[6]
	tle.ElementSetType, _ = strconv.Atoi(firstArr[7])
	tle.ElementNumber, _ = strconv.Atoi(firstArr[8][:len(firstArr[8])-1])
	tle.ChecksumOne, _ = strconv.Atoi(string(firstArr[8][len(firstArr[8])-1]))
	tle.SatelliteCatalogNumber, _ = strconv.Atoi(secondArr[1])
	tle.OrbitInclination, _ = strconv.ParseFloat(secondArr[2], 64)
	tle.RightAscension, _ = strconv.ParseFloat(secondArr[3], 64)
	tle.Eccentrcity, _ = strconv.ParseFloat("0." + secondArr[4], 64)
	tle.Perigee, _ = strconv.ParseFloat(secondArr[5], 64)
	tle.MeanAnamoly, _ = strconv.ParseFloat(secondArr[6], 64)
	tle.MeanMotion, _ = strconv.ParseFloat(secondArr[7][:11], 64)
	tle.RevolutionNumber, _ = strconv.Atoi(secondArr[7][11:16])
	tle.ChecksumTwo, _ = strconv.Atoi(string(secondArr[7][len(secondArr[7])-1]))

	return tle
}

func PrintTLE (tle TLE) {
	fmt.Println(color.Ize(color.Purple, "\n╔═════════════════════════════════════════════════════════════╗"))
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
	
	fmt.Println(color.Ize(color.Purple, "╚═════════════════════════════════════════════════════════════╝ \n\n"))
}