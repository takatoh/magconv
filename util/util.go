package util

import (
	"os"
	"strings"
	"bufio"
	"bytes"
	"encoding/binary"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func DecodeShiftJIS(b []byte) string {
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
