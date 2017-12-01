package mag

import (
	"os"
//	"io"
	"bufio"
)

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
	buf := make([]byte, 18)
	file.Read(buf)
	return string(buf)
}

func Comment(file *os.File) string {
	r := bufio.NewReader(file)
	buf := make([]byte, 0)
	for {
		c, _ := r.ReadByte()
		if c == 0x1A { break }
		buf = append(buf, c)
	}
	return string(buf)
}
