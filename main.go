package main

import (
	"fmt"
	"os"
	"flag"
	"image"
	"image/color"
	"image/png"
	"path"
	"strings"

	"github.com/takatoh/magconv/mag"
)

const (
	progVersion = "v0.1.1"
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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
`Usage:
  %s [options] <file.mag>
Options:
`, os.Args[0])
		flag.PrintDefaults()
	}
	opt_info := flag.Bool("info", false, "Display informations.")
	opt_flags := flag.Bool("flags", false, "Print flag A and B.")
	opt_palettes := flag.Bool("palettes", false, "Print palettes.")
	opt_pixels := flag.Bool("pixels", false, "Print pixels.")
	opt_version := flag.Bool("version", false, "Show version.")
	flag.Parse()

	if *opt_version {
		fmt.Println(progVersion)
		os.Exit(0)
	}

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

	if *opt_palettes {
		palettes := mag.ReadPalettes(magfile, header.Colors)
		fmt.Println("Palettes:")
		for i, palette := range palettes {
			fmt.Printf("%d: r=%02x, g=%02x, b=%02x\n", i, palette.R, palette.G, palette.B)
		}
		os.Exit(0)
	}

	if *opt_flags {
		mag.ReadPalettes(magfile, header.Colors)
		flagA := mag.ReadFlagA(magfile, header.FlgASize)
		flagB := mag.ReadFlagB(magfile, header.FlgBSize)
		printFlag(flagA, "Flag A", header.FlgASize)
		printFlag(flagB, "Flag B", header.FlgBSize)
		os.Exit(0)
	}

	if *opt_pixels {
		mag.ReadPalettes(magfile, header.Colors)
		mag.ReadFlagA(magfile, header.FlgASize)
		mag.ReadFlagB(magfile, header.FlgBSize)
		pixel := mag.ReadPixel(magfile, header.PxSize)
		printFlag(pixel, "Pixels", header.PxSize)
		os.Exit(0)
	}

	result := mag.Load(magfile, header)

	w := int(header.Width)
	h := int(header.Height)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r := result[y][x].R * 16
			g := result[y][x].G * 16
			b := result[y][x].B * 16
			c := color.RGBA{r, g, b, 255}
			img.Set(x, y, c)
		}
	}

	ext := path.Ext(filename)
	pngFilename := strings.Replace(filename, ext, ".png", 1)
	f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
