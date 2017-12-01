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
	fmt.Println("MAG format")

	machineCode := mag.MachineCode(magfile)
	fmt.Println(machineCode)

	user := mag.User(magfile)
	fmt.Println(user)

	comment := mag.Comment(magfile)
	fmt.Println(comment)

	for i := 0; i < 3; i++ {
		mag.ReadUint8(magfile)
//		fmt.Printf("%d\n", x)
	}

	var colors int
	mode := mag.ReadUint8(magfile)
	mode = mode >> 7
	if mode == 1 {
		colors = 256
	} else {
		colors = 16
	}
	fmt.Printf("colors=%d\n", colors)

	sx := mag.ReadUint16(magfile)
	sy := mag.ReadUint16(magfile)
	ex := mag.ReadUint16(magfile)
	ey := mag.ReadUint16(magfile)
	fmt.Println(sx, sy, ex, ey)

	flgAOffset := mag.ReadUint32(magfile)
	flgBOffset := mag.ReadUint32(magfile)
	flgASize := flgBOffset - flgAOffset
	flgBSize := mag.ReadUint32(magfile)
	pxOffset := mag.ReadUint32(magfile)
	pxSize := mag.ReadUint32(magfile)
	fmt.Printf("FlagA: offset=%d size=%d\n", flgAOffset, flgASize)
	fmt.Printf("FlagB: offset=%d size=%d\n", flgBOffset, flgBSize)
	fmt.Printf("Pixel: offset=%d size=%d\n", pxOffset, pxSize)
}
