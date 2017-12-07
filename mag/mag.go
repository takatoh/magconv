package mag

import (
	"os"

	"github.com/takatoh/magconv/util"
)

type Loader struct {
	CheckMag bool
	MachineCode string
	User string
	Comment string
	Header *Header
	Palettes []*Palette
	FlagA []byte
	FlagB []byte
	Pixel []byte	
}

func NewLoader(file *os.File) *Loader {
	loader := new(Loader)
	loader.CheckMag = CheckMag(file)
	loader.MachineCode = MachineCode(file)
	loader.User = User(file)
	loader.Comment = Comment(file)
	loader.Header = ReadHeader(file)
	loader.Palettes = ReadPalettes(file, loader.Header.Colors)
	loader.FlagA = ReadFlagA(file, loader.Header.FlgASize)
	loader.FlagB = ReadFlagB(file, loader.Header.FlgBSize)
	loader.Pixel = ReadPixel(file, loader.Header.PxSize)
	return loader
}

func (l *Loader) Load() [][]*Palette {
	header := l.Header
	palettes := l.Palettes
	flagA := l.FlagA
	flagB := l.FlagB
	pixel := l.Pixel

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

	data := make([]*Palette, 0)
	src := 0
	dest := 0

	flagBuf := make([]byte, flagSize)
	var flagAPos int
	var flagBPos int

	var mask uint8 = 0x80
	for y := 0; y < int(header.Height); y++ {
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
				if header.Colors == 16 {
					c := (pixel[src] >> 4)
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					data = append(data, palettes[c])
					src++
					c = (pixel[src] >> 4)
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					data = append(data, palettes[c])
					src++
					dest += 4
				}
			} else {
				if header.Colors == 16 {
					copySrc := dest + copypos[v]
					data = append(data, data[copySrc])
					data = append(data, data[copySrc + 1])
					data = append(data, data[copySrc + 2])
					data = append(data, data[copySrc + 3])
					dest += 4
				}
			}
			v = vv & 0xf
			if v == 0 {
				if header.Colors == 16 {
					c := (pixel[src] >> 4)
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					data = append(data, palettes[c])
					src++
					c = (pixel[src] >> 4)
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					data = append(data, palettes[c])
					src++
					dest += 4
				}
			} else {
				if header.Colors == 16 {
					copySrc := dest + copypos[v]
					data = append(data, data[copySrc])
					data = append(data, data[copySrc + 1])
					data = append(data, data[copySrc + 2])
					data = append(data, data[copySrc + 3])
					dest += 4
				}
			}
		}
	}

	result := make([][]*Palette, 0)
	for y := 0; y < int(header.Height); y++ {
		s := y * int(header.Width)
		e := s + int(header.Width)
		result = append(result, data[s:e])
	}
	return result
}

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
	return util.ConvertFromShiftJIS(buf[0:18])
}

func Comment(file *os.File) string {
	c := make([]byte, 1)
	buf := make([]byte, 0)
	for {
		file.Read(c)
		if c[0] == 0x1A { break }
		buf = append(buf, c[0])
	}
	return util.ConvertFromShiftJIS(buf)
}

func ReadHeader(file *os.File) *Header {
	header := NewHeader()

	for i := 0; i < 3; i++ {
		util.ReadUint8(file)
	}

	mode := util.ReadUint8(file)
	mode = mode >> 7
	if mode == 1 {
		header.Colors = 256
	} else {
		header.Colors = 16
	}

	header.StartX = util.ReadUint16(file)
	header.StartY = util.ReadUint16(file)
	header.EndX = util.ReadUint16(file)
	header.EndY = util.ReadUint16(file)
	header.FlgAOffset = util.ReadUint32(file)
	header.FlgBOffset = util.ReadUint32(file)
	header.FlgASize = header.FlgBOffset - header.FlgAOffset
	header.FlgBSize = util.ReadUint32(file)
	header.PxOffset = util.ReadUint32(file)
	header.PxSize = util.ReadUint32(file)
	header.Width = header.EndX - header.StartX + 1
	header.Height = header.EndY - header.StartY + 1

	return header
}

func ReadPalettes(file *os.File, n int) []*Palette {
	palettes := make([]*Palette, 0)
	for i := 0; i < n; i++ {
		g := util.ReadUint8(file) >> 4
		r := util.ReadUint8(file) >> 4
		b := util.ReadUint8(file) >> 4
		palettes = append(palettes, NewPalette(g, r, b))
	}
	return palettes
}

func ReadFlagA(file *os.File, size uint32) []byte {
	flgA := make([]byte, size)
	file.Read(flgA)
	return flgA
}

func ReadFlagB(file *os.File, size uint32) []byte {
	flgB := make([]byte, size)
	file.Read(flgB)
	return flgB
}

func ReadPixel(file *os.File, size uint32) []byte {
	pxl := make([]byte, size)
	file.Read(pxl)
	return pxl
}

func Load(file *os.File, header *Header) [][]*Palette {
	palettes := ReadPalettes(file, header.Colors)
	flagA := ReadFlagA(file, header.FlgASize)
	flagB := ReadFlagB(file, header.FlgBSize)
	pixel := ReadPixel(file, header.PxSize)

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

	data := make([]*Palette, 0)
	src := 0
	dest := 0

	flagBuf := make([]byte, flagSize)
	var flagAPos int
	var flagBPos int

	var mask uint8 = 0x80
	for y := 0; y < int(header.Height); y++ {
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
				if header.Colors == 16 {
					c := (pixel[src] >> 4)
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					data = append(data, palettes[c])
					src++
					c = (pixel[src] >> 4)
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					data = append(data, palettes[c])
					src++
					dest += 4
				}
			} else {
				if header.Colors == 16 {
					copySrc := dest + copypos[v]
					data = append(data, data[copySrc])
					data = append(data, data[copySrc + 1])
					data = append(data, data[copySrc + 2])
					data = append(data, data[copySrc + 3])
					dest += 4
				}
			}
			v = vv & 0xf
			if v == 0 {
				if header.Colors == 16 {
					c := (pixel[src] >> 4)
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					data = append(data, palettes[c])
					src++
					c = (pixel[src] >> 4)
					data = append(data, palettes[c])
					c = (pixel[src] & 0xf)
					data = append(data, palettes[c])
					src++
					dest += 4
				}
			} else {
				if header.Colors == 16 {
					copySrc := dest + copypos[v]
					data = append(data, data[copySrc])
					data = append(data, data[copySrc + 1])
					data = append(data, data[copySrc + 2])
					data = append(data, data[copySrc + 3])
					dest += 4
				}
			}
		}
	}

	result := make([][]*Palette, 0)
	for y := 0; y < int(header.Height); y++ {
		s := y * int(header.Width)
		e := s + int(header.Width)
		result = append(result, data[s:e])
	}
	return result
}
