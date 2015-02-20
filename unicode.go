package JSONParse

import (
	"fmt"
	"unicode/utf16"
	"strings"
)

const (
	replacementChar = '\uFFFD'
)

func unicodeLen(str string) int {
	index := strings.Index(str, `\u`)
	for {
		if index == -1 {
			break
		}

		r1str:= str[index+2:index+6] + " " + str[index+8:index+12]

		var r1, r2	rune

		fmt.Sscanf(r1str, "%x %x", &r1, &r2)
		Trace.Println("runes", r1, r2)
		if utf16.IsSurrogate(r1) && utf16.IsSurrogate(r2) {
			r := utf16.DecodeRune(r1, r2)
			if r == replacementChar {
				Warning.Println("Unable to deccode runes")
				str = strings.Replace(str, `\u`, "0x", 2)
			} else {
				str = str[:index] + " " + str[index+12:len(str)]
			}
		}
		index = strings.Index(str, `\u`)
	}
	return len(str)
}
