package main

import (
	"fmt"
	"os"

	"github.com/takatoh/magconv/mag"
)

func main() {
	filename := os.Args[1]
	magfile, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open file: %s\n", filename)
		os.Exit(1)
	}
	defer magfile.Close()

	check := mag.CheckMag(magfile)
	if !check {
		fmt.Fprintln(os.Stderr, "Not MAG format")
		os.Exit(0)
	}
	mmt.Println("Not MAG format")
}
