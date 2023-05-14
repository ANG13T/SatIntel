/*
Copyright © 2023 Angelina Tsuboi angelinatsuboi@proton.me
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ANG13T/SatIntel/cli"
)

func setEnvironmentalVariable(envKey string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", envKey)
	input, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	input = strings.TrimSuffix(input, "\n")

	if err := os.Setenv(envKey, input); err != nil {
		fmt.Printf("Error setting environment variable %s: %v\n", envKey, err)
		os.Exit(1)
	}

	return input
}


func checkEnvironmentalVariable(envKey string) {
	_, found := os.LookupEnv(envKey)
	if !found {
		setEnvironmentalVariable(envKey)
	}
}


func main() {
	checkEnvironmentalVariable("SPACE_TRACK_USERNAME")
	checkEnvironmentalVariable("SPACE_TRACK_PASSWORD")
	checkEnvironmentalVariable("N2YO_API_KEY")

	cli.SatIntel()
}
