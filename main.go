package main

import (
	"fmt"
	"os"
	"flag"

	"github.com/takatoh/magconv/mag"
)

func main() {
	opt_information := flag.Bool("information", false, "Display informations.")
	flag.Parse()

	filename := flag.Args()[0]
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
//	fmt.Println("MAG format")

	mag.MachineCode(magfile)
//	fmt.Println(machineCode)
	user := mag.User(magfile)
//	fmt.Println(user)
	comment := mag.Comment(magfile)
//	fmt.Println(comment)

	header := mag.ReadHeader(magfile)

	if *opt_information {
		fmt.Printf("user=%s\n", user)
		fmt.Printf("comment=%s\n", comment)
		fmt.Printf("colors=%d\n", header.Colors)
//		fmt.Println(header.StartX, header.StartY, header.EndX, header.EndY)
		fmt.Printf("width=%d, height=%d\n", header.Width, header.Height)
		fmt.Printf("FlagA: offset=%d size=%d\n", header.FlgAOffset, header.FlgASize)
		fmt.Printf("FlagB: offset=%d size=%d\n", header.FlgBOffset, header.FlgBSize)
		fmt.Printf("Pixel: offset=%d size=%d\n", header.PxOffset, header.PxSize)
	}
}
