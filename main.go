package main

import (
	"fmt"
	"os"
	"flag"

	"github.com/takatoh/magconv/mag"
)

func printFlag(flag []byte, name string, size uint32) {
	flagLen := len(flag)
	fmt.Printf("%s: %d\n", name, flagLen)
	var i uint32
	for i = 0; i < size; i++ {
		fmt.Printf("%08b\n", flag[i])
	}
}

func main() {
	opt_info := flag.Bool("info", false, "Display informations.")
	opt_printflag := flag.Bool("printflag", false, "Print flag A and B.")
	opt_palettes := flag.Bool("palettes", false, "Print palettes.")
	opt_pixels := flag.Bool("pixels", false, "Print pixels.")
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

	if *opt_info {
		fmt.Printf("user=%s\n", user)
		fmt.Printf("comment=%s\n", comment)
		fmt.Printf("colors=%d\n", header.Colors)
		fmt.Printf("width=%d, height=%d\n", header.Width, header.Height)
		fmt.Printf("FlagA: offset=%d size=%d\n", header.FlgAOffset, header.FlgASize)
		fmt.Printf("FlagB: offset=%d size=%d\n", header.FlgBOffset, header.FlgBSize)
		fmt.Printf("Pixel: offset=%d size=%d\n", header.PxOffset, header.PxSize)

		os.Exit(0)
	}

	palettes := mag.ReadPalettes(magfile, header.Colors)
	if *opt_palettes {
		fmt.Println("Palettes:")
		for i, palette := range palettes {
			fmt.Printf("%d: r=%02x, g=%02x, b=%02x\n", i, palette.R, palette.G, palette.B)
		}
	}

	flagA := mag.ReadFlagA(magfile, header.FlgASize)
	flagB := mag.ReadFlagB(magfile, header.FlgBSize)
	if *opt_printflag {
		printFlag(flagA, "Flag A", header.FlgASize)
		printFlag(flagB, "Flag B", header.FlgBSize)
	}

	pixel := mag.ReadPixel(magfile, header.PxSize)
	if *opt_pixels {
		printFlag(pixel, "Pixels", header.PxSize)
	}

	var pixelUnitLog uint
	if header.Colors == 256 {
		pixelUnitLog = 1
	} else {
		pixelUnitLog = 2
	}
	flagSize := header.Width >> (pixelUnitLog + 1)
	fmt.Printf("flag size=%d\n", flagSize)
}
