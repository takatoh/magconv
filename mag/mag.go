package mag

import (
	"os"
//	"io"
	"bufio"
	"bytes"
	"encoding/binary"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Header struct {
	Colors     int
	StartX     uint16
	StartY     uint16
	EndX       uint16
	EndY       uint16
	FlgAOffset uint32
	FlgASize   uint32
	FlgBOffset uint32
	FlgBSize   uint32
	PxOffset   uint32
	PxSize     uint32
	Width      uint16
	Height     uint16
}

func NewHeader() *Header {
	p := new(Header)
	return p
}

type Palette struct {
	R uint8
	G uint8
	B uint8
}

func NewPalette(g, r, b uint8) *Palette {
	p := new(Palette)
	p.R = r
	p.G = g
	p.B = b
	return p
}

func CheckMag(file *os.File) bool {
	buf := make([]byte, 8)
	n, _ := file.Read(buf)
	if n != 8 {
		return false
	} else if string(buf) == "MAKI02  " {
		return true
	} else {
		return false
	}
}

func MachineCode(file *os.File) string {
	buf := make([]byte, 4)
	file.Read(buf)
	return string(buf)
}

func User(file *os.File) string {
	buf := make([]byte, 18 + 2)
	file.Read(buf)
	return convertFromShiftJIS(buf[0:18])
}

func Comment(file *os.File) string {
	c := make([]byte, 1)
	buf := make([]byte, 0)
	for {
		file.Read(c)
		if c[0] == 0x1A { break }
		buf = append(buf, c[0])
	}
	return convertFromShiftJIS(buf)
}

func convertFromShiftJIS(b []byte) string {
	r := strings.NewReader(string(b))
	s := bufio.NewScanner(transform.NewReader(r, japanese.ShiftJIS.NewDecoder()))
	list := make([]string, 0)
	for s.Scan() {
		list = append(list, s.Text())
	}
	return strings.Join(list, "")
}

func ReadUint8(file *os.File) uint8 {
	b := make([]byte, 1)
	file.Read(b)
	var val uint8
	binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &val)
	return val
}

func ReadUint16(file *os.File) uint16 {
	b := make([]byte, 2)
	file.Read(b)
	var val uint16
	binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &val)
	return val
}

func ReadUint32(file *os.File) uint32 {
	b := make([]byte, 4)
	file.Read(b)
	var val uint32
	binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &val)
	return val
}

func ReadHeader(file *os.File) *Header {
	header := NewHeader()

	for i := 0; i < 3; i++ {
		ReadUint8(file)
	}

	mode := ReadUint8(file)
	mode = mode >> 7
	if mode == 1 {
		header.Colors = 256
	} else {
		header.Colors = 16
	}

	header.StartX = ReadUint16(file)
	header.StartY = ReadUint16(file)
	header.EndX = ReadUint16(file)
	header.EndY = ReadUint16(file)
	header.FlgAOffset = ReadUint32(file)
	header.FlgBOffset = ReadUint32(file)
	header.FlgASize = header.FlgBOffset - header.FlgAOffset
	header.FlgBSize = ReadUint32(file)
	header.PxOffset = ReadUint32(file)
	header.PxSize = ReadUint32(file)
	header.Width = header.EndX - header.StartX + 1
	header.Height = header.EndY - header.StartY + 1

	return header
}

func ReadPalettes(file *os.File, n int) []*Palette {
	palettes := make([]*Palette, 0)
	for i := 0; i < n; i++ {
		g := ReadUint8(file) >> 4
		r := ReadUint8(file) >> 4
		b := ReadUint8(file) >> 4
		palettes = append(palettes, NewPalette(g, r, b))
	}
	return palettes
}
