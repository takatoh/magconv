package mag

import (
	"os"
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
