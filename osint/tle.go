package osint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/manifoldco/promptui"
)

type TLE struct {
	CommonName                 string
	SatelliteCatalogNumber     int
	ElsetClassificiation       string
	InternationalDesignator    string
	ElementSetEpoch            float64
	FirstDerivativeMeanMotion  float64
	SecondDerivativeMeanMotion string
	BDragTerm                  string
	ElementSetType             int
	ElementNumber              int
	ChecksumOne                int
	OrbitInclination           float64
	RightAscension             float64
	Eccentrcity                float64
	Perigee                    float64
	MeanAnamoly                float64
	MeanMotion                 float64
	RevolutionNumber           int
	ChecksumTwo                int
}

// ConstructTLE parses two-line element data into a TLE struct.
// It handles variable field counts gracefully and returns an empty TLE if parsing fails.
func ConstructTLE(one string, two string, three string) TLE {
	tle := TLE{}
	tle.CommonName = one
	firstArr := strings.Fields(two)
	secondArr := strings.Fields(three)

	if len(firstArr) < 4 {
		return tle
	}
	if len(secondArr) < 3 {
		return tle
	}

	if len(firstArr) > 1 && len(firstArr[1]) > 0 {
		catalogStr := firstArr[1]
		if len(catalogStr) > 1 {
			tle.SatelliteCatalogNumber, _ = strconv.Atoi(catalogStr[:len(catalogStr)-1])
			tle.ElsetClassificiation = string(catalogStr[len(catalogStr)-1])
		} else {
			tle.SatelliteCatalogNumber, _ = strconv.Atoi(catalogStr)
		}
	}
	if len(firstArr) > 2 {
		tle.InternationalDesignator = firstArr[2]
	}
	if len(firstArr) > 3 {
		tle.ElementSetEpoch, _ = strconv.ParseFloat(firstArr[3], 64)
	}
	if len(firstArr) > 4 {
		tle.FirstDerivativeMeanMotion, _ = strconv.ParseFloat(firstArr[4], 64)
	}
	if len(firstArr) > 5 {
		tle.SecondDerivativeMeanMotion = firstArr[5]
	}
	if len(firstArr) > 6 {
		tle.BDragTerm = firstArr[6]
	}
	if len(firstArr) > 7 {
		tle.ElementSetType, _ = strconv.Atoi(firstArr[7])
	}
	if len(firstArr) > 8 {
		lastField := firstArr[8]
		if len(lastField) > 1 {
			tle.ElementNumber, _ = strconv.Atoi(lastField[:len(lastField)-1])
			tle.ChecksumOne, _ = strconv.Atoi(string(lastField[len(lastField)-1]))
		} else if len(lastField) > 0 {
			tle.ElementNumber, _ = strconv.Atoi(lastField)
		}
	}

	if len(secondArr) > 1 {
		tle.SatelliteCatalogNumber, _ = strconv.Atoi(secondArr[1])
	}
	if len(secondArr) > 2 {
		tle.OrbitInclination, _ = strconv.ParseFloat(secondArr[2], 64)
	}
	if len(secondArr) > 3 {
		tle.RightAscension, _ = strconv.ParseFloat(secondArr[3], 64)
	}
	if len(secondArr) > 4 {
		tle.Eccentrcity, _ = strconv.ParseFloat("0."+secondArr[4], 64)
	}
	if len(secondArr) > 5 {
		tle.Perigee, _ = strconv.ParseFloat(secondArr[5], 64)
	}
	if len(secondArr) > 6 {
		tle.MeanAnamoly, _ = strconv.ParseFloat(secondArr[6], 64)
	}
	if len(secondArr) > 7 {
		meanMotionStr := secondArr[7]
		if len(meanMotionStr) >= 11 {
			tle.MeanMotion, _ = strconv.ParseFloat(meanMotionStr[:11], 64)
			if len(meanMotionStr) >= 16 {
				tle.RevolutionNumber, _ = strconv.Atoi(meanMotionStr[11:16])
			}
		} else {
			tle.MeanMotion, _ = strconv.ParseFloat(meanMotionStr, 64)
		}
		if len(meanMotionStr) > 0 {
			tle.ChecksumTwo, _ = strconv.Atoi(string(meanMotionStr[len(meanMotionStr)-1]))
		}
	}

	return tle
}

// PrintTLE displays the TLE data in a formatted table.
func PrintTLE(tle TLE) {
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

	// Offer export option
	exportPrompt := promptui.Prompt{
		Label:     "Export TLE data? (y/n)",
		Default:   "n",
		AllowEdit: true,
	}
	exportAnswer, _ := exportPrompt.Run()
	if strings.ToLower(strings.TrimSpace(exportAnswer)) == "y" {
		defaultFilename := fmt.Sprintf("tle_%s_%d", strings.ReplaceAll(tle.CommonName, " ", "_"), tle.SatelliteCatalogNumber)
		format, filePath, err := showExportMenu(defaultFilename)
		if err == nil {
			if err := ExportTLE(tle, format, filePath); err != nil {
				fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to export: "+err.Error()))
			} else {
				fmt.Println(color.Ize(color.Green, fmt.Sprintf("  [+] Exported to: %s", filePath)))
			}
		}
	}
}
