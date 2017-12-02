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

	mag.MachineCode(magfile)
	user := mag.User(magfile)
	comment := mag.Comment(magfile)
	header := mag.ReadHeader(magfile)

	if *opt_information {
		fmt.Printf("user=%s\n", user)
		fmt.Printf("comment=%s\n", comment)
		fmt.Printf("colors=%d\n", header.Colors)
		fmt.Printf("width=%d, height=%d\n", header.Width, header.Height)
		fmt.Printf("FlagA: offset=%d size=%d\n", header.FlgAOffset, header.FlgASize)
		fmt.Printf("FlagB: offset=%d size=%d\n", header.FlgBOffset, header.FlgBSize)
		fmt.Printf("Pixel: offset=%d size=%d\n", header.PxOffset, header.PxSize)
		os.Exit(0)
	}

	pallets := mag.ReadPallets(magfile, header.Colors)
	fmt.Println("Pallets:")
	for i, pallet := range pallets {
		fmt.Printf("%d: r=%02x, g=%02x, b=%02x\n", i, pallet.R, pallet.G, pallet.B)
	}
}
