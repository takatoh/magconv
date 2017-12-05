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
		os.Exit(0)
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

	copyx := [16]int{ 0,1,2,4,0,1,0,1,2,0,1,2,0,1,2,0 } 
	copyy := [16]int{ 0,0,0,0,1,1,2,2,2,4,4,4,8,8,8,16 }
	var copypos [16]int
	for i := 0; i < 16; i++ {
		copypos[i] = -(copyy[i] * int(header.Width) + (copyx[i] << pixelUnitLog))
	}

	data := make([]*mag.Palette, 0)
//	data := make([]byte, 0)
	src := 0
	dest := 0

	flagBuf := make([]byte, flagSize)
	var flagAPos int
	var flagBPos int

	var mask uint8 = 0x80
	for y := 0; y < int(header.Height); y++ {
		fmt.Println("---")
		for x := 0; x < int(flagSize); x++ {
			if flagA[flagAPos] & mask != 0x00 {
				flagBuf[x] = flagBuf[x] ^ flagB[flagBPos]
				flagBPos++
			}
			mask = mask >> 1
			if mask == 0 {
				mask = 0x80
				flagAPos++
			}
		}
		for x := 0; x < int(flagSize); x++ {
			vv := flagBuf[x]
			v := vv >> 4
			if v == 0 {
				fmt.Printf("(%d, %d) %d: ", x, y, v)
				if header.Colors == 16 {
					c := (pixel[src] >> 4)
					fmt.Printf("%v,", palettes[c])
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					fmt.Printf("%v,", palettes[c])
					data = append(data, palettes[c])
					src++
					c = (pixel[src] >> 4)
					fmt.Printf("%v,", palettes[c])
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					fmt.Printf("%v\n", palettes[c])
					data = append(data, palettes[c])
					src++
					dest += 4
				}
			} else {
				fmt.Printf("(%d, %d) %d: ", x, y, v)
				if header.Colors == 16 {
					copySrc := dest + copypos[v]
					fmt.Printf("%v,%v,%v,%v\n", data[copySrc], data[copySrc + 1], data[copySrc + 2], data[copySrc + 3])
					data = append(data, data[copySrc])
					data = append(data, data[copySrc + 1])
					data = append(data, data[copySrc + 2])
					data = append(data, data[copySrc + 3])
					dest += 4
				}
			}
			v = vv & 0xf
			if v == 0 {
				fmt.Printf("(%d, %d) %d: ", x, y, v)
				if header.Colors == 16 {
					c := (pixel[src] >> 4)
					fmt.Printf("%v,", palettes[c])
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					fmt.Printf("%v,", palettes[c])
					data = append(data, palettes[c])
					src++
					c = (pixel[src] >> 4)
					fmt.Printf("%v,", palettes[c])
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					fmt.Printf("%v\n", palettes[c])
					data = append(data, palettes[c])
					src++
					dest += 4
				}
			} else {
				fmt.Printf("(%d, %d) %d: ", x, y, v)
				if header.Colors == 16 {
					copySrc := dest + copypos[v]
					fmt.Printf("%v,%v,%v,%v\n", data[copySrc], data[copySrc + 1], data[copySrc + 2], data[copySrc + 3])
					data = append(data, data[copySrc])
					data = append(data, data[copySrc + 1])
					data = append(data, data[copySrc + 2])
					data = append(data, data[copySrc + 3])
					dest += 4
				}
			}
		}
	}
}
