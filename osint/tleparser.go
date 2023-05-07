package osint

import (
	"fmt"
	"os"
	"github.com/TwiN/go-color"
	"github.com/iskaa02/qalam/gradient"	
	"io/ioutil"
	"bufio"
)

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
