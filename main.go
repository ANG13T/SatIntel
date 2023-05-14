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


func set_environmental_variable(env_key string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(env_key + ": ")
	input, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	os.Setenv(env_key, strings.TrimSpace(input))
}


func check_environmental_variable(env_key string) {
	if os.Getenv(env_key) == "" {
		set_environmental_variable(env_key)
	}
}


func main() {
	check_environmental_variable("SPACE_TRACK_USERNAME")
	check_environmental_variable("SPACE_TRACK_PASSWORD")
	check_environmental_variable("N2YO_API_KEY")

	cli.SatIntel()
}
