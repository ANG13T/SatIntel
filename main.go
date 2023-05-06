/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"os"
	"github.com/ANG13T/SatIntel/cli"
)

func main() {
	os.Setenv("SPACE_TRACK_USERNAME", "username")
	os.Setenv("SPACE_TRACK_PASSWORD", "password")
	cli.SatIntel()
}
